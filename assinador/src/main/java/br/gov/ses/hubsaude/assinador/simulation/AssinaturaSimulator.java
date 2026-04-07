package br.gov.ses.hubsaude.assinador.simulation;

import br.gov.ses.hubsaude.assinador.model.CriarAssinaturaRequest;
import br.gov.ses.hubsaude.assinador.model.OperationOutcome;
import br.gov.ses.hubsaude.assinador.validation.CriarAssinaturaValidator;
import br.gov.ses.hubsaude.assinador.validation.ValidationException;
import com.google.gson.Gson;
import com.google.gson.JsonObject;

import java.nio.charset.StandardCharsets;
import java.time.Instant;
import java.util.Base64;
import java.util.UUID;

/**
 * Simula a criação de uma assinatura digital, retornando um JWS
 * pré-construído quando os parâmetros são válidos.
 */
public class AssinaturaSimulator {

    private static final Gson GSON = new Gson();
    private final CriarAssinaturaValidator validator = new CriarAssinaturaValidator();

    /**
     * Valida parâmetros e, se tudo estiver correto, retorna um
     * JWS simulado (header.payload.signature) + OperationOutcome de sucesso.
     */
    public SimulationResult criar(CriarAssinaturaRequest request) throws ValidationException {
        validator.validate(request);

        String jws = buildSimulatedJws(request);
        OperationOutcome outcome = OperationOutcome.success("Assinatura simulada criada com sucesso.");
        return new SimulationResult(jws, outcome);
    }

    private String buildSimulatedJws(CriarAssinaturaRequest request) {
        var header = new JsonObject();
        header.addProperty("alg", "RS256");
        header.addProperty("typ", "JWS");
        header.addProperty("kid", UUID.randomUUID().toString());

        var payload = new JsonObject();
        payload.addProperty("iss", "runner-assinador-simulado");
        payload.addProperty("iat", Instant.now().getEpochSecond());
        payload.addProperty("jti", UUID.randomUUID().toString());
        payload.addProperty("sub", "simulacao");

        String headerB64  = b64(GSON.toJson(header));
        String payloadB64 = b64(GSON.toJson(payload));
        String signatureB64 = b64("SIMULATED_SIGNATURE_" + UUID.randomUUID());

        return headerB64 + "." + payloadB64 + "." + signatureB64;
    }

    private static String b64(String input) {
        return Base64.getUrlEncoder().withoutPadding()
                .encodeToString(input.getBytes(StandardCharsets.UTF_8));
    }

    public record SimulationResult(String jws, OperationOutcome outcome) {}
}
