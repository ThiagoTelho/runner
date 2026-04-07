package br.gov.ses.hubsaude.assinador.validation;

import br.gov.ses.hubsaude.assinador.model.*;
import com.google.gson.Gson;
import com.google.gson.JsonSyntaxException;

import java.util.ArrayList;
import java.util.List;

/**
 * Valida todos os parâmetros de uma requisição de criação de assinatura
 * conforme a especificação FHIR "caso-de-uso-criar-assinatura".
 */
public class CriarAssinaturaValidator {

    private static final Gson GSON = new Gson();

    public void validate(CriarAssinaturaRequest request) throws ValidationException {
        var issues = new ArrayList<Issue>();

        validateBundle(request.getBundle(), issues);
        validateProvenance(request.getProvenance(), issues);
        validateMaterial(request.getMaterialCriptografico(), issues);

        if (!issues.isEmpty()) {
            var outcome = new OperationOutcome();
            issues.forEach(outcome::addIssue);
            throw new ValidationException(outcome);
        }
    }

    private void validateBundle(String bundle, List<Issue> issues) {
        if (isBlank(bundle)) {
            issues.add(required("bundle"));
            return;
        }
        if (!isValidJson(bundle)) {
            issues.add(invalid("bundle", "O bundle deve ser um JSON válido (FHIR Bundle R4)."));
        }
    }

    private void validateProvenance(String provenance, List<Issue> issues) {
        if (isBlank(provenance)) {
            issues.add(required("provenance"));
            return;
        }
        if (!isValidJson(provenance)) {
            issues.add(invalid("provenance", "O provenance deve ser um JSON válido (FHIR Provenance R4)."));
        }
    }

    private void validateMaterial(MaterialCriptografico material, List<Issue> issues) {
        if (material == null) {
            issues.add(required("materialCriptografico"));
            return;
        }

        var tipo = TipoMaterial.fromString(material.getTipo());
        if (tipo == null) {
            issues.add(invalid("materialCriptografico.tipo",
                    "Tipo de material inválido ou ausente. Valores aceitos: " + TipoMaterial.valoresAceitos()));
            return;
        }

        switch (tipo) {
            case PEM      -> validatePem(material, issues);
            case PKCS12   -> validatePkcs12(material, issues);
            case SMARTCARD, TOKEN -> validatePkcs11(material, issues);
            case REMOTE   -> validateRemote(material, issues);
        }
    }

    private void validatePem(MaterialCriptografico m, List<Issue> issues) {
        if (isBlank(m.getChavePrivada())) {
            issues.add(required("materialCriptografico.chavePrivada",
                    "A chave privada PKCS#8 (PEM) é obrigatória para o tipo PEM."));
        }
    }

    private void validatePkcs12(MaterialCriptografico m, List<Issue> issues) {
        if (isBlank(m.getAlias())) {
            issues.add(required("materialCriptografico.alias",
                    "O alias é obrigatório para o tipo PKCS12."));
        }
        if (isBlank(m.getSenha())) {
            issues.add(required("materialCriptografico.senha",
                    "A senha é obrigatória para o tipo PKCS12."));
        }
        if (isBlank(m.getConteudo())) {
            issues.add(required("materialCriptografico.conteudo",
                    "O conteúdo base64 do PKCS#12 é obrigatório."));
        }
    }

    private void validatePkcs11(MaterialCriptografico m, List<Issue> issues) {
        String tipoLabel = m.getTipo().toUpperCase();
        if (isBlank(m.getTokenLabel())) {
            issues.add(required("materialCriptografico.tokenLabel",
                    "O token label é obrigatório para o tipo " + tipoLabel + " (PKCS#11)."));
        }
        if (m.getSlotId() == null) {
            issues.add(required("materialCriptografico.slotId",
                    "O slot ID é obrigatório para o tipo " + tipoLabel + " (PKCS#11)."));
        }
        if (isBlank(m.getIdentificador())) {
            issues.add(required("materialCriptografico.identificador",
                    "O identificador da chave é obrigatório para o tipo " + tipoLabel + " (PKCS#11)."));
        }
        if (isBlank(m.getPin())) {
            issues.add(required("materialCriptografico.pin",
                    "O PIN é obrigatório para o tipo " + tipoLabel + " (PKCS#11)."));
        }
    }

    private void validateRemote(MaterialCriptografico m, List<Issue> issues) {
        if (isBlank(m.getCredenciais())) {
            issues.add(required("materialCriptografico.credenciais",
                    "As credenciais são obrigatórias para o tipo REMOTE."));
        }
        if (isBlank(m.getEnderecoServico())) {
            issues.add(required("materialCriptografico.enderecoServico",
                    "O endereço do serviço remoto é obrigatório para o tipo REMOTE."));
        }
    }

    // --- Helpers ---

    private static boolean isBlank(String s) {
        return s == null || s.isBlank();
    }

    private static boolean isValidJson(String s) {
        try {
            GSON.fromJson(s, Object.class);
            return true;
        } catch (JsonSyntaxException e) {
            return false;
        }
    }

    private static Issue required(String field) {
        return new Issue(Issue.Severity.ERROR, Issue.Code.REQUIRED,
                "Campo obrigatório ausente: " + field, List.of(field));
    }

    private static Issue required(String field, String detail) {
        return new Issue(Issue.Severity.ERROR, Issue.Code.REQUIRED, detail, List.of(field));
    }

    private static Issue invalid(String field, String detail) {
        return new Issue(Issue.Severity.ERROR, Issue.Code.INVALID, detail, List.of(field));
    }
}
