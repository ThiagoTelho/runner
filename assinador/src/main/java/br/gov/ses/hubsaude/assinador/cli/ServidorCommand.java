package br.gov.ses.hubsaude.assinador.cli;

import br.gov.ses.hubsaude.assinador.server.AssinadorHttpServer;
import picocli.CommandLine.Command;
import picocli.CommandLine.Option;

import java.util.concurrent.Callable;

@Command(
    name = "servidor",
    description = "Inicia o assinador em modo servidor HTTP (warm start).",
    mixinStandardHelpOptions = true
)
public class ServidorCommand implements Callable<Integer> {

    @Option(names = {"-p", "--porta"}, defaultValue = "8190",
            description = "Porta HTTP (padrão: ${DEFAULT-VALUE}).")
    private int porta;

    @Override
    public Integer call() {
        var server = new AssinadorHttpServer();
        try {
            server.start(porta);
            Runtime.getRuntime().addShutdownHook(new Thread(server::stop));
            System.out.println("Pressione Ctrl+C para encerrar.");
            Thread.currentThread().join();
            return 0;
        } catch (java.net.BindException e) {
            System.err.printf("Erro: porta %d já está em uso. Verifique se há outra instância em execução.%n", porta);
            return 1;
        } catch (Exception e) {
            System.err.println("Erro ao iniciar servidor: " + e.getMessage());
            return 2;
        }
    }
}
