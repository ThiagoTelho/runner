package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// ── cores ANSI ──────────────────────────────────────────────────────────────

const (
	reset  = "\033[0m"
	bold   = "\033[1m"
	cyan   = "\033[36m"
	green  = "\033[32m"
	yellow = "\033[33m"
	red    = "\033[31m"
	gray   = "\033[90m"
	bCyan  = "\033[96m"
)

// ── estado global da sessão ──────────────────────────────────────────────────

var (
	binAssinatura string
	binSimulador  string
	jarAssinador  string
	workDir       string
	lastJWSFile   string
	reader        = bufio.NewReader(os.Stdin)
)

func main() {
	clearScreen()
	if err := setup(); err != nil {
		fmt.Printf("%s%sErro na inicialização: %v%s\n", bold, red, err, reset)
		os.Exit(1)
	}
	menuPrincipal()
}

// ── setup ────────────────────────────────────────────────────────────────────

func setup() error {
	// Localiza o diretório raiz do repositório (pai de demo/).
	self, err := os.Executable()
	if err != nil {
		self = os.Args[0]
	}
	self, _ = filepath.EvalSymlinks(self)
	repoRoot := filepath.Dir(filepath.Dir(self))

	// Binários Go.
	binAssinatura = filepath.Join(repoRoot, "assinatura", exeName("assinatura"))
	binSimulador = filepath.Join(repoRoot, "simulador", exeName("simulador"))

	// JAR do assinador.
	jarAssinador = os.Getenv("RUNNER_ASSINADOR_JAR")
	if jarAssinador == "" {
		candidate := filepath.Join(repoRoot, "assinador", "build", "libs", "assinador-0.1.0.jar")
		if _, err := os.Stat(candidate); err == nil {
			jarAssinador = candidate
		}
	}

	// Diretório de trabalho com arquivos de exemplo.
	workDir = filepath.Join(os.TempDir(), "runner-demo")
	if err := os.MkdirAll(workDir, 0o755); err != nil {
		return err
	}
	return createSampleFiles()
}

func exeName(base string) string {
	if runtime.GOOS == "windows" {
		return base + ".exe"
	}
	return base
}

func createSampleFiles() error {
	files := map[string]string{
		"bundle.json": `{
  "resourceType": "Bundle",
  "id": "demo-bundle",
  "type": "document",
  "entry": [{"fullUrl": "urn:uuid:demo"}]
}`,
		"provenance.json": `{
  "resourceType": "Provenance",
  "id": "demo-provenance",
  "recorded": "2026-01-01T00:00:00Z",
  "agent": [{"type": {"text": "author"}}]
}`,
		"material-pem.json": `{
  "tipo": "PEM",
  "chavePrivada": "-----BEGIN PRIVATE KEY-----\nMIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQC7\n-----END PRIVATE KEY-----"
}`,
		"material-pkcs12.json": `{
  "tipo": "PKCS12",
  "alias": "meu-certificado",
  "senha": "senha123",
  "conteudo": "MIIJkAIBAzCCCVYGCSqGSIb3DQEHAaCCCUcEgglDMIIJPzCC"
}`,
		"material-token.json": `{
  "tipo": "TOKEN",
  "tokenLabel": "MeuToken",
  "slotId": 0,
  "identificador": "chave-01",
  "pin": "1234"
}`,
	}
	for name, content := range files {
		path := filepath.Join(workDir, name)
		if _, err := os.Stat(path); err != nil {
			if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
				return err
			}
		}
	}
	return nil
}

// ── menus ────────────────────────────────────────────────────────────────────

func menuPrincipal() {
	for {
		clearScreen()
		printHeader()
		fmt.Println()
		printItem("1", "Assinador — criar assinatura")
		printItem("2", "Assinador — validar assinatura")
		printItem("3", "Assinador — servidor HTTP  (iniciar / parar / status)")
		fmt.Println()
		printItem("4", "Simulador HubSaúde  (iniciar / parar / status)")
		fmt.Println()
		printItem("0", "Sair")
		fmt.Println()

		switch prompt("Escolha") {
		case "1":
			menuCriar()
		case "2":
			menuValidar()
		case "3":
			menuServidor()
		case "4":
			menuSimulador()
		case "0", "q", "sair":
			clearScreen()
			fmt.Printf("%sTé mais!%s\n\n", cyan, reset)
			return
		}
	}
}

// ── menu: criar assinatura ───────────────────────────────────────────────────

func menuCriar() {
	for {
		clearScreen()
		printHeader()
		printSubtitle("Criar Assinatura")
		fmt.Println()
		printItem("1", "Material PEM        (chave privada)")
		printItem("2", "Material PKCS12     (certificado + senha)")
		printItem("3", "Material TOKEN      (PKCS#11 / smart card)")
		fmt.Println()
		printItem("v", "Ver arquivos de entrada")
		printItem("0", "Voltar")
		fmt.Println()

		op := prompt("Escolha o tipo de material")
		switch op {
		case "0":
			return
		case "v":
			verArquivos()
		case "1", "2", "3":
			materiais := map[string]string{
				"1": "material-pem.json",
				"2": "material-pkcs12.json",
				"3": "material-token.json",
			}
			runCriar(materiais[op])
		}
	}
}

func runCriar(materialFile string) {
	clearScreen()
	printHeader()
	printSubtitle("Criar Assinatura — executando")
	fmt.Println()

	jws, err := runAssinatura(
		"criar",
		"-b", filepath.Join(workDir, "bundle.json"),
		"-p", filepath.Join(workDir, "provenance.json"),
		"-m", filepath.Join(workDir, materialFile),
	)

	if err == nil {
		// Extrai o JWS da saída para uso posterior em "validar".
		for _, line := range strings.Split(jws, "\n") {
			if strings.HasPrefix(strings.TrimSpace(line), "ey") {
				lastJWSFile = filepath.Join(workDir, "ultimo.jws")
				_ = os.WriteFile(lastJWSFile, []byte(strings.TrimSpace(line)), 0o644)
				break
			}
		}
	}

	aguardar()
}

// ── menu: validar assinatura ─────────────────────────────────────────────────

func menuValidar() {
	clearScreen()
	printHeader()
	printSubtitle("Validar Assinatura")
	fmt.Println()

	if lastJWSFile == "" {
		fmt.Printf("%sNenhum JWS disponível na sessão.%s\n", yellow, reset)
		fmt.Printf("%sDica: execute \"Criar Assinatura\" primeiro para gerar um JWS.%s\n\n", gray, reset)
		aguardar()
		return
	}

	fmt.Printf("%sUsando JWS da última criação:%s %s\n\n", gray, reset, lastJWSFile)

	runAssinatura(
		"validar",
		"-j", lastJWSFile,
		"-r", "warn",
		"-t", fmt.Sprintf("%d", time.Now().Unix()),
		"-a", "https://politica.exemplo.gov.br/v1",
	)

	aguardar()
}

// ── menu: servidor HTTP ───────────────────────────────────────────────────────

func menuServidor() {
	for {
		clearScreen()
		printHeader()
		printSubtitle("Servidor HTTP do Assinador")
		fmt.Println()
		printItem("1", "Status")
		printItem("2", "Iniciar  (porta 8190)")
		printItem("3", "Iniciar com timeout de inatividade")
		printItem("4", "Parar")
		fmt.Println()
		printItem("0", "Voltar")
		fmt.Println()

		switch prompt("Escolha") {
		case "0":
			return
		case "1":
			runAssinatura("servidor", "status")
			aguardar()
		case "2":
			runAssinatura("servidor", "iniciar")
			aguardar()
		case "3":
			minutos := prompt("Encerrar após quantos minutos sem interação?")
			runAssinatura("servidor", "iniciar", "--inatividade-minutos", minutos)
			aguardar()
		case "4":
			runAssinatura("servidor", "parar")
			aguardar()
		}
	}
}

// ── menu: simulador ───────────────────────────────────────────────────────────

func menuSimulador() {
	for {
		clearScreen()
		printHeader()
		printSubtitle("Simulador HubSaúde")
		fmt.Println()

		repo := os.Getenv("RUNNER_SIMULADOR_REPO")
		jar := os.Getenv("RUNNER_SIMULADOR_JAR")
		if repo != "" {
			fmt.Printf("  %sRepositório:%s %s\n\n", gray, reset, repo)
		} else if jar != "" {
			fmt.Printf("  %sJAR local:%s %s\n\n", gray, reset, jar)
		} else {
			fmt.Printf("  %s⚠  Configure RUNNER_SIMULADOR_REPO=owner/repo%s\n", yellow, reset)
			fmt.Printf("  %s   ou RUNNER_SIMULADOR_JAR=/caminho/simulador.jar%s\n\n", yellow, reset)
		}

		printItem("1", "Status")
		printItem("2", "Iniciar  (porta 8080)")
		printItem("3", "Parar")
		fmt.Println()
		printItem("0", "Voltar")
		fmt.Println()

		switch prompt("Escolha") {
		case "0":
			return
		case "1":
			runSimulador("status")
			aguardar()
		case "2":
			runSimulador("iniciar")
			aguardar()
		case "3":
			runSimulador("parar")
			aguardar()
		}
	}
}

// ── execução dos CLIs ─────────────────────────────────────────────────────────

func runAssinatura(args ...string) (string, error) {
	return runBin(binAssinatura, map[string]string{"RUNNER_ASSINADOR_JAR": jarAssinador}, args...)
}

func runSimulador(args ...string) (string, error) {
	return runBin(binSimulador, nil, args...)
}

func runBin(bin string, extraEnv map[string]string, args ...string) (string, error) {
	if _, err := os.Stat(bin); err != nil {
		name := filepath.Base(bin)
		fmt.Printf("%sBinário não encontrado: %s%s\n", red, bin, reset)
		fmt.Printf("%sConstrua primeiro com: CGO_ENABLED=0 go build -o %s .%s\n\n", yellow, name, reset)
		return "", err
	}

	cmd := exec.Command(bin, args...)
	env := os.Environ()
	for k, v := range extraEnv {
		if v != "" {
			env = append(env, k+"="+v)
		}
	}
	cmd.Env = env
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	printCmd(bin, args)
	fmt.Println()

	err := cmd.Run()
	if err != nil {
		fmt.Printf("\n%sexit: %v%s\n", red, err, reset)
	}
	return "", err
}

func printCmd(bin string, args []string) {
	name := filepath.Base(bin)
	fmt.Printf("  %s$ %s%s %s%s\n", gray, bCyan, name, strings.Join(args, " "), reset)
	fmt.Printf("  %s%s%s\n\n", gray, strings.Repeat("─", 50), reset)
}

// ── utilitários visuais ───────────────────────────────────────────────────────

func clearScreen() {
	if runtime.GOOS == "windows" {
		exec.Command("cmd", "/c", "cls").Run()
	} else {
		fmt.Print("\033[H\033[2J")
	}
}

func printHeader() {
	fmt.Printf("%s%s", bold, cyan)
	fmt.Println("╔══════════════════════════════════════════════╗")
	fmt.Println("║         Sistema Runner — HubSaúde            ║")
	fmt.Println("║    Assinador Digital Simulado · Demo CLI      ║")
	fmt.Println("╚══════════════════════════════════════════════╝")
	fmt.Print(reset)
}

func printSubtitle(title string) {
	fmt.Printf("  %s%s▸ %s%s\n", bold, yellow, title, reset)
}

func printItem(key, label string) {
	fmt.Printf("  %s[%s]%s %s\n", cyan, key, reset, label)
}

func prompt(label string) string {
	fmt.Printf("\n%s%s › %s", bold, label, reset)
	line, _ := reader.ReadString('\n')
	return strings.TrimSpace(line)
}

func aguardar() {
	fmt.Printf("\n%sPressione Enter para continuar...%s", gray, reset)
	reader.ReadString('\n')
}

func verArquivos() {
	clearScreen()
	printHeader()
	printSubtitle("Arquivos de entrada (gerados automaticamente)")
	fmt.Println()

	files := []string{"bundle.json", "provenance.json", "material-pem.json", "material-pkcs12.json", "material-token.json"}
	for _, f := range files {
		path := filepath.Join(workDir, f)
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		fmt.Printf("  %s%s%s\n", yellow, path, reset)
		for _, line := range strings.Split(string(data), "\n") {
			fmt.Printf("  %s%s%s\n", gray, line, reset)
		}
		fmt.Println()
	}
	aguardar()
}
