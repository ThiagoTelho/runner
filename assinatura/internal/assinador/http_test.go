package assinador_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/thiagotelho/runner/assinatura/internal/assinador"
)

// portFromURL extrai a porta de um httptest.Server.URL.
func portFromURL(t *testing.T, raw string) int {
	t.Helper()
	u, err := url.Parse(raw)
	if err != nil {
		t.Fatalf("parse url %q: %v", raw, err)
	}
	p, err := strconv.Atoi(u.Port())
	if err != nil {
		t.Fatalf("parse port %q: %v", u.Port(), err)
	}
	return p
}

func TestCheckHealth_OK(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/health" || r.Method != "GET" {
			http.Error(w, "bad", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	c := assinador.NewHTTPClient(portFromURL(t, srv.URL))
	if !c.CheckHealth(context.Background()) {
		t.Fatal("expected CheckHealth=true")
	}
}

func TestCheckHealth_Non200(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	c := assinador.NewHTTPClient(portFromURL(t, srv.URL))
	if c.CheckHealth(context.Background()) {
		t.Fatal("expected CheckHealth=false on 500")
	}
}

func TestCheckHealth_NoServer(t *testing.T) {
	// Random unused-ish port. We don't need to guarantee it's closed —
	// we just need *some* port that is overwhelmingly unlikely to respond.
	c := assinador.NewHTTPClient(1)
	if c.CheckHealth(context.Background()) {
		t.Fatal("expected CheckHealth=false when nothing is listening")
	}
}

func TestShutdown_Success(t *testing.T) {
	called := false
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/shutdown" && r.Method == "POST" {
			called = true
			w.WriteHeader(http.StatusOK)
			return
		}
		http.Error(w, "bad", http.StatusBadRequest)
	}))
	defer srv.Close()

	c := assinador.NewHTTPClient(portFromURL(t, srv.URL))
	if err := c.Shutdown(context.Background()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Fatal("server handler not invoked")
	}
}

func TestShutdown_NoServer(t *testing.T) {
	c := assinador.NewHTTPClient(1)
	if err := c.Shutdown(context.Background()); err == nil {
		t.Fatal("expected error when server unreachable")
	}
}

func TestRunCriarHTTP_Success(t *testing.T) {
	dir := t.TempDir()
	bundle := filepath.Join(dir, "bundle.json")
	prov := filepath.Join(dir, "prov.json")
	mat := filepath.Join(dir, "mat.json")
	mustWrite(t, bundle, `{"resourceType":"Bundle"}`)
	mustWrite(t, prov, `{"resourceType":"Provenance"}`)
	mustWrite(t, mat, `{"tipo":"PEM","chavePrivada":"key"}`)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/criar-assinatura" || r.Method != "POST" {
			http.Error(w, "bad", http.StatusBadRequest)
			return
		}
		body, _ := io.ReadAll(r.Body)
		var got map[string]json.RawMessage
		if err := json.Unmarshal(body, &got); err != nil {
			http.Error(w, "bad json", http.StatusBadRequest)
			return
		}
		// Sanity: the three top-level fields are present and non-empty.
		for _, k := range []string{"bundle", "provenance", "materialCriptografico"} {
			if len(got[k]) == 0 {
				http.Error(w, "missing "+k, http.StatusBadRequest)
				return
			}
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"jws":"a.b.c","outcome":{"resourceType":"OperationOutcome","issue":[]}}`))
	}))
	defer srv.Close()

	c := assinador.NewHTTPClient(portFromURL(t, srv.URL))
	var stdout, stderr bytes.Buffer
	err := c.RunCriarHTTP(context.Background(), bundle, prov, mat, &stdout, &stderr)
	if err != nil {
		t.Fatalf("unexpected error: %v\nstderr=%s", err, stderr.String())
	}
	if !strings.Contains(stdout.String(), "a.b.c") {
		t.Fatalf("stdout missing JWS: %s", stdout.String())
	}
}

func TestRunCriarHTTP_ServerReturns422(t *testing.T) {
	dir := t.TempDir()
	bundle := filepath.Join(dir, "bundle.json")
	prov := filepath.Join(dir, "prov.json")
	mat := filepath.Join(dir, "mat.json")
	mustWrite(t, bundle, `{}`)
	mustWrite(t, prov, `{}`)
	mustWrite(t, mat, `{}`)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(422)
		w.Write([]byte(`{"resourceType":"OperationOutcome","issue":[{"severity":"error"}]}`))
	}))
	defer srv.Close()

	c := assinador.NewHTTPClient(portFromURL(t, srv.URL))
	var stdout, stderr bytes.Buffer
	err := c.RunCriarHTTP(context.Background(), bundle, prov, mat, &stdout, &stderr)
	if err == nil {
		t.Fatal("expected error on 422 response")
	}
	if !strings.Contains(stderr.String(), "OperationOutcome") {
		t.Fatalf("expected OperationOutcome on stderr, got: %s", stderr.String())
	}
}

func TestRunCriarHTTP_FileMissing(t *testing.T) {
	c := assinador.NewHTTPClient(1)
	var stdout, stderr bytes.Buffer
	err := c.RunCriarHTTP(context.Background(),
		"/nonexistent/bundle.json",
		"/nonexistent/prov.json",
		"/nonexistent/mat.json",
		&stdout, &stderr)
	if err == nil {
		t.Fatal("expected error when input files do not exist")
	}
}

func TestRunValidarHTTP_ValidOutcome(t *testing.T) {
	dir := t.TempDir()
	jws := filepath.Join(dir, "sig.jws")
	mustWrite(t, jws, "a.b.c")

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/validar-assinatura" || r.Method != "POST" {
			http.Error(w, "bad", http.StatusBadRequest)
			return
		}
		w.Write([]byte(`{"resourceType":"OperationOutcome","issue":[{"severity":"information"}]}`))
	}))
	defer srv.Close()

	c := assinador.NewHTTPClient(portFromURL(t, srv.URL))
	var stdout, stderr bytes.Buffer
	err := c.RunValidarHTTP(context.Background(), assinador.ValidarParams{
		JwsPath:            jws,
		PoliticaRevogacao:  "warn",
		TimestampUnixUTC:   1700000000,
		PoliticaAssinatura: "https://example/policy",
	}, &stdout, &stderr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(stdout.String(), "OperationOutcome") {
		t.Fatalf("stdout missing OperationOutcome: %s", stdout.String())
	}
}

func TestRunValidarHTTP_ErrorSeverityFails(t *testing.T) {
	dir := t.TempDir()
	jws := filepath.Join(dir, "sig.jws")
	mustWrite(t, jws, "a.b.c")

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"resourceType":"OperationOutcome","issue":[{"severity":"error"}]}`))
	}))
	defer srv.Close()

	c := assinador.NewHTTPClient(portFromURL(t, srv.URL))
	var stdout, stderr bytes.Buffer
	err := c.RunValidarHTTP(context.Background(), assinador.ValidarParams{
		JwsPath:            jws,
		PoliticaRevogacao:  "warn",
		TimestampUnixUTC:   1700000000,
		PoliticaAssinatura: "https://example/policy",
	}, &stdout, &stderr)
	if err == nil {
		t.Fatal("expected error when outcome has severity=error")
	}
}

func TestRunValidarHTTP_TrimsJwsWhitespace(t *testing.T) {
	dir := t.TempDir()
	jws := filepath.Join(dir, "sig.jws")
	mustWrite(t, jws, "  a.b.c\n")

	var got string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var v map[string]any
		_ = json.Unmarshal(body, &v)
		got, _ = v["jws"].(string)
		w.Write([]byte(`{"issue":[]}`))
	}))
	defer srv.Close()

	c := assinador.NewHTTPClient(portFromURL(t, srv.URL))
	var stdout, stderr bytes.Buffer
	_ = c.RunValidarHTTP(context.Background(), assinador.ValidarParams{
		JwsPath:            jws,
		PoliticaRevogacao:  "warn",
		TimestampUnixUTC:   1700000000,
		PoliticaAssinatura: "p",
	}, &stdout, &stderr)
	if got != "a.b.c" {
		t.Fatalf("expected JWS to be trimmed, got %q", got)
	}
}

func TestRunValidarHTTP_OptionalBundleSent(t *testing.T) {
	dir := t.TempDir()
	jws := filepath.Join(dir, "sig.jws")
	bundle := filepath.Join(dir, "bundle.json")
	mustWrite(t, jws, "a.b.c")
	mustWrite(t, bundle, `{"resourceType":"Bundle"}`)

	var sentBundle string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var v map[string]any
		_ = json.Unmarshal(body, &v)
		sentBundle, _ = v["bundle"].(string)
		w.Write([]byte(`{"issue":[]}`))
	}))
	defer srv.Close()

	c := assinador.NewHTTPClient(portFromURL(t, srv.URL))
	var stdout, stderr bytes.Buffer
	if err := c.RunValidarHTTP(context.Background(), assinador.ValidarParams{
		JwsPath:            jws,
		PoliticaRevogacao:  "warn",
		TimestampUnixUTC:   1700000000,
		PoliticaAssinatura: "p",
		BundlePathOpcional: bundle,
	}, &stdout, &stderr); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(sentBundle, "Bundle") {
		t.Fatalf("expected bundle field to be forwarded, got %q", sentBundle)
	}
}

func mustWrite(t *testing.T, path, contents string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
		t.Fatal(err)
	}
}
