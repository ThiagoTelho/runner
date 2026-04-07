package br.gov.ses.hubsaude.assinador.simulation;

import br.gov.ses.hubsaude.assinador.model.Issue;
import br.gov.ses.hubsaude.assinador.model.OperationOutcome;
import br.gov.ses.hubsaude.assinador.model.ValidarAssinaturaRequest;
import br.gov.ses.hubsaude.assinador.validation.ValidarAssinaturaValidator;
import br.gov.ses.hubsaude.assinador.validation.ValidationException;

import java.nio.charset.StandardCharsets;
import java.util.Base64;
import java.util.List;

/**
 * Simula a validação de uma assinatura digital.
 * <p>
 * Regra de simulação: o JWS é considerado <b>válido</b> se o payload
 * decodificado contém {@code "runner-assinador-simulado"} (ou seja, foi
 * gerado pelo próprio simulador). Caso contrário, é <b>inválido</b>.
 */
public class ValidacaoSimulator {

    private static final String ISSUER_MARKER = "runner-assinador-simulado";

    private final ValidarAssinaturaValidator validator = new ValidarAssinaturaValidator();

    public OperationOutcome validar(ValidarAssinaturaRequest request) throws ValidationException {
        validator.validate(request);

        boolean valid = isSimulatedSignature(request.getJws());

        if (valid) {
            return OperationOutcome.success("Assinatura válida (simulação).");
        } else {
            return new OperationOutcome().addIssue(new Issue(
                    Issue.Severity.ERROR,
                    Issue.Code.PROCESSING,
                    "Assinatura inválida: o JWS não foi gerado pelo assinador simulado.",
                    List.of("jws")
            ));
        }
    }

    private boolean isSimulatedSignature(String jws) {
        try {
            String[] parts = jws.split("\\.");
            if (parts.length < 2) return false;
            String payload = new String(
                    Base64.getUrlDecoder().decode(parts[1]),
                    StandardCharsets.UTF_8);
            return payload.contains(ISSUER_MARKER);
        } catch (Exception e) {
            return false;
        }
    }
}
