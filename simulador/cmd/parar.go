package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func pararCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "parar",
		Short: "Para o Simulador HubSaúde",
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: stop running simulador.jar process
			fmt.Println("parar: ainda não implementado")
			return nil
		},
	}
}
