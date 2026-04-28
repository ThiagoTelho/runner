package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/thiagotelho/runner/simulador/internal/process"
)

func statusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Exibe o status atual do Simulador HubSaúde",
		RunE: func(cmd *cobra.Command, args []string) error {
			running, pid, porta := process.Status()
			if running {
				fmt.Fprintf(cmd.OutOrStdout(), "Simulador em execução (PID %d, porta %d).\n", pid, porta)
			} else {
				fmt.Fprintln(cmd.OutOrStdout(), "Simulador não está em execução.")
			}
			return nil
		},
	}
}
