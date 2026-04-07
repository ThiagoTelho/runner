package br.gov.ses.hubsaude.assinador.cli;

import br.gov.ses.hubsaude.assinador.model.CriarAssinaturaRequest;
import br.gov.ses.hubsaude.assinador.model.MaterialCriptografico;
import br.gov.ses.hubsaude.assinador.simulation.AssinaturaSimulator;
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
    name = "criar",
    description = "Cria uma assinatura digital simulada a partir de um Bundle e Provenance FHIR.",
    mixinStandardHelpOptions = true
)
public class CriarCommand implements Callable<Integer> {

    private static final Gson GSON = new GsonBuilder().setPrettyPrinting().create();

    @Option(names = {"-b", "--bundle"}, required = true,
            description = "Caminho para o arquivo JSON do Bundle FHIR R4.")
    private Path bundlePath;

    @Option(names = {"-p", "--provenance"}, required = true,
            description = "Caminho para o arquivo JSON do Provenance FHIR R4.")
    private Path provenancePath;

    @Option(names = {"-m", "--material"}, required = true,
            description = "Caminho para o arquivo JSON com o material criptográfico do signatário.")
    private Path materialPath;

    @Override
    public Integer call() {
        try {
            String bundle     = readFile(bundlePath);
            String provenance = readFile(provenancePath);
            String materialJson = readFile(materialPath);

            var material = GSON.fromJson(materialJson, MaterialCriptografico.class);

            var request = new CriarAssinaturaRequest();
            request.setBundle(bundle);
            request.setProvenance(provenance);
            request.setMaterialCriptografico(material);

            var simulator = new AssinaturaSimulator();
            var result = simulator.criar(request);

            System.out.println("JWS gerado:");
            System.out.println(result.jws());
            System.out.println();
            System.out.println(GSON.toJson(result.outcome()));
            return 0;

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
