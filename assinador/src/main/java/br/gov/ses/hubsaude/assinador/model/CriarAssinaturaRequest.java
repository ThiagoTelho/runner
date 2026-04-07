package br.gov.ses.hubsaude.assinador.model;

/**
 * Requisição para criação de assinatura digital (simulada).
 * Conforme caso de uso FHIR "criar-assinatura".
 */
public class CriarAssinaturaRequest {

    private String bundle;
    private String provenance;
    private MaterialCriptografico materialCriptografico;

    public String getBundle()                           { return bundle; }
    public String getProvenance()                       { return provenance; }
    public MaterialCriptografico getMaterialCriptografico() { return materialCriptografico; }

    public void setBundle(String bundle)                { this.bundle = bundle; }
    public void setProvenance(String provenance)        { this.provenance = provenance; }
    public void setMaterialCriptografico(MaterialCriptografico m) { this.materialCriptografico = m; }
}
