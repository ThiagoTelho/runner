package assinador

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const DefaultPort = 8190

// HTTPClient é o cliente HTTP para o assinador.jar em modo servidor.
type HTTPClient struct {
	baseURL string
	client  *http.Client
}

func NewHTTPClient(port int) *HTTPClient {
	return &HTTPClient{
		baseURL: fmt.Sprintf("http://localhost:%d", port),
		client:  &http.Client{Timeout: 30 * time.Second},
	}
}

// CheckHealth retorna true se o servidor estiver respondendo em /health.
func (c *HTTPClient) CheckHealth(ctx context.Context) bool {
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/health", nil)
	if err != nil {
		return false
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return false
	}
	resp.Body.Close()
	return resp.StatusCode == 200
}

// Shutdown envia POST /shutdown ao servidor.
func (c *HTTPClient) Shutdown(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/shutdown", nil)
	if err != nil {
		return err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("assinador não está em execução ou não respondeu: %w", err)
	}
	resp.Body.Close()
	return nil
}

// RunCriarHTTP envia os arquivos ao endpoint /criar-assinatura e exibe o resultado.
func (c *HTTPClient) RunCriarHTTP(ctx context.Context, bundlePath, provenancePath, materialPath string, stdout, stderr io.Writer) error {
	bundle, err := os.ReadFile(bundlePath)
	if err != nil {
		return fmt.Errorf("ler bundle: %w", err)
	}
	provenance, err := os.ReadFile(provenancePath)
	if err != nil {
		return fmt.Errorf("ler provenance: %w", err)
	}
	materialRaw, err := os.ReadFile(materialPath)
	if err != nil {
		return fmt.Errorf("ler material: %w", err)
	}

	type reqBody struct {
		Bundle               string          `json:"bundle"`
		Provenance           string          `json:"provenance"`
		MaterialCriptografico json.RawMessage `json:"materialCriptografico"`
	}
	payload, _ := json.Marshal(reqBody{
		Bundle:               string(bundle),
		Provenance:           string(provenance),
		MaterialCriptografico: json.RawMessage(materialRaw),
	})

	respBody, status, err := c.post(ctx, "/criar-assinatura", payload)
	if err != nil {
		return err
	}
	if status == 422 || status >= 400 {
		fmt.Fprintf(stderr, "%s\n", prettyJSON(respBody))
		return fmt.Errorf("assinador retornou status %d", status)
	}

	var result struct {
		JWS     string          `json:"jws"`
		Outcome json.RawMessage `json:"outcome"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		fmt.Fprintf(stdout, "%s\n", respBody)
		return nil
	}
	fmt.Fprintf(stdout, "JWS gerado:\n%s\n\n%s\n", result.JWS, prettyJSON(result.Outcome))
	return nil
}

// RunValidarHTTP envia os parâmetros ao endpoint /validar-assinatura e exibe o resultado.
func (c *HTTPClient) RunValidarHTTP(ctx context.Context, p ValidarParams, stdout, stderr io.Writer) error {
	jwsBytes, err := os.ReadFile(p.JwsPath)
	if err != nil {
		return fmt.Errorf("ler JWS: %w", err)
	}

	type reqBody struct {
		JWS                 string `json:"jws"`
		PoliticaRevogacao   string `json:"politicaRevogacao"`
		TimestampReferencia int64  `json:"timestampReferencia"`
		PoliticaAssinatura  string `json:"politicaAssinatura"`
		Bundle              string `json:"bundle,omitempty"`
	}
	body := reqBody{
		JWS:                 strings.TrimSpace(string(jwsBytes)),
		PoliticaRevogacao:   p.PoliticaRevogacao,
		TimestampReferencia: p.TimestampUnixUTC,
		PoliticaAssinatura:  p.PoliticaAssinatura,
	}
	if p.BundlePathOpcional != "" {
		b, err := os.ReadFile(p.BundlePathOpcional)
		if err != nil {
			return fmt.Errorf("ler bundle: %w", err)
		}
		body.Bundle = string(b)
	}

	payload, _ := json.Marshal(body)
	respBody, status, err := c.post(ctx, "/validar-assinatura", payload)
	if err != nil {
		return err
	}
	if status == 422 {
		fmt.Fprintf(stderr, "%s\n", prettyJSON(respBody))
		return fmt.Errorf("assinador retornou status %d", status)
	}

	fmt.Fprintf(stdout, "%s\n", prettyJSON(respBody))

	var outcome struct {
		Issue []struct {
			Severity string `json:"severity"`
		} `json:"issue"`
	}
	if err := json.Unmarshal(respBody, &outcome); err == nil {
		for _, iss := range outcome.Issue {
			if iss.Severity == "error" {
				return fmt.Errorf("assinatura inválida")
			}
		}
	}
	return nil
}

func (c *HTTPClient) post(ctx context.Context, path string, body []byte) ([]byte, int, error) {
	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+path, bytes.NewReader(body))
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("chamar assinador HTTP (%s): %w", path, err)
	}
	defer resp.Body.Close()
	respBody, _ := io.ReadAll(resp.Body)
	return respBody, resp.StatusCode, nil
}

func prettyJSON(data []byte) string {
	var buf bytes.Buffer
	if err := json.Indent(&buf, data, "", "  "); err != nil {
		return string(data)
	}
	return buf.String()
}
