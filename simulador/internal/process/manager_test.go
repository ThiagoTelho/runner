package process

import (
	"context"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// redirectCacheDir aponta os.UserCacheDir() para um diretório temporário, para que
// os testes não toquem no cache real do usuário.
func redirectCacheDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	switch runtime.GOOS {
	case "windows":
		t.Setenv("LocalAppData", dir)
	case "darwin":
		t.Setenv("HOME", dir)
	default:
		t.Setenv("XDG_CACHE_HOME", dir)
	}
	return dir
}

// freePort allocates a TCP port and returns it (the listener is closed before return).
// There is a small race window, but it is fine for tests.
func freePort(t *testing.T) int {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	port := ln.Addr().(*net.TCPAddr).Port
	_ = ln.Close()
	return port
}

// --- CacheDir ---

func TestCacheDir_PointsToRunnerSimulador(t *testing.T) {
	redirectCacheDir(t)
	dir, err := CacheDir()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasSuffix(filepath.ToSlash(dir), "runner/simulador") {
		t.Errorf("expected path ending in runner/simulador, got %q", dir)
	}
}

// --- CheckPort ---

func TestCheckPort_FreePort(t *testing.T) {
	port := freePort(t)
	if err := CheckPort(port); err != nil {
		t.Fatalf("expected nil for free port, got %v", err)
	}
}

func TestCheckPort_InUse(t *testing.T) {
	// Bind no mesmo endereço wildcard que CheckPort usa (":port"),
	// para garantir conflito real em qualquer plataforma.
	ln, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()
	port := ln.Addr().(*net.TCPAddr).Port

	if err := CheckPort(port); err == nil {
		t.Fatalf("expected error for port %d in use", port)
	}
}

// --- Status / Stop com PID file ---

func writePIDFile(t *testing.T, pid, port int) {
	t.Helper()
	dir, err := CacheDir()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	content := fmt.Sprintf("%d\n%d\n", pid, port)
	if err := os.WriteFile(filepath.Join(dir, "simulador.pid"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestStatus_NoPIDFile(t *testing.T) {
	redirectCacheDir(t)
	running, pid, port := Status()
	if running || pid != 0 || port != 0 {
		t.Errorf("expected (false, 0, 0); got (%v, %d, %d)", running, pid, port)
	}
}

func TestStatus_RunningProcess(t *testing.T) {
	redirectCacheDir(t)
	// Próprio processo de teste — sempre vivo.
	writePIDFile(t, os.Getpid(), 8080)

	running, pid, port := Status()
	if !running {
		t.Fatal("expected running=true for current process")
	}
	if pid != os.Getpid() {
		t.Errorf("pid: want %d, got %d", os.Getpid(), pid)
	}
	if port != 8080 {
		t.Errorf("port: want 8080, got %d", port)
	}
}

func TestStatus_StalePIDFileIsCleanedUp(t *testing.T) {
	redirectCacheDir(t)
	// PID extremamente improvável de existir como processo ativo.
	writePIDFile(t, 99999999, 8080)

	running, _, _ := Status()
	if running {
		t.Fatal("expected running=false for nonexistent PID")
	}

	// O PID file deve ter sido removido.
	dir, _ := CacheDir()
	if _, err := os.Stat(filepath.Join(dir, "simulador.pid")); !os.IsNotExist(err) {
		t.Errorf("expected stale PID file to be removed; stat err=%v", err)
	}
}

func TestStatus_MalformedPIDFile(t *testing.T) {
	redirectCacheDir(t)
	dir, _ := CacheDir()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "simulador.pid"), []byte("garbage"), 0o644); err != nil {
		t.Fatal(err)
	}

	running, _, _ := Status()
	if running {
		t.Error("expected running=false for malformed PID file")
	}
}

func TestStop_NoPIDFile(t *testing.T) {
	redirectCacheDir(t)
	if err := Stop(); err == nil {
		t.Fatal("expected error when PID file missing")
	}
}

func TestStop_StalePIDFileRemoved(t *testing.T) {
	redirectCacheDir(t)
	writePIDFile(t, 99999999, 8080)

	if err := Stop(); err == nil {
		t.Fatal("expected error when PID does not point to a running process")
	}

	dir, _ := CacheDir()
	if _, err := os.Stat(filepath.Join(dir, "simulador.pid")); !os.IsNotExist(err) {
		t.Errorf("expected stale PID file to be removed; stat err=%v", err)
	}
}

// --- FindOrDownload (somente caminhos sem rede) ---

func TestFindOrDownload_EnvJarValid(t *testing.T) {
	redirectCacheDir(t)
	jar := filepath.Join(t.TempDir(), "simulador.jar")
	if err := os.WriteFile(jar, []byte("fake"), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Setenv(EnvSimuladorJar, jar)
	t.Setenv(EnvSimuladorRepo, "")

	got, err := FindOrDownload(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != jar {
		t.Errorf("got %q, want %q", got, jar)
	}
}

func TestFindOrDownload_EnvJarMissing(t *testing.T) {
	redirectCacheDir(t)
	t.Setenv(EnvSimuladorJar, "/nonexistent/simulador.jar")

	if _, err := FindOrDownload(context.Background()); err == nil {
		t.Fatal("expected error when env var points to missing file")
	}
}

func TestFindOrDownload_NoEnvVarsConfigured(t *testing.T) {
	redirectCacheDir(t)
	t.Setenv(EnvSimuladorJar, "")
	t.Setenv(EnvSimuladorRepo, "")

	_, err := FindOrDownload(context.Background())
	if err == nil {
		t.Fatal("expected error when neither env var is configured")
	}
	if !strings.Contains(err.Error(), EnvSimuladorRepo) {
		t.Errorf("error should mention %s, got: %v", EnvSimuladorRepo, err)
	}
}

// --- helpers internos: PID file roundtrip ---

func TestReadPIDFile_RoundTrip(t *testing.T) {
	redirectCacheDir(t)
	writePIDFile(t, 1234, 5678)

	pid, port, err := readPIDFile()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pid != 1234 || port != 5678 {
		t.Errorf("got (pid=%d, port=%d); want (1234, 5678)", pid, port)
	}
}

func TestReadPIDFile_Truncated(t *testing.T) {
	redirectCacheDir(t)
	dir, _ := CacheDir()
	_ = os.MkdirAll(dir, 0o755)
	if err := os.WriteFile(filepath.Join(dir, "simulador.pid"), []byte("1234"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, _, err := readPIDFile(); err == nil {
		t.Fatal("expected error when PID file lacks port line")
	}
}
