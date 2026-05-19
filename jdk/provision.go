// Package jdk provisiona um JDK Temurin compatível para o Runner.
//
// A ordem de busca é: JAVA_HOME → PATH → cache local → download. O cache fica
// em <UserCacheDir>/runner/jdk/<major>/ e é compartilhado entre os CLIs do
// projeto (simulador, assinatura).
package jdk

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const RequiredMajor = 21

var versionRe = regexp.MustCompile(`version "([^"]+)"`)

// FindOrProvision retorna um executável java com versão >= RequiredMajor.
// Ordem de busca: JAVA_HOME → PATH → cache local → download Temurin.
func FindOrProvision(ctx context.Context) (string, error) {
	if jh := os.Getenv("JAVA_HOME"); jh != "" {
		if p, err := javaExeIn(jh); err == nil {
			if v, err := MajorVersion(p); err == nil && v >= RequiredMajor {
				return p, nil
			}
		}
	}

	if p, err := exec.LookPath("java"); err == nil {
		if v, err := MajorVersion(p); err == nil && v >= RequiredMajor {
			return p, nil
		}
	}

	if p, err := findCached(); err == nil {
		return p, nil
	}

	fmt.Fprintf(os.Stderr, "JDK %d não encontrado. Baixando Temurin JDK %d...\n", RequiredMajor, RequiredMajor)
	return download(ctx)
}

// MajorVersion executa java -version e retorna o número major (ex.: 21).
func MajorVersion(javaExe string) (int, error) {
	out, err := exec.Command(javaExe, "-version").CombinedOutput()
	if err != nil {
		return 0, err
	}
	m := versionRe.FindSubmatch(out)
	if m == nil {
		return 0, fmt.Errorf("versão não encontrada em: %s", out)
	}
	parts := strings.Split(string(m[1]), ".")
	major := parts[0]
	// Formato legado "1.8.0": major real é o segundo segmento.
	if major == "1" && len(parts) > 1 {
		major = parts[1]
	}
	return strconv.Atoi(major)
}

// CacheDir retorna o diretório de cache do JDK provisionado pelo Runner.
func CacheDir() (string, error) {
	base, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(base, "runner", "jdk", strconv.Itoa(RequiredMajor)), nil
}

func javaExeIn(home string) (string, error) {
	bin := filepath.Join(home, "bin")
	for _, name := range []string{"java", "java.exe"} {
		p := filepath.Join(bin, name)
		if fi, err := os.Stat(p); err == nil && !fi.IsDir() {
			return p, nil
		}
	}
	return "", fmt.Errorf("java não encontrado em %s", home)
}

func findCached() (string, error) {
	dir, err := CacheDir()
	if err != nil {
		return "", err
	}
	marker := filepath.Join(dir, "java.path")
	data, err := os.ReadFile(marker)
	if err != nil {
		return "", err
	}
	p := strings.TrimSpace(string(data))
	if _, err := os.Stat(p); err != nil {
		return "", fmt.Errorf("executável em cache não encontrado: %s", p)
	}
	return p, nil
}

// --- download ---

type adoptiumAsset struct {
	Binary struct {
		Package struct {
			Link string `json:"link"`
			Name string `json:"name"`
		} `json:"package"`
	} `json:"binary"`
}

func download(ctx context.Context) (string, error) {
	dir, err := CacheDir()
	if err != nil {
		return "", fmt.Errorf("diretório de cache: %w", err)
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("criar diretório de cache: %w", err)
	}

	dlURL, filename, err := resolveURL(ctx)
	if err != nil {
		return "", fmt.Errorf("resolver URL de download: %w", err)
	}

	tmpFile := filepath.Join(dir, filename)
	if err := downloadFile(ctx, dlURL, tmpFile); err != nil {
		_ = os.Remove(tmpFile)
		return "", fmt.Errorf("baixar JDK: %w", err)
	}
	defer os.Remove(tmpFile)

	fmt.Fprintln(os.Stderr, "Extraindo JDK...")
	extractDir := filepath.Join(dir, "home")
	_ = os.RemoveAll(extractDir)
	if err := os.MkdirAll(extractDir, 0o755); err != nil {
		return "", err
	}

	if strings.HasSuffix(filename, ".zip") {
		err = extractZip(tmpFile, extractDir)
	} else {
		err = extractTarGz(tmpFile, extractDir)
	}
	if err != nil {
		return "", fmt.Errorf("extrair JDK: %w", err)
	}

	javaPath, err := findJavaBinary(extractDir)
	if err != nil {
		return "", fmt.Errorf("localizar java no JDK extraído: %w", err)
	}

	marker := filepath.Join(dir, "java.path")
	if err := os.WriteFile(marker, []byte(javaPath), 0o644); err != nil {
		return "", err
	}

	fmt.Fprintf(os.Stderr, "JDK instalado em: %s\n", filepath.Dir(filepath.Dir(javaPath)))
	return javaPath, nil
}

func resolveURL(ctx context.Context) (dlURL, filename string, err error) {
	goos := adoptiumOS()
	goarch := adoptiumArch()
	apiURL := fmt.Sprintf(
		"https://api.adoptium.net/v3/assets/latest/%d/hotspot?architecture=%s&image_type=jdk&os=%s&vendor=eclipse",
		RequiredMajor, goarch, goos,
	)

	client := &http.Client{Timeout: 15 * time.Second}
	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return "", "", err
	}
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("API Adoptium: %w", err)
	}
	defer resp.Body.Close()

	var assets []adoptiumAsset
	if err := json.NewDecoder(resp.Body).Decode(&assets); err != nil {
		return "", "", fmt.Errorf("decodificar resposta Adoptium: %w", err)
	}
	if len(assets) == 0 {
		return "", "", fmt.Errorf("nenhum asset encontrado (os=%s arch=%s)", goos, goarch)
	}

	pkg := assets[0].Binary.Package
	return pkg.Link, pkg.Name, nil
}

func downloadFile(ctx context.Context, url, dest string) error {
	client := &http.Client{Timeout: 10 * time.Minute}
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return err
	}
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
			if _, werr := f.Write(buf[:n]); werr != nil {
				return werr
			}
			downloaded += int64(n)
			if total > 0 {
				fmt.Fprintf(os.Stderr, "\r  %d%% (%d / %d MB)",
					downloaded*100/total, downloaded>>20, total>>20)
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

// --- extração ---

func extractTarGz(src, dest string) error {
	f, err := os.Open(src)
	if err != nil {
		return err
	}
	defer f.Close()

	gz, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer gz.Close()

	tr := tar.NewReader(gz)
	cleanDest := filepath.Clean(dest) + string(os.PathSeparator)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		target := filepath.Join(dest, filepath.FromSlash(hdr.Name))
		if !strings.HasPrefix(target, cleanDest) {
			continue
		}
		switch hdr.Typeflag {
		case tar.TypeDir:
			_ = os.MkdirAll(target, 0o755)
		case tar.TypeReg:
			_ = os.MkdirAll(filepath.Dir(target), 0o755)
			out, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(hdr.Mode))
			if err != nil {
				return err
			}
			_, err = io.Copy(out, tr)
			out.Close()
			if err != nil {
				return err
			}
		case tar.TypeSymlink:
			_ = os.MkdirAll(filepath.Dir(target), 0o755)
			_ = os.Symlink(hdr.Linkname, target)
		}
	}
	return nil
}

func extractZip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	cleanDest := filepath.Clean(dest) + string(os.PathSeparator)
	for _, f := range r.File {
		target := filepath.Join(dest, filepath.FromSlash(f.Name))
		if !strings.HasPrefix(target, cleanDest) {
			continue
		}
		if f.FileInfo().IsDir() {
			_ = os.MkdirAll(target, 0o755)
			continue
		}
		_ = os.MkdirAll(filepath.Dir(target), 0o755)
		rc, err := f.Open()
		if err != nil {
			return err
		}
		out, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, f.Mode())
		if err != nil {
			rc.Close()
			return err
		}
		_, err = io.Copy(out, rc)
		out.Close()
		rc.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func findJavaBinary(dir string) (string, error) {
	var found string
	_ = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil || found != "" {
			return nil
		}
		if info.IsDir() {
			return nil
		}
		name := info.Name()
		if (name == "java" || name == "java.exe") && filepath.Base(filepath.Dir(path)) == "bin" {
			found = path
			return filepath.SkipAll
		}
		return nil
	})
	if found == "" {
		return "", fmt.Errorf("java não encontrado em %s", dir)
	}
	return found, nil
}

func adoptiumOS() string {
	switch runtime.GOOS {
	case "darwin":
		return "mac"
	case "windows":
		return "windows"
	default:
		return "linux"
	}
}

func adoptiumArch() string {
	if runtime.GOARCH == "arm64" {
		return "aarch64"
	}
	return "x64"
}
