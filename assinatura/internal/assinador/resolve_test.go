package assinador_test

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/thiagotelho/runner/assinatura/internal/assinador"
)

func TestFindAssinadorJar_EnvVar_ValidFile(t *testing.T) {
	dir := t.TempDir()
	jar := filepath.Join(dir, "assinador.jar")
	if err := os.WriteFile(jar, []byte("fake"), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Setenv(assinador.EnvVarAssinadorJar, jar)

	got, err := assinador.FindAssinadorJar()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	abs, _ := filepath.Abs(jar)
	if got != abs {
		t.Errorf("got %q, want %q", got, abs)
	}
}

func TestFindAssinadorJar_EnvVar_NotFound(t *testing.T) {
	t.Setenv(assinador.EnvVarAssinadorJar, "/nonexistent/assinador.jar")
	_, err := assinador.FindAssinadorJar()
	if err == nil {
		t.Fatal("expected error for non-existent path")
	}
}

func TestFindAssinadorJar_EnvVar_IsDirectory(t *testing.T) {
	t.Setenv(assinador.EnvVarAssinadorJar, t.TempDir())
	_, err := assinador.FindAssinadorJar()
	if err == nil {
		t.Fatal("expected error when env var points to directory")
	}
}

func TestFindJava_JAVA_HOME_Valid(t *testing.T) {
	home := t.TempDir()
	bin := filepath.Join(home, "bin")
	if err := os.MkdirAll(bin, 0o755); err != nil {
		t.Fatal(err)
	}
	name := "java"
	if runtime.GOOS == "windows" {
		name = "java.exe"
	}
	javaPath := filepath.Join(bin, name)
	if err := os.WriteFile(javaPath, []byte("fake"), 0o755); err != nil {
		t.Fatal(err)
	}

	t.Setenv("JAVA_HOME", home)
	got, err := assinador.FindJava()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	abs, _ := filepath.Abs(javaPath)
	if got != abs {
		t.Errorf("got %q, want %q", got, abs)
	}
}

func TestFindJava_JAVA_HOME_NoBinary(t *testing.T) {
	// JAVA_HOME exists but has no bin/java — falls through to PATH.
	// We just verify no panic; error is expected only if java is absent from PATH.
	t.Setenv("JAVA_HOME", t.TempDir())
	_, _ = assinador.FindJava()
}
