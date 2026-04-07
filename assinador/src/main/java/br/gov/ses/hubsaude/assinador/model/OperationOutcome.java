package br.gov.ses.hubsaude.assinador.model;

import java.util.ArrayList;
import java.util.List;

/**
 * Resultado padronizado no formato FHIR OperationOutcome.
 * Usado tanto para respostas de sucesso quanto de erro.
 */
public class OperationOutcome {

    private final String resourceType = "OperationOutcome";
    private final List<Issue> issue = new ArrayList<>();

    public String getResourceType()   { return resourceType; }
    public List<Issue> getIssue()     { return issue; }

    public OperationOutcome addIssue(Issue issue) {
        this.issue.add(issue);
        return this;
    }

    public boolean hasErrors() {
        return issue.stream().anyMatch(i -> "error".equals(i.getSeverity()));
    }

    // --- Factory helpers ---

    public static OperationOutcome success(String message) {
        return new OperationOutcome().addIssue(new Issue(
                Issue.Severity.INFORMATION,
                Issue.Code.INFORMATIONAL,
                message,
                List.of()
        ));
    }

    public static OperationOutcome error(String diagnostics, String location) {
        return new OperationOutcome().addIssue(new Issue(
                Issue.Severity.ERROR,
                Issue.Code.INVALID,
                diagnostics,
                location == null ? List.of() : List.of(location)
        ));
    }

    public static OperationOutcome requiredField(String field) {
        return new OperationOutcome().addIssue(new Issue(
                Issue.Severity.ERROR,
                Issue.Code.REQUIRED,
                "Campo obrigatório ausente: " + field,
                List.of(field)
        ));
    }
}
