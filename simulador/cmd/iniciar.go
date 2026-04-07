package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func iniciarCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "iniciar",
		Short: "Inicia o Simulador HubSaúde (baixa automaticamente se necessário)",
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: download simulador.jar if needed, check ports, start
			fmt.Println("iniciar: ainda não implementado")
			return nil
		},
	}
}
