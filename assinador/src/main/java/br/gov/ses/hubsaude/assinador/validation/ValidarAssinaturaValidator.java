package br.gov.ses.hubsaude.assinador.validation;

import br.gov.ses.hubsaude.assinador.model.*;

import java.time.Instant;
import java.util.ArrayList;
import java.util.List;

/**
 * Valida todos os parâmetros de uma requisição de validação de assinatura
 * conforme a especificação FHIR "caso-de-uso-validar-assinatura".
 */
public class ValidarAssinaturaValidator {

    static final long TS_MIN = 1751328000L;   // 2025-07-01
    static final long TS_MAX = 4102444800L;   // 2100-01-01
    static final long CLOCK_SKEW_MAX = 300L;  // ±5 min

    public void validate(ValidarAssinaturaRequest request) throws ValidationException {
        var issues = new ArrayList<Issue>();

        validateJws(request.getJws(), issues);
        validatePoliticaRevogacao(request.getPoliticaRevogacao(), issues);
        validateTimestamp(request.getTimestampReferencia(), issues);
        validatePoliticaAssinatura(request.getPoliticaAssinatura(), issues);

        if (!issues.isEmpty()) {
            var outcome = new OperationOutcome();
            issues.forEach(outcome::addIssue);
            throw new ValidationException(outcome);
        }
    }

    private void validateJws(String jws, List<Issue> issues) {
        if (isBlank(jws)) {
            issues.add(required("jws"));
        }
    }

    private void validatePoliticaRevogacao(String value, List<Issue> issues) {
        if (isBlank(value)) {
            issues.add(required("politicaRevogacao"));
            return;
        }
        if (PoliticaRevogacao.fromString(value) == null) {
            issues.add(invalid("politicaRevogacao",
                    "Política de revogação inválida. Valores aceitos: " + PoliticaRevogacao.valoresAceitos()));
        }
    }

    private void validateTimestamp(Long ts, List<Issue> issues) {
        if (ts == null) {
            issues.add(required("timestampReferencia"));
            return;
        }
        if (ts < TS_MIN || ts > TS_MAX) {
            issues.add(invalid("timestampReferencia",
                    String.format("Timestamp fora do intervalo aceito [%d, %d].", TS_MIN, TS_MAX)));
            return;
        }
        long now = Instant.now().getEpochSecond();
        if (Math.abs(ts - now) > CLOCK_SKEW_MAX) {
            issues.add(new Issue(Issue.Severity.WARNING, Issue.Code.VALUE,
                    String.format("Timestamp difere mais de %ds do relógio local (diferença: %ds).",
                            CLOCK_SKEW_MAX, Math.abs(ts - now)),
                    List.of("timestampReferencia")));
        }
    }

    private void validatePoliticaAssinatura(String uri, List<Issue> issues) {
        if (isBlank(uri)) {
            issues.add(required("politicaAssinatura"));
        }
    }

    // --- Helpers ---

    private static boolean isBlank(String s) {
        return s == null || s.isBlank();
    }

    private static Issue required(String field) {
        return new Issue(Issue.Severity.ERROR, Issue.Code.REQUIRED,
                "Campo obrigatório ausente: " + field, List.of(field));
    }

    private static Issue invalid(String field, String detail) {
        return new Issue(Issue.Severity.ERROR, Issue.Code.INVALID, detail, List.of(field));
    }
}
