package br.gov.ses.hubsaude.assinador.simulation;

import br.gov.ses.hubsaude.assinador.model.CriarAssinaturaRequest;
import br.gov.ses.hubsaude.assinador.model.MaterialCriptografico;
import br.gov.ses.hubsaude.assinador.validation.ValidationException;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

class AssinaturaSimulatorTest {

    private final AssinaturaSimulator simulator = new AssinaturaSimulator();

    private CriarAssinaturaRequest validRequest() {
        var mat = new MaterialCriptografico();
        mat.setTipo("PEM");
        mat.setChavePrivada("-----BEGIN PRIVATE KEY-----\nfake\n-----END PRIVATE KEY-----");
        var req = new CriarAssinaturaRequest();
        req.setBundle("{\"resourceType\":\"Bundle\"}");
        req.setProvenance("{\"resourceType\":\"Provenance\"}");
        req.setMaterialCriptografico(mat);
        return req;
    }

    @Test
    void criarReturnsJwsWithThreeParts() throws ValidationException {
        var result = simulator.criar(validRequest());
        assertNotNull(result.jws());
        assertEquals(3, result.jws().split("\\.").length, "JWS should have header.payload.signature");
    }

    @Test
    void criarReturnsSuccessOutcome() throws ValidationException {
        var result = simulator.criar(validRequest());
        assertFalse(result.outcome().hasErrors());
    }

    @Test
    void criarRejectsInvalidRequest() {
        var req = new CriarAssinaturaRequest();
        assertThrows(ValidationException.class, () -> simulator.criar(req));
    }

    @Test
    void criarGeneratesUniqueJwsPerCall() throws ValidationException {
        var r1 = simulator.criar(validRequest());
        var r2 = simulator.criar(validRequest());
        assertNotEquals(r1.jws(), r2.jws());
    }
}
