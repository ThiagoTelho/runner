package br.gov.ses.hubsaude.assinador.cli;

import picocli.CommandLine.Command;

@Command(
    name = "assinador",
    description = "Assinador digital simulado — Sistema Runner / HubSaúde.",
    version = "assinador 0.1.0",
    mixinStandardHelpOptions = true,
    subcommands = {
        CriarCommand.class,
        ValidarCommand.class,
        ServidorCommand.class
    }
)
public class AssinadorCli {
}
