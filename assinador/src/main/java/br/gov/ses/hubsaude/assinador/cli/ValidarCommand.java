package br.gov.ses.hubsaude.assinador.cli;

import br.gov.ses.hubsaude.assinador.model.ValidarAssinaturaRequest;
import br.gov.ses.hubsaude.assinador.simulation.ValidacaoSimulator;
import br.gov.ses.hubsaude.assinador.validation.ValidationException;
import com.google.gson.Gson;
import com.google.gson.GsonBuilder;
import picocli.CommandLine.Command;
import picocli.CommandLine.Option;

import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.concurrent.Callable;

@Command(
    name = "validar",
    description = "Valida uma assinatura digital simulada (JWS).",
    mixinStandardHelpOptions = true
)
public class ValidarCommand implements Callable<Integer> {

    private static final Gson GSON = new GsonBuilder().setPrettyPrinting().create();

    @Option(names = {"-j", "--jws"}, required = true,
            description = "Caminho para o arquivo contendo o JWS (assinatura).")
    private Path jwsPath;

    @Option(names = {"-r", "--politica-revogacao"}, required = true,
            description = "Política de revogação: warn, soft-fail ou strict.")
    private String politicaRevogacao;

    @Option(names = {"-t", "--timestamp"}, required = true,
            description = "Timestamp de referência (Unix UTC, em segundos).")
    private long timestamp;

    @Option(names = {"-a", "--politica-assinatura"}, required = true,
            description = "URI versionada da política de assinatura.")
    private String politicaAssinatura;

    @Option(names = {"-b", "--bundle"},
            description = "Caminho para o Bundle original (opcional, para verificação de integridade).")
    private Path bundlePath;

    @Override
    public Integer call() {
        try {
            String jws = readFile(jwsPath);

            var request = new ValidarAssinaturaRequest();
            request.setJws(jws);
            request.setPoliticaRevogacao(politicaRevogacao);
            request.setTimestampReferencia(timestamp);
            request.setPoliticaAssinatura(politicaAssinatura);

            if (bundlePath != null) {
                request.setBundle(readFile(bundlePath));
            }

            var simulator = new ValidacaoSimulator();
            var outcome = simulator.validar(request);

            System.out.println(GSON.toJson(outcome));
            return outcome.hasErrors() ? 1 : 0;

        } catch (ValidationException e) {
            System.err.println(e.getMessage());
            System.err.println(GSON.toJson(e.getOutcome()));
            return 1;
        } catch (IOException e) {
            System.err.println("Erro ao ler arquivo: " + e.getMessage());
            return 2;
        } catch (Exception e) {
            System.err.println("Erro inesperado: " + e.getMessage());
            return 3;
        }
    }

    private String readFile(Path path) throws IOException {
        if (!Files.exists(path)) {
            throw new IOException("Arquivo não encontrado: " + path);
        }
        return Files.readString(path);
    }
}
