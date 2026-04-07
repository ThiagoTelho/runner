package br.gov.ses.hubsaude.assinador.model;

import com.google.gson.annotations.SerializedName;

/**
 * Material criptográfico do signatário conforme a especificação FHIR
 * de criação de assinatura. Os campos obrigatórios variam por {@link TipoMaterial}.
 *
 * <ul>
 *   <li><b>PEM</b>: chavePrivada (obrigatório), senha (opcional)</li>
 *   <li><b>PKCS12</b>: alias, senha, conteudo (todos obrigatórios)</li>
 *   <li><b>SMARTCARD / TOKEN</b>: tokenLabel, slotId, identificador, pin</li>
 *   <li><b>REMOTE</b>: credenciais, enderecoServico</li>
 * </ul>
 */
public class MaterialCriptografico {

    private String tipo;

    // PEM
    private String chavePrivada;
    private String senha;         // also used by PKCS12

    // PKCS12
    private String alias;
    private String conteudo;      // base64

    // SMARTCARD / TOKEN (PKCS#11)
    private String tokenLabel;
    private Integer slotId;
    private String identificador;
    private String pin;

    // REMOTE
    private String credenciais;
    @SerializedName("enderecoServico")
    private String enderecoServico;

    public String getTipo()            { return tipo; }
    public String getChavePrivada()    { return chavePrivada; }
    public String getSenha()           { return senha; }
    public String getAlias()           { return alias; }
    public String getConteudo()        { return conteudo; }
    public String getTokenLabel()      { return tokenLabel; }
    public Integer getSlotId()         { return slotId; }
    public String getIdentificador()   { return identificador; }
    public String getPin()             { return pin; }
    public String getCredenciais()     { return credenciais; }
    public String getEnderecoServico() { return enderecoServico; }

    public void setTipo(String tipo)                     { this.tipo = tipo; }
    public void setChavePrivada(String chavePrivada)     { this.chavePrivada = chavePrivada; }
    public void setSenha(String senha)                   { this.senha = senha; }
    public void setAlias(String alias)                   { this.alias = alias; }
    public void setConteudo(String conteudo)             { this.conteudo = conteudo; }
    public void setTokenLabel(String tokenLabel)         { this.tokenLabel = tokenLabel; }
    public void setSlotId(Integer slotId)                { this.slotId = slotId; }
    public void setIdentificador(String identificador)   { this.identificador = identificador; }
    public void setPin(String pin)                       { this.pin = pin; }
    public void setCredenciais(String credenciais)       { this.credenciais = credenciais; }
    public void setEnderecoServico(String enderecoServico) { this.enderecoServico = enderecoServico; }
}
