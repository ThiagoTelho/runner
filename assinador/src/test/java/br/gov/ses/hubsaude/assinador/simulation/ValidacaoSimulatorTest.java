package br.gov.ses.hubsaude.assinador.simulation;

import br.gov.ses.hubsaude.assinador.model.CriarAssinaturaRequest;
import br.gov.ses.hubsaude.assinador.model.MaterialCriptografico;
import br.gov.ses.hubsaude.assinador.model.ValidarAssinaturaRequest;
import br.gov.ses.hubsaude.assinador.validation.ValidationException;
import org.junit.jupiter.api.Test;

import java.time.Instant;

import static org.junit.jupiter.api.Assertions.*;

class ValidacaoSimulatorTest {

    private final AssinaturaSimulator criarSim = new AssinaturaSimulator();
    private final ValidacaoSimulator validarSim = new ValidacaoSimulator();

    private String generateSimulatedJws() throws ValidationException {
        var mat = new MaterialCriptografico();
        mat.setTipo("PEM");
        mat.setChavePrivada("-----BEGIN PRIVATE KEY-----\nfake\n-----END PRIVATE KEY-----");
        var req = new CriarAssinaturaRequest();
        req.setBundle("{\"resourceType\":\"Bundle\"}");
        req.setProvenance("{\"resourceType\":\"Provenance\"}");
        req.setMaterialCriptografico(mat);
        return criarSim.criar(req).jws();
    }

    private ValidarAssinaturaRequest validRequest(String jws) {
        var req = new ValidarAssinaturaRequest();
        req.setJws(jws);
        req.setPoliticaRevogacao("warn");
        req.setTimestampReferencia(Instant.now().getEpochSecond());
        req.setPoliticaAssinatura("https://fhir.saude.go.gov.br/r4/seguranca/assinatura/v1");
        return req;
    }

    @Test
    void acceptsSimulatedSignature() throws ValidationException {
        String jws = generateSimulatedJws();
        var outcome = validarSim.validar(validRequest(jws));
        assertFalse(outcome.hasErrors(), "A simulated JWS should be valid");
    }

    @Test
    void rejectsNonSimulatedSignature() throws ValidationException {
        String fakeJws = "eyJhbGciOiJSUzI1NiJ9.eyJpc3MiOiJ1bmtub3duIn0.c2ln";
        var outcome = validarSim.validar(validRequest(fakeJws));
        assertTrue(outcome.hasErrors(), "A non-simulated JWS should be invalid");
    }

    @Test
    void rejectsInvalidParameters() {
        var req = new ValidarAssinaturaRequest();
        assertThrows(ValidationException.class, () -> validarSim.validar(req));
    }
}
