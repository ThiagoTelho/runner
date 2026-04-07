package br.gov.ses.hubsaude.assinador.model;

import java.util.List;

/**
 * Uma issue dentro de um OperationOutcome, seguindo o modelo FHIR.
 */
public class Issue {

    public enum Severity { ERROR, WARNING, INFORMATION }
    public enum Code     { INVALID, REQUIRED, VALUE, PROCESSING, INFORMATIONAL }

    private final String severity;
    private final String code;
    private final String diagnostics;
    private final List<String> location;

    public Issue(Severity severity, Code code, String diagnostics, List<String> location) {
        this.severity    = severity.name().toLowerCase();
        this.code        = code.name().toLowerCase();
        this.diagnostics = diagnostics;
        this.location    = location;
    }

    public String getSeverity()       { return severity; }
    public String getCode()           { return code; }
    public String getDiagnostics()    { return diagnostics; }
    public List<String> getLocation() { return location; }
}
