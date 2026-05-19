package assinador

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/thiagotelho/runner/jdk"
)

// EnvVarAssinadorJar é a variável de ambiente que, se definida, aponta para o arquivo assinador.jar.
// Tem precedência sobre a busca ao lado do executável do CLI.
const EnvVarAssinadorJar = "RUNNER_ASSINADOR_JAR"

// FindAssinadorJar localiza o JAR do Assinador:
// 1) RUNNER_ASSINADOR_JAR, se definido e válido;
// 2) arquivo "assinador.jar" no mesmo diretório do binário assinatura.
func FindAssinadorJar() (string, error) {
	if p := os.Getenv(EnvVarAssinadorJar); p != "" {
		abs, err := filepath.Abs(p)
		if err != nil {
			return "", fmt.Errorf("%s: %w", EnvVarAssinadorJar, err)
		}
		fi, err := os.Stat(abs)
		if err != nil {
			return "", fmt.Errorf("%s: caminho inválido %q: %w", EnvVarAssinadorJar, abs, err)
		}
		if fi.IsDir() {
			return "", fmt.Errorf("%s não pode ser um diretório: %s", EnvVarAssinadorJar, abs)
		}
		return abs, nil
	}

	exe, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("obter caminho do executável: %w", err)
	}
	exe, err = filepath.EvalSymlinks(exe)
	if err != nil {
		return "", fmt.Errorf("resolver symlinks do executável: %w", err)
	}
	dir := filepath.Dir(exe)
	candidate := filepath.Join(dir, "assinador.jar")
	fi, err := os.Stat(candidate)
	if err != nil || fi.IsDir() {
		return "", fmt.Errorf(
			`não foi possível localizar assinador.jar ao lado do executável (%s). Defina %s ou copie o JAR para %s`,
			exe, EnvVarAssinadorJar, candidate)
	}
	abs, err := filepath.Abs(candidate)
	if err != nil {
		return "", err
	}
	return abs, nil
}

// FindOrProvisionJava retorna um executável java com versão >= 21.
// Se não encontrar no sistema, baixa e faz cache do Temurin JDK 21 automaticamente.
func FindOrProvisionJava(ctx context.Context) (string, error) {
	return jdk.FindOrProvision(ctx)
}

// FindJava resolve o executável `java`: JAVA_HOME/bin/java quando existir; senão `java` no PATH.
func FindJava() (string, error) {
	if jh := os.Getenv("JAVA_HOME"); jh != "" {
		if p, err := javaInsideJAVA_HOME(jh); err == nil {
			return p, nil
		}
	}
	path, err := exec.LookPath("java")
	if err != nil {
		return "", fmt.Errorf("java não encontrado (configure JAVA_HOME ou PATH): %w", err)
	}
	return path, nil
}

func javaInsideJAVA_HOME(home string) (string, error) {
	bin := filepath.Join(home, "bin")
	candidates := []string{
		filepath.Join(bin, "java"),
		filepath.Join(bin, "java.exe"),
	}
	if runtime.GOOS == "windows" {
		candidates = []string{
			filepath.Join(bin, "java.exe"),
			filepath.Join(bin, "java"),
		}
	}
	for _, p := range candidates {
		fi, err := os.Stat(p)
		if err != nil || fi.IsDir() {
			continue
		}
		return filepath.Abs(p)
	}
	return "", fmt.Errorf("executável java não encontrado em JAVA_HOME=%s", home)
}
