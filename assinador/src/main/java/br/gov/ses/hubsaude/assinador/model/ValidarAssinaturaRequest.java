package br.gov.ses.hubsaude.assinador.model;

/**
 * Requisição para validação de assinatura digital (simulada).
 * Conforme caso de uso FHIR "validar-assinatura".
 */
public class ValidarAssinaturaRequest {

    private String jws;
    private String politicaRevogacao;
    private Long timestampReferencia;
    private String politicaAssinatura;
    private String bundle;  // optional

    public String getJws()                  { return jws; }
    public String getPoliticaRevogacao()    { return politicaRevogacao; }
    public Long getTimestampReferencia()    { return timestampReferencia; }
    public String getPoliticaAssinatura()   { return politicaAssinatura; }
    public String getBundle()               { return bundle; }

    public void setJws(String jws)                              { this.jws = jws; }
    public void setPoliticaRevogacao(String politicaRevogacao)  { this.politicaRevogacao = politicaRevogacao; }
    public void setTimestampReferencia(Long timestampReferencia){ this.timestampReferencia = timestampReferencia; }
    public void setPoliticaAssinatura(String politicaAssinatura){ this.politicaAssinatura = politicaAssinatura; }
    public void setBundle(String bundle)                        { this.bundle = bundle; }
}
