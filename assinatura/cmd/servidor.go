package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func servidorCmd() *cobra.Command {
	parent := &cobra.Command{
		Use:   "servidor",
		Short: "Gerencia o assinador.jar em modo servidor HTTP",
	}

	var porta int

	parent.AddCommand(&cobra.Command{
		Use:   "iniciar",
		Short: "Inicia o assinador.jar em modo servidor",
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: start assinador.jar in server mode
			fmt.Printf("servidor iniciar: porta=%d (ainda não implementado)\n", porta)
			return nil
		},
	})

	parent.AddCommand(&cobra.Command{
		Use:   "parar",
		Short: "Para o assinador.jar em modo servidor",
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: POST /shutdown to running instance
			fmt.Printf("servidor parar: porta=%d (ainda não implementado)\n", porta)
			return nil
		},
	})

	parent.AddCommand(&cobra.Command{
		Use:   "status",
		Short: "Verifica se o assinador.jar está em execução",
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: GET /health on running instance
			fmt.Printf("servidor status: porta=%d (ainda não implementado)\n", porta)
			return nil
		},
	})

	parent.PersistentFlags().IntVarP(&porta, "porta", "p", 8190, "Porta do servidor HTTP")

	return parent
}
