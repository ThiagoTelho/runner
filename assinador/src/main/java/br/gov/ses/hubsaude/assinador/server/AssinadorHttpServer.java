package br.gov.ses.hubsaude.assinador.server;

import br.gov.ses.hubsaude.assinador.model.*;
import br.gov.ses.hubsaude.assinador.simulation.AssinaturaSimulator;
import br.gov.ses.hubsaude.assinador.simulation.ValidacaoSimulator;
import br.gov.ses.hubsaude.assinador.validation.ValidationException;
import com.google.gson.Gson;
import com.google.gson.JsonObject;
import com.google.gson.JsonSyntaxException;
import com.sun.net.httpserver.HttpExchange;
import com.sun.net.httpserver.HttpServer;

import java.io.IOException;
import java.io.OutputStream;
import java.net.InetSocketAddress;
import java.nio.charset.StandardCharsets;
import java.util.concurrent.Executors;

/**
 * Modo servidor HTTP do Assinador (warm start).
 * Endpoints:
 * <ul>
 *   <li>POST /criar-assinatura</li>
 *   <li>POST /validar-assinatura</li>
 *   <li>POST /shutdown</li>
 *   <li>GET  /health</li>
 * </ul>
 */
public class AssinadorHttpServer {

    public static final int DEFAULT_PORT = 8190;

    private final Gson gson = new Gson();
    private final AssinaturaSimulator criarSimulator = new AssinaturaSimulator();
    private final ValidacaoSimulator validarSimulator = new ValidacaoSimulator();
    private HttpServer server;
    private final java.util.concurrent.atomic.AtomicLong lastActivity =
            new java.util.concurrent.atomic.AtomicLong(System.currentTimeMillis());
    private java.util.concurrent.ScheduledExecutorService watchdog;

    public void start(int port) throws IOException {
        start(port, 0);
    }

    public void start(int port, int inativadeMinutos) throws IOException {
        server = HttpServer.create(new InetSocketAddress(port), 0);
        server.setExecutor(Executors.newVirtualThreadPerTaskExecutor());

        server.createContext("/criar-assinatura", this::handleCriar);
        server.createContext("/validar-assinatura", this::handleValidar);
        server.createContext("/shutdown", this::handleShutdown);
        server.createContext("/health", this::handleHealth);

        server.start();
        System.out.printf("Assinador HTTP iniciado na porta %d%n", port);

        if (inativadeMinutos > 0) {
            long intervalMs = 30_000L;
            long limiteMs = inativadeMinutos * 60_000L;
            watchdog = Executors.newSingleThreadScheduledExecutor(r -> {
                Thread t = new Thread(r, "idle-watchdog");
                t.setDaemon(true);
                return t;
            });
            watchdog.scheduleAtFixedRate(() -> {
                if (System.currentTimeMillis() - lastActivity.get() >= limiteMs) {
                    System.out.printf("Encerrando por inatividade (%d min).%n", inativadeMinutos);
                    stop();
                    System.exit(0);
                }
            }, intervalMs, intervalMs, java.util.concurrent.TimeUnit.MILLISECONDS);
        }
    }

    public void stop() {
        if (watchdog != null) {
            watchdog.shutdownNow();
        }
        if (server != null) {
            server.stop(1);
            System.out.println("Assinador HTTP encerrado.");
        }
    }

    // --- Handlers ---

    private void handleCriar(HttpExchange exchange) throws IOException {
        if (!requirePost(exchange)) return;
        lastActivity.set(System.currentTimeMillis());

        String body = readBody(exchange);
        try {
            var request = gson.fromJson(body, CriarAssinaturaRequest.class);
            if (request == null) {
                sendJson(exchange, 400, OperationOutcome.error("Corpo da requisição vazio ou inválido.", null));
                return;
            }
            var result = criarSimulator.criar(request);
            var response = new JsonObject();
            response.addProperty("jws", result.jws());
            response.add("outcome", gson.toJsonTree(result.outcome()));
            sendJson(exchange, 200, response);
        } catch (ValidationException e) {
            sendJson(exchange, 422, e.getOutcome());
        } catch (JsonSyntaxException e) {
            sendJson(exchange, 400, OperationOutcome.error("JSON inválido: " + e.getMessage(), null));
        }
    }

    private void handleValidar(HttpExchange exchange) throws IOException {
        if (!requirePost(exchange)) return;
        lastActivity.set(System.currentTimeMillis());

        String body = readBody(exchange);
        try {
            var request = gson.fromJson(body, ValidarAssinaturaRequest.class);
            if (request == null) {
                sendJson(exchange, 400, OperationOutcome.error("Corpo da requisição vazio ou inválido.", null));
                return;
            }
            var outcome = validarSimulator.validar(request);
            int status = outcome.hasErrors() ? 200 : 200;
            sendJson(exchange, status, outcome);
        } catch (ValidationException e) {
            sendJson(exchange, 422, e.getOutcome());
        } catch (JsonSyntaxException e) {
            sendJson(exchange, 400, OperationOutcome.error("JSON inválido: " + e.getMessage(), null));
        }
    }

    private void handleShutdown(HttpExchange exchange) throws IOException {
        if (!requirePost(exchange)) return;
        sendJson(exchange, 200, OperationOutcome.success("Servidor encerrando."));
        new Thread(() -> {
            try { Thread.sleep(200); } catch (InterruptedException ignored) {}
            stop();
        }).start();
    }

    private void handleHealth(HttpExchange exchange) throws IOException {
        if (!"GET".equalsIgnoreCase(exchange.getRequestMethod())) {
            sendJson(exchange, 405, OperationOutcome.error("Método não permitido. Use GET.", null));
            return;
        }
        sendJson(exchange, 200, OperationOutcome.success("OK"));
    }

    // --- Helpers ---

    private boolean requirePost(HttpExchange exchange) throws IOException {
        if (!"POST".equalsIgnoreCase(exchange.getRequestMethod())) {
            sendJson(exchange, 405, OperationOutcome.error("Método não permitido. Use POST.", null));
            return false;
        }
        return true;
    }

    private String readBody(HttpExchange exchange) throws IOException {
        return new String(exchange.getRequestBody().readAllBytes(), StandardCharsets.UTF_8);
    }

    private void sendJson(HttpExchange exchange, int status, Object body) throws IOException {
        byte[] bytes = gson.toJson(body).getBytes(StandardCharsets.UTF_8);
        exchange.getResponseHeaders().set("Content-Type", "application/json; charset=utf-8");
        exchange.sendResponseHeaders(status, bytes.length);
        try (OutputStream os = exchange.getResponseBody()) {
            os.write(bytes);
        }
    }
}
