package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/thiagotelho/runner/simulador/internal/process"
)

func pararCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "parar",
		Short: "Para o Simulador HubSaúde",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := process.Stop(); err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), "Simulador encerrado.")
			return nil
		},
	}
}
