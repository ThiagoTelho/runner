package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func statusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Exibe o status atual do Simulador HubSaúde",
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: check if simulador.jar is running, show port/PID
			fmt.Println("status: ainda não implementado")
			return nil
		},
	}
}
