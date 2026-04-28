package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/thiagotelho/runner/assinatura/internal/assinador"
)

func servidorCmd() *cobra.Command {
	parent := &cobra.Command{
		Use:   "servidor",
		Short: "Gerencia o assinador.jar em modo servidor HTTP",
	}

	var porta int

	// --- iniciar ---
	var inatividade int
	iniciar := &cobra.Command{
		Use:   "iniciar",
		Short: "Inicia o assinador.jar em modo servidor (background)",
		RunE: func(cmd *cobra.Command, args []string) error {
			java, err := assinador.FindJava()
			if err != nil {
				return err
			}
			jar, err := assinador.FindAssinadorJar()
			if err != nil {
				return err
			}
			client := assinador.NewHTTPClient(porta)
			if client.CheckHealth(cmd.Context()) {
				fmt.Fprintf(cmd.OutOrStdout(), "Assinador já está em execução na porta %d.\n", porta)
				return nil
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Iniciando assinador na porta %d...\n", porta)
			if err := assinador.StartServer(java, jar, porta, inatividade); err != nil {
				return err
			}
			if err := assinador.WaitForServer(cmd.Context(), client, 15*time.Second); err != nil {
				return err
			}
			if inatividade > 0 {
				fmt.Fprintf(cmd.OutOrStdout(), "Assinador iniciado na porta %d (timeout de inatividade: %d min).\n", porta, inatividade)
			} else {
				fmt.Fprintf(cmd.OutOrStdout(), "Assinador iniciado na porta %d.\n", porta)
			}
			return nil
		},
	}
	iniciar.Flags().IntVar(&inatividade, "inatividade-minutos", 0, "Encerrar após N minutos sem interação (0 = nunca)")

	// --- parar ---
	parar := &cobra.Command{
		Use:   "parar",
		Short: "Para o assinador.jar em modo servidor",
		RunE: func(cmd *cobra.Command, args []string) error {
			client := assinador.NewHTTPClient(porta)
			if err := client.Shutdown(cmd.Context()); err != nil {
				return err
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Assinador na porta %d encerrado.\n", porta)
			return nil
		},
	}

	// --- status ---
	status := &cobra.Command{
		Use:   "status",
		Short: "Verifica se o assinador.jar está em execução",
		RunE: func(cmd *cobra.Command, args []string) error {
			client := assinador.NewHTTPClient(porta)
			if client.CheckHealth(cmd.Context()) {
				fmt.Fprintf(cmd.OutOrStdout(), "Assinador em execução na porta %d.\n", porta)
			} else {
				fmt.Fprintf(cmd.OutOrStdout(), "Assinador não está em execução na porta %d.\n", porta)
			}
			return nil
		},
	}

	parent.PersistentFlags().IntVarP(&porta, "porta", "p", assinador.DefaultPort, "Porta do servidor HTTP")
	parent.AddCommand(iniciar, parar, status)

	return parent
}
