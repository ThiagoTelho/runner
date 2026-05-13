package br.gov.ses.hubsaude.assinador.model;

import org.junit.jupiter.api.Test;

import static org.junit.jupiter.api.Assertions.*;

class OperationOutcomeTest {

    @Test
    void emptyOutcomeHasNoErrors() {
        var outcome = new OperationOutcome();
        assertFalse(outcome.hasErrors());
        assertTrue(outcome.getIssue().isEmpty());
        assertEquals("OperationOutcome", outcome.getResourceType());
    }

    @Test
    void successFactoryProducesInformationIssue() {
        var outcome = OperationOutcome.success("Tudo certo");
        assertFalse(outcome.hasErrors());
        assertEquals(1, outcome.getIssue().size());
        var issue = outcome.getIssue().get(0);
        assertEquals("information", issue.getSeverity());
        assertEquals("informational", issue.getCode());
        assertEquals("Tudo certo", issue.getDiagnostics());
        assertTrue(issue.getLocation().isEmpty());
    }

    @Test
    void errorFactoryProducesErrorIssue() {
        var outcome = OperationOutcome.error("falhou", "campo.x");
        assertTrue(outcome.hasErrors());
        var issue = outcome.getIssue().get(0);
        assertEquals("error", issue.getSeverity());
        assertEquals("invalid", issue.getCode());
        assertEquals("falhou", issue.getDiagnostics());
        assertEquals(1, issue.getLocation().size());
        assertEquals("campo.x", issue.getLocation().get(0));
    }

    @Test
    void errorFactoryAcceptsNullLocation() {
        var outcome = OperationOutcome.error("erro sem campo", null);
        assertTrue(outcome.hasErrors());
        assertTrue(outcome.getIssue().get(0).getLocation().isEmpty());
    }

    @Test
    void requiredFieldFactoryProducesRequiredIssue() {
        var outcome = OperationOutcome.requiredField("bundle");
        assertTrue(outcome.hasErrors());
        var issue = outcome.getIssue().get(0);
        assertEquals("error", issue.getSeverity());
        assertEquals("required", issue.getCode());
        assertTrue(issue.getDiagnostics().contains("bundle"));
        assertEquals("bundle", issue.getLocation().get(0));
    }

    @Test
    void addIssueIsChainable() {
        var outcome = new OperationOutcome()
                .addIssue(new Issue(Issue.Severity.WARNING, Issue.Code.PROCESSING, "aviso", java.util.List.of()))
                .addIssue(new Issue(Issue.Severity.ERROR, Issue.Code.INVALID, "erro", java.util.List.of("x")));
        assertEquals(2, outcome.getIssue().size());
        assertTrue(outcome.hasErrors());
    }

    @Test
    void hasErrorsIgnoresWarningsAndInfo() {
        var outcome = new OperationOutcome()
                .addIssue(new Issue(Issue.Severity.WARNING, Issue.Code.PROCESSING, "aviso", java.util.List.of()))
                .addIssue(new Issue(Issue.Severity.INFORMATION, Issue.Code.INFORMATIONAL, "info", java.util.List.of()));
        assertFalse(outcome.hasErrors());
    }
}
