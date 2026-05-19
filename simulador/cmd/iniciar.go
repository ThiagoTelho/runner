package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/thiagotelho/runner/jdk"
	"github.com/thiagotelho/runner/simulador/internal/process"
)

func iniciarCmd() *cobra.Command {
	var porta int

	cmd := &cobra.Command{
		Use:   "iniciar",
		Short: "Inicia o Simulador HubSaúde (baixa automaticamente se necessário)",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()

			// Verifica se já está rodando.
			if running, pid, pt := process.Status(); running {
				fmt.Fprintf(cmd.OutOrStdout(), "Simulador já está em execução (PID %d, porta %d).\n", pid, pt)
				return nil
			}

			// Garante que o JAR está disponível.
			jarPath, err := process.FindOrDownload(ctx)
			if err != nil {
				return err
			}

			// Garante que o JDK está disponível.
			javaExe, err := jdk.FindOrProvision(ctx)
			if err != nil {
				return err
			}

			// Verifica porta antes de subir.
			if err := process.CheckPort(porta); err != nil {
				return err
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Iniciando Simulador na porta %d...\n", porta)
			if err := process.Start(javaExe, jarPath, porta); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Simulador iniciado (porta %d).\n", porta)
			return nil
		},
	}

	cmd.Flags().IntVarP(&porta, "porta", "p", process.DefaultPort, "Porta HTTP do Simulador")
	return cmd
}
