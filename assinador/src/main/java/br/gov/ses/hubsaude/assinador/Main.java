package br.gov.ses.hubsaude.assinador;

import br.gov.ses.hubsaude.assinador.cli.AssinadorCli;
import picocli.CommandLine;

public class Main {
    public static void main(String[] args) {
        int exitCode = new CommandLine(new AssinadorCli()).execute(args);
        System.exit(exitCode);
    }
}
