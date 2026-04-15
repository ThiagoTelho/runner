package assinador

import (
	"context"
	"fmt"
	"io"
	"os/exec"
)

// RunCriarLocal executa java -jar assinador.jar criar com os mesmos parâmetros do JAR.
func RunCriarLocal(ctx context.Context, javaExe, jarPath string, bundlePath, provenancePath, materialPath string, stdout, stderr io.Writer) error {
	args := []string{"-jar", jarPath, "criar",
		"-b", bundlePath,
		"-p", provenancePath,
		"-m", materialPath,
	}
	return runJava(ctx, javaExe, args, stdout, stderr)
}

// ValidarParams agrega os parâmetros do subcomando validar do assinador.jar.
type ValidarParams struct {
	JwsPath              string
	PoliticaRevogacao    string
	TimestampUnixUTC     int64
	PoliticaAssinatura   string
	BundlePathOpcional   string
}

// RunValidarLocal executa java -jar assinador.jar validar.
func RunValidarLocal(ctx context.Context, javaExe, jarPath string, p ValidarParams, stdout, stderr io.Writer) error {
	args := []string{"-jar", jarPath, "validar",
		"-j", p.JwsPath,
		"-r", p.PoliticaRevogacao,
		"-t", fmtInt64(p.TimestampUnixUTC),
		"-a", p.PoliticaAssinatura,
	}
	if p.BundlePathOpcional != "" {
		args = append(args, "-b", p.BundlePathOpcional)
	}
	return runJava(ctx, javaExe, args, stdout, stderr)
}

func fmtInt64(v int64) string {
	return fmt.Sprintf("%d", v)
}

func runJava(ctx context.Context, java string, args []string, stdout, stderr io.Writer) error {
	cmd := exec.CommandContext(ctx, java, args...)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	err := cmd.Run()
	if err == nil {
		return nil
	}
	if exitErr, ok := err.(*exec.ExitError); ok {
		return fmt.Errorf("assinador.jar encerrou com código %d", exitErr.ExitCode())
	}
	return fmt.Errorf("executar java: %w", err)
}
