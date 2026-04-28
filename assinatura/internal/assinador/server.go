package assinador

import (
	"context"
	"fmt"
	"os/exec"
	"time"
)

// StartServer inicia o assinador.jar em modo servidor como processo em background.
// inativadeMinutos=0 significa sem timeout de inatividade.
func StartServer(javaExe, jarPath string, port, inativadeMinutos int) error {
	args := []string{"-jar", jarPath, "servidor", "--porta", fmt.Sprintf("%d", port)}
	if inativadeMinutos > 0 {
		args = append(args, "--inatividade-minutos", fmt.Sprintf("%d", inativadeMinutos))
	}
	cmd := exec.Command(javaExe, args...)
	detachProcess(cmd)
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("iniciar assinador servidor: %w", err)
	}
	// Libera o processo para rodar independente do CLI.
	_ = cmd.Process.Release()
	return nil
}

// WaitForServer aguarda o servidor responder em /health até o timeout expirar.
func WaitForServer(ctx context.Context, client *HTTPClient, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if client.CheckHealth(ctx) {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(250 * time.Millisecond):
		}
	}
	return fmt.Errorf("assinador servidor não respondeu em %v", timeout)
}
