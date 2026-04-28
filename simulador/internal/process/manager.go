package process

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	// EnvSimuladorRepo define o repositório GitHub (owner/repo) onde o simulador.jar
	// é publicado. Exemplo: "minha-org/hubsaude-simulador".
	EnvSimuladorRepo = "RUNNER_SIMULADOR_REPO"

	// EnvSimuladorJar permite apontar diretamente para um JAR local, ignorando download.
	EnvSimuladorJar = "RUNNER_SIMULADOR_JAR"

	DefaultPort = 8080
)

// CacheDir retorna o diretório de cache do simulador.
func CacheDir() (string, error) {
	base, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(base, "runner", "simulador"), nil
}

// FindOrDownload retorna o caminho para o simulador.jar.
// Se RUNNER_SIMULADOR_JAR estiver definido, usa esse caminho diretamente.
// Caso contrário, baixa do GitHub Releases se necessário.
func FindOrDownload(ctx context.Context) (string, error) {
	if p := os.Getenv(EnvSimuladorJar); p != "" {
		if _, err := os.Stat(p); err != nil {
			return "", fmt.Errorf("%s: arquivo não encontrado: %s", EnvSimuladorJar, p)
		}
		return p, nil
	}

	repo := os.Getenv(EnvSimuladorRepo)
	if repo == "" {
		return "", fmt.Errorf(
			"repositório do simulador não configurado.\n"+
				"Defina %s=<owner>/<repo> com o repositório GitHub que publica o simulador.jar.\n"+
				"Ou defina %s=<caminho> para usar um JAR local.",
			EnvSimuladorRepo, EnvSimuladorJar,
		)
	}

	dir, err := CacheDir()
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}

	latest, dlURL, err := latestRelease(ctx, repo)
	if err != nil {
		// Se não conseguir verificar, usa JAR em cache se existir.
		cached := filepath.Join(dir, "simulador.jar")
		if _, serr := os.Stat(cached); serr == nil {
			fmt.Fprintf(os.Stderr, "Aviso: não foi possível verificar atualizações (%v). Usando versão em cache.\n", err)
			return cached, nil
		}
		return "", fmt.Errorf("verificar release: %w", err)
	}

	jarPath := filepath.Join(dir, "simulador.jar")
	versionFile := filepath.Join(dir, "version.txt")

	if cached, _ := os.ReadFile(versionFile); strings.TrimSpace(string(cached)) == latest {
		if _, err := os.Stat(jarPath); err == nil {
			return jarPath, nil
		}
	}

	fmt.Fprintf(os.Stderr, "Baixando simulador.jar (%s)...\n", latest)
	if err := downloadFile(ctx, dlURL, jarPath); err != nil {
		return "", fmt.Errorf("baixar simulador.jar: %w", err)
	}
	_ = os.WriteFile(versionFile, []byte(latest), 0o644)
	fmt.Fprintln(os.Stderr, "simulador.jar atualizado.")
	return jarPath, nil
}

// CheckPort retorna erro se a porta TCP já estiver em uso.
// Testa tanto IPv4 quanto IPv6 para detectar conflito em qualquer stack.
func CheckPort(port int) error {
	addr := fmt.Sprintf(":%d", port)
	for _, network := range []string{"tcp4", "tcp6"} {
		ln, err := net.Listen(network, addr)
		if err != nil {
			return fmt.Errorf("porta %d já está em uso: %w", port, err)
		}
		ln.Close()
	}
	return nil
}

// Start inicia o simulador.jar como processo em background após verificar a porta.
func Start(javaExe, jarPath string, port int) error {
	if err := CheckPort(port); err != nil {
		return err
	}

	cmd := exec.Command(javaExe, "-jar", jarPath)
	detachProcess(cmd)

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("iniciar simulador: %w", err)
	}

	pid := cmd.Process.Pid
	_ = cmd.Process.Release()

	dir, err := CacheDir()
	if err != nil {
		return err
	}
	pidContent := fmt.Sprintf("%d\n%d\n", pid, port)
	return os.WriteFile(filepath.Join(dir, "simulador.pid"), []byte(pidContent), 0o644)
}

// Stop encerra o processo do simulador usando o PID salvo.
func Stop() error {
	pid, _, err := readPIDFile()
	if err != nil {
		return fmt.Errorf("simulador não está em execução (PID não encontrado): %w", err)
	}

	proc, err := os.FindProcess(pid)
	if err != nil {
		return fmt.Errorf("processo PID %d não encontrado: %w", pid, err)
	}

	if !isProcessRunning(proc) {
		_ = removePIDFile()
		return fmt.Errorf("simulador (PID %d) não está em execução", pid)
	}

	if err := terminateProcess(proc); err != nil {
		return fmt.Errorf("encerrar processo %d: %w", pid, err)
	}

	_ = removePIDFile()
	return nil
}

// Status retorna se o simulador está em execução, seu PID e porta.
func Status() (running bool, pid, port int) {
	p, pt, err := readPIDFile()
	if err != nil {
		return false, 0, 0
	}
	proc, err := os.FindProcess(p)
	if err != nil {
		return false, 0, 0
	}
	if !isProcessRunning(proc) {
		_ = removePIDFile()
		return false, 0, 0
	}
	return true, p, pt
}

// --- helpers internos ---

type githubRelease struct {
	TagName string `json:"tag_name"`
	Assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

func latestRelease(ctx context.Context, repo string) (tag, downloadURL string, err error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", repo)
	client := &http.Client{Timeout: 15 * time.Second}
	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return "", "", fmt.Errorf("release não encontrada em %s (repositório existe?)", repo)
	}
	if resp.StatusCode != 200 {
		return "", "", fmt.Errorf("GitHub API retornou %d", resp.StatusCode)
	}

	var rel githubRelease
	if err := json.NewDecoder(resp.Body).Decode(&rel); err != nil {
		return "", "", err
	}

	for _, a := range rel.Assets {
		if strings.HasSuffix(strings.ToLower(a.Name), ".jar") &&
			strings.Contains(strings.ToLower(a.Name), "simulador") {
			return rel.TagName, a.BrowserDownloadURL, nil
		}
	}
	return "", "", fmt.Errorf("asset simulador*.jar não encontrado na release %s do repo %s", rel.TagName, repo)
}

func downloadFile(ctx context.Context, url, dest string) error {
	client := &http.Client{Timeout: 10 * time.Minute}
	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	f, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer f.Close()

	total := resp.ContentLength
	var downloaded int64
	buf := make([]byte, 64*1024)
	for {
		n, readErr := resp.Body.Read(buf)
		if n > 0 {
			f.Write(buf[:n])
			downloaded += int64(n)
			if total > 0 {
				fmt.Fprintf(os.Stderr, "\r  %d%% (%d / %d MB)", downloaded*100/total, downloaded>>20, total>>20)
			} else {
				fmt.Fprintf(os.Stderr, "\r  %d MB baixados", downloaded>>20)
			}
		}
		if readErr == io.EOF {
			break
		}
		if readErr != nil {
			return readErr
		}
	}
	fmt.Fprintln(os.Stderr)
	return nil
}

func readPIDFile() (pid, port int, err error) {
	dir, err := CacheDir()
	if err != nil {
		return 0, 0, err
	}
	data, err := os.ReadFile(filepath.Join(dir, "simulador.pid"))
	if err != nil {
		return 0, 0, err
	}
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if len(lines) < 2 {
		return 0, 0, fmt.Errorf("PID file inválido")
	}
	pid, err = strconv.Atoi(strings.TrimSpace(lines[0]))
	if err != nil {
		return 0, 0, err
	}
	port, err = strconv.Atoi(strings.TrimSpace(lines[1]))
	return pid, port, err
}

func removePIDFile() error {
	dir, err := CacheDir()
	if err != nil {
		return err
	}
	return os.Remove(filepath.Join(dir, "simulador.pid"))
}
