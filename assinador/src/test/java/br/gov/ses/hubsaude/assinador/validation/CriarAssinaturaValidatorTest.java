package br.gov.ses.hubsaude.assinador.validation;

import br.gov.ses.hubsaude.assinador.model.CriarAssinaturaRequest;
import br.gov.ses.hubsaude.assinador.model.MaterialCriptografico;
import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

class CriarAssinaturaValidatorTest {

    private final CriarAssinaturaValidator validator = new CriarAssinaturaValidator();

    // --- Helpers ---

    private static CriarAssinaturaRequest validRequest(String tipoMaterial) {
        var mat = new MaterialCriptografico();
        mat.setTipo(tipoMaterial);
        switch (tipoMaterial.toUpperCase()) {
            case "PEM" -> mat.setChavePrivada("-----BEGIN PRIVATE KEY-----\nfake\n-----END PRIVATE KEY-----");
            case "PKCS12" -> { mat.setAlias("a"); mat.setSenha("s"); mat.setConteudo("Y29udGV1ZG8="); }
            case "SMARTCARD", "TOKEN" -> { mat.setTokenLabel("t"); mat.setSlotId(0); mat.setIdentificador("id"); mat.setPin("1234"); }
            case "REMOTE" -> { mat.setCredenciais("cred"); mat.setEnderecoServico("https://remote.example"); }
        }
        var req = new CriarAssinaturaRequest();
        req.setBundle("{\"resourceType\":\"Bundle\"}");
        req.setProvenance("{\"resourceType\":\"Provenance\"}");
        req.setMaterialCriptografico(mat);
        return req;
    }

    // --- Bundle ---

    @Test
    void rejectsMissingBundle() {
        var req = validRequest("PEM");
        req.setBundle(null);
        var ex = assertThrows(ValidationException.class, () -> validator.validate(req));
        assertTrue(ex.getOutcome().getIssue().stream().anyMatch(
                i -> i.getLocation().contains("bundle")));
    }

    @Test
    void rejectsInvalidJsonBundle() {
        var req = validRequest("PEM");
        req.setBundle("not json");
        var ex = assertThrows(ValidationException.class, () -> validator.validate(req));
        assertTrue(ex.getOutcome().getIssue().stream().anyMatch(
                i -> i.getLocation().contains("bundle") && "invalid".equals(i.getCode())));
    }

    // --- Provenance ---

    @Test
    void rejectsMissingProvenance() {
        var req = validRequest("PEM");
        req.setProvenance(null);
        var ex = assertThrows(ValidationException.class, () -> validator.validate(req));
        assertTrue(ex.getOutcome().getIssue().stream().anyMatch(
                i -> i.getLocation().contains("provenance")));
    }

    // --- Material: missing ---

    @Test
    void rejectsMissingMaterial() {
        var req = validRequest("PEM");
        req.setMaterialCriptografico(null);
        var ex = assertThrows(ValidationException.class, () -> validator.validate(req));
        assertTrue(ex.getOutcome().getIssue().stream().anyMatch(
                i -> i.getLocation().contains("materialCriptografico")));
    }

    // --- Material: bad type ---

    @Test
    void rejectsUnknownMaterialType() {
        var req = validRequest("PEM");
        req.getMaterialCriptografico().setTipo("INVALID_TYPE");
        var ex = assertThrows(ValidationException.class, () -> validator.validate(req));
        assertTrue(ex.getMessage().contains("Tipo de material inválido"));
    }

    // --- PEM ---

    @Test
    void acceptsValidPem() {
        assertDoesNotThrow(() -> validator.validate(validRequest("PEM")));
    }

    @Test
    void rejectsPemWithoutPrivateKey() {
        var req = validRequest("PEM");
        req.getMaterialCriptografico().setChavePrivada(null);
        assertThrows(ValidationException.class, () -> validator.validate(req));
    }

    // --- PKCS12 ---

    @Test
    void acceptsValidPkcs12() {
        assertDoesNotThrow(() -> validator.validate(validRequest("PKCS12")));
    }

    @Test
    void rejectsPkcs12WithoutAlias() {
        var req = validRequest("PKCS12");
        req.getMaterialCriptografico().setAlias(null);
        assertThrows(ValidationException.class, () -> validator.validate(req));
    }

    @Test
    void rejectsPkcs12WithoutSenha() {
        var req = validRequest("PKCS12");
        req.getMaterialCriptografico().setSenha(null);
        assertThrows(ValidationException.class, () -> validator.validate(req));
    }

    @Test
    void rejectsPkcs12WithoutConteudo() {
        var req = validRequest("PKCS12");
        req.getMaterialCriptografico().setConteudo(null);
        assertThrows(ValidationException.class, () -> validator.validate(req));
    }

    // --- SMARTCARD (PKCS#11) ---

    @Test
    void acceptsValidSmartcard() {
        assertDoesNotThrow(() -> validator.validate(validRequest("SMARTCARD")));
    }

    @Test
    void rejectsSmartcardWithoutPin() {
        var req = validRequest("SMARTCARD");
        req.getMaterialCriptografico().setPin(null);
        assertThrows(ValidationException.class, () -> validator.validate(req));
    }

    @Test
    void rejectsSmartcardWithoutSlotId() {
        var req = validRequest("SMARTCARD");
        req.getMaterialCriptografico().setSlotId(null);
        assertThrows(ValidationException.class, () -> validator.validate(req));
    }

    // --- TOKEN (PKCS#11) ---

    @Test
    void acceptsValidToken() {
        assertDoesNotThrow(() -> validator.validate(validRequest("TOKEN")));
    }

    // --- REMOTE ---

    @Test
    void acceptsValidRemote() {
        assertDoesNotThrow(() -> validator.validate(validRequest("REMOTE")));
    }

    @Test
    void rejectsRemoteWithoutCredenciais() {
        var req = validRequest("REMOTE");
        req.getMaterialCriptografico().setCredenciais(null);
        assertThrows(ValidationException.class, () -> validator.validate(req));
    }

    @Test
    void rejectsRemoteWithoutEnderecoServico() {
        var req = validRequest("REMOTE");
        req.getMaterialCriptografico().setEnderecoServico(null);
        assertThrows(ValidationException.class, () -> validator.validate(req));
    }

    // --- Multiple errors reported at once ---

    @Test
    void reportsMultipleErrors() {
        var req = new CriarAssinaturaRequest();
        // all fields null
        var ex = assertThrows(ValidationException.class, () -> validator.validate(req));
        assertTrue(ex.getOutcome().getIssue().size() >= 3,
                "Should report at least 3 errors (bundle, provenance, material)");
    }
}
