package assinador

import (
	"context"
	"fmt"
	"io"
	"time"
)

// RunCriar executa o comando criar em modo local (--local) ou HTTP (padrão).
// Em modo HTTP: reutiliza instância em execução ou inicia uma nova.
func RunCriar(ctx context.Context, javaExe, jarPath string, port int, local bool,
	bundlePath, provenancePath, materialPath string, stdout, stderr io.Writer) error {
	if local {
		return RunCriarLocal(ctx, javaExe, jarPath, bundlePath, provenancePath, materialPath, stdout, stderr)
	}
	client := NewHTTPClient(port)
	if err := ensureServer(ctx, client, javaExe, jarPath, port, stderr); err != nil {
		return err
	}
	return client.RunCriarHTTP(ctx, bundlePath, provenancePath, materialPath, stdout, stderr)
}

// RunValidar executa o comando validar em modo local (--local) ou HTTP (padrão).
func RunValidar(ctx context.Context, javaExe, jarPath string, port int, local bool,
	p ValidarParams, stdout, stderr io.Writer) error {
	if local {
		return RunValidarLocal(ctx, javaExe, jarPath, p, stdout, stderr)
	}
	client := NewHTTPClient(port)
	if err := ensureServer(ctx, client, javaExe, jarPath, port, stderr); err != nil {
		return err
	}
	return client.RunValidarHTTP(ctx, p, stdout, stderr)
}

// ensureServer garante que um servidor esteja em execução na porta indicada.
// Se não houver instância, inicia uma nova e aguarda até 15 segundos.
func ensureServer(ctx context.Context, client *HTTPClient, javaExe, jarPath string, port int, stderr io.Writer) error {
	if client.CheckHealth(ctx) {
		return nil
	}
	fmt.Fprintf(stderr, "Iniciando assinador em modo servidor (porta %d)...\n", port)
	if err := StartServer(javaExe, jarPath, port, 0); err != nil {
		return err
	}
	if err := WaitForServer(ctx, client, 15*time.Second); err != nil {
		return fmt.Errorf("%w\nDica: use --local para invocação direta sem servidor", err)
	}
	return nil
}
