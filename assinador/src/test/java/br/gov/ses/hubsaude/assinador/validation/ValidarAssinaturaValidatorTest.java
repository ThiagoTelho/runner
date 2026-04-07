package br.gov.ses.hubsaude.assinador.validation;

import br.gov.ses.hubsaude.assinador.model.ValidarAssinaturaRequest;
import org.junit.jupiter.api.Test;

import java.time.Instant;

import static org.junit.jupiter.api.Assertions.*;

class ValidarAssinaturaValidatorTest {

    private final ValidarAssinaturaValidator validator = new ValidarAssinaturaValidator();

    private ValidarAssinaturaRequest validRequest() {
        var req = new ValidarAssinaturaRequest();
        req.setJws("eyJhbGciOiJSUzI1NiJ9.eyJpc3MiOiJ0ZXN0In0.c2ln");
        req.setPoliticaRevogacao("warn");
        req.setTimestampReferencia(Instant.now().getEpochSecond());
        req.setPoliticaAssinatura("https://fhir.saude.go.gov.br/r4/seguranca/assinatura/v1");
        return req;
    }

    @Test
    void acceptsValidRequest() {
        assertDoesNotThrow(() -> validator.validate(validRequest()));
    }

    @Test
    void rejectsMissingJws() {
        var req = validRequest();
        req.setJws(null);
        var ex = assertThrows(ValidationException.class, () -> validator.validate(req));
        assertTrue(ex.getOutcome().getIssue().stream().anyMatch(i -> i.getLocation().contains("jws")));
    }

    @Test
    void rejectsMissingPoliticaRevogacao() {
        var req = validRequest();
        req.setPoliticaRevogacao(null);
        assertThrows(ValidationException.class, () -> validator.validate(req));
    }

    @Test
    void rejectsInvalidPoliticaRevogacao() {
        var req = validRequest();
        req.setPoliticaRevogacao("invalid-value");
        var ex = assertThrows(ValidationException.class, () -> validator.validate(req));
        assertTrue(ex.getMessage().contains("Política de revogação inválida"));
    }

    @Test
    void acceptsAllValidPoliticaRevogacao() {
        for (String value : new String[]{"warn", "soft-fail", "strict"}) {
            var req = validRequest();
            req.setPoliticaRevogacao(value);
            assertDoesNotThrow(() -> validator.validate(req));
        }
    }

    @Test
    void rejectsMissingTimestamp() {
        var req = validRequest();
        req.setTimestampReferencia(null);
        assertThrows(ValidationException.class, () -> validator.validate(req));
    }

    @Test
    void rejectsTimestampBelowMinimum() {
        var req = validRequest();
        req.setTimestampReferencia(1000L);
        var ex = assertThrows(ValidationException.class, () -> validator.validate(req));
        assertTrue(ex.getMessage().contains("fora do intervalo"));
    }

    @Test
    void rejectsTimestampAboveMaximum() {
        var req = validRequest();
        req.setTimestampReferencia(5000000000L);
        var ex = assertThrows(ValidationException.class, () -> validator.validate(req));
        assertTrue(ex.getMessage().contains("fora do intervalo"));
    }

    @Test
    void rejectsMissingPoliticaAssinatura() {
        var req = validRequest();
        req.setPoliticaAssinatura(null);
        assertThrows(ValidationException.class, () -> validator.validate(req));
    }

    @Test
    void reportsMultipleErrors() {
        var req = new ValidarAssinaturaRequest();
        var ex = assertThrows(ValidationException.class, () -> validator.validate(req));
        assertTrue(ex.getOutcome().getIssue().size() >= 4,
                "Should report at least 4 errors (jws, politicaRevogacao, timestamp, politicaAssinatura)");
    }
}
