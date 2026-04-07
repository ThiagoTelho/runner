package br.gov.ses.hubsaude.assinador.validation;

import br.gov.ses.hubsaude.assinador.model.OperationOutcome;

/**
 * Exceção lançada quando a validação de parâmetros falha.
 * Carrega um {@link OperationOutcome} com todos os erros encontrados.
 */
public class ValidationException extends Exception {

    private final OperationOutcome outcome;

    public ValidationException(OperationOutcome outcome) {
        super(buildMessage(outcome));
        this.outcome = outcome;
    }

    public OperationOutcome getOutcome() {
        return outcome;
    }

    private static String buildMessage(OperationOutcome outcome) {
        var sb = new StringBuilder("Erro(s) de validação:");
        for (var issue : outcome.getIssue()) {
            sb.append("\n  - ").append(issue.getDiagnostics());
            if (issue.getLocation() != null && !issue.getLocation().isEmpty()) {
                sb.append(" [").append(String.join(", ", issue.getLocation())).append("]");
            }
        }
        return sb.toString();
    }
}
