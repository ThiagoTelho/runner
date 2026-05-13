package br.gov.ses.hubsaude.assinador.server;

import com.google.gson.Gson;
import com.google.gson.JsonObject;
import org.junit.jupiter.api.AfterEach;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;

import java.io.IOException;
import java.net.ServerSocket;
import java.net.URI;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;
import java.time.Duration;
import java.time.Instant;

import static org.junit.jupiter.api.Assertions.*;

class AssinadorHttpServerTest {

    private AssinadorHttpServer server;
    private int port;
    private HttpClient http;
    private final Gson gson = new Gson();

    @BeforeEach
    void setUp() throws IOException {
        port = freePort();
        server = new AssinadorHttpServer();
        server.start(port);
        http = HttpClient.newBuilder().connectTimeout(Duration.ofSeconds(5)).build();
    }

    @AfterEach
    void tearDown() {
        if (server != null) server.stop();
    }

    @Test
    void healthReturnsOk() throws Exception {
        var resp = http.send(
                HttpRequest.newBuilder(URI.create(base() + "/health")).GET().build(),
                HttpResponse.BodyHandlers.ofString());
        assertEquals(200, resp.statusCode());
        assertTrue(resp.body().contains("OperationOutcome"));
    }

    @Test
    void healthRejectsPost() throws Exception {
        var resp = http.send(
                HttpRequest.newBuilder(URI.create(base() + "/health"))
                        .POST(HttpRequest.BodyPublishers.noBody())
                        .build(),
                HttpResponse.BodyHandlers.ofString());
        assertEquals(405, resp.statusCode());
    }

    @Test
    void criarRejectsGet() throws Exception {
        var resp = http.send(
                HttpRequest.newBuilder(URI.create(base() + "/criar-assinatura")).GET().build(),
                HttpResponse.BodyHandlers.ofString());
        assertEquals(405, resp.statusCode());
    }

    @Test
    void criarReturnsBadRequestOnEmptyBody() throws Exception {
        var resp = http.send(
                HttpRequest.newBuilder(URI.create(base() + "/criar-assinatura"))
                        .header("Content-Type", "application/json")
                        .POST(HttpRequest.BodyPublishers.ofString(""))
                        .build(),
                HttpResponse.BodyHandlers.ofString());
        assertEquals(400, resp.statusCode());
    }

    @Test
    void criarReturnsBadRequestOnInvalidJson() throws Exception {
        var resp = http.send(
                HttpRequest.newBuilder(URI.create(base() + "/criar-assinatura"))
                        .header("Content-Type", "application/json")
                        .POST(HttpRequest.BodyPublishers.ofString("{invalid"))
                        .build(),
                HttpResponse.BodyHandlers.ofString());
        assertEquals(400, resp.statusCode());
        assertTrue(resp.body().contains("error"));
    }

    @Test
    void criarReturns422OnValidationFailure() throws Exception {
        // Empty object — missing all required fields.
        var resp = http.send(
                HttpRequest.newBuilder(URI.create(base() + "/criar-assinatura"))
                        .header("Content-Type", "application/json")
                        .POST(HttpRequest.BodyPublishers.ofString("{}"))
                        .build(),
                HttpResponse.BodyHandlers.ofString());
        assertEquals(422, resp.statusCode());
        assertTrue(resp.body().contains("error"));
    }

    @Test
    void criarReturnsJwsOnValidPemRequest() throws Exception {
        String body = """
                {
                  "bundle":"{\\"resourceType\\":\\"Bundle\\"}",
                  "provenance":"{\\"resourceType\\":\\"Provenance\\"}",
                  "materialCriptografico":{
                    "tipo":"PEM",
                    "chavePrivada":"-----BEGIN PRIVATE KEY-----\\nfake\\n-----END PRIVATE KEY-----"
                  }
                }""";
        var resp = http.send(
                HttpRequest.newBuilder(URI.create(base() + "/criar-assinatura"))
                        .header("Content-Type", "application/json")
                        .POST(HttpRequest.BodyPublishers.ofString(body))
                        .build(),
                HttpResponse.BodyHandlers.ofString());
        assertEquals(200, resp.statusCode());
        JsonObject obj = gson.fromJson(resp.body(), JsonObject.class);
        assertNotNull(obj.get("jws"));
        String jws = obj.get("jws").getAsString();
        assertEquals(3, jws.split("\\.").length);
        assertNotNull(obj.get("outcome"));
    }

    @Test
    void validarRoundTripsSimulatedSignature() throws Exception {
        // First call /criar-assinatura to get a valid JWS.
        String criarBody = """
                {
                  "bundle":"{\\"resourceType\\":\\"Bundle\\"}",
                  "provenance":"{\\"resourceType\\":\\"Provenance\\"}",
                  "materialCriptografico":{
                    "tipo":"PEM",
                    "chavePrivada":"-----BEGIN PRIVATE KEY-----\\nfake\\n-----END PRIVATE KEY-----"
                  }
                }""";
        var criarResp = http.send(
                HttpRequest.newBuilder(URI.create(base() + "/criar-assinatura"))
                        .header("Content-Type", "application/json")
                        .POST(HttpRequest.BodyPublishers.ofString(criarBody))
                        .build(),
                HttpResponse.BodyHandlers.ofString());
        assertEquals(200, criarResp.statusCode());
        String jws = gson.fromJson(criarResp.body(), JsonObject.class).get("jws").getAsString();

        String validarBody = String.format("""
                {
                  "jws":"%s",
                  "politicaRevogacao":"warn",
                  "timestampReferencia":%d,
                  "politicaAssinatura":"https://fhir.saude.go.gov.br/r4/seguranca/assinatura/v1"
                }""", jws, Instant.now().getEpochSecond());

        var validarResp = http.send(
                HttpRequest.newBuilder(URI.create(base() + "/validar-assinatura"))
                        .header("Content-Type", "application/json")
                        .POST(HttpRequest.BodyPublishers.ofString(validarBody))
                        .build(),
                HttpResponse.BodyHandlers.ofString());
        assertEquals(200, validarResp.statusCode());
        // Outcome should not contain "error" severity for a known simulated JWS.
        assertFalse(validarResp.body().contains("\"severity\":\"error\""),
                "Expected no error severity in: " + validarResp.body());
    }

    @Test
    void validarReturns422WhenRequestInvalid() throws Exception {
        var resp = http.send(
                HttpRequest.newBuilder(URI.create(base() + "/validar-assinatura"))
                        .header("Content-Type", "application/json")
                        .POST(HttpRequest.BodyPublishers.ofString("{}"))
                        .build(),
                HttpResponse.BodyHandlers.ofString());
        assertEquals(422, resp.statusCode());
    }

    @Test
    void shutdownReturns200AndStopsServer() throws Exception {
        var resp = http.send(
                HttpRequest.newBuilder(URI.create(base() + "/shutdown"))
                        .POST(HttpRequest.BodyPublishers.noBody()).build(),
                HttpResponse.BodyHandlers.ofString());
        assertEquals(200, resp.statusCode());

        // The handler sleeps 200ms then calls server.stop(1) which can wait up to 1s
        // for in-flight requests to drain. Poll until the port is no longer accepting.
        long deadline = System.currentTimeMillis() + 5_000;
        boolean down = false;
        while (System.currentTimeMillis() < deadline) {
            try (java.net.Socket s = new java.net.Socket()) {
                s.connect(new java.net.InetSocketAddress("localhost", port), 200);
                Thread.sleep(100);
            } catch (IOException expected) {
                down = true;
                break;
            }
        }
        assertTrue(down, "server still accepting connections after shutdown");
    }

    // --- helpers ---

    private String base() {
        return "http://localhost:" + port;
    }

    private static int freePort() throws IOException {
        try (ServerSocket s = new ServerSocket(0)) {
            return s.getLocalPort();
        }
    }
}
