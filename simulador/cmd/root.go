package cmd

import "github.com/spf13/cobra"

const version = "0.1.0"

func Root() *cobra.Command {
	root := &cobra.Command{
		Use:     "simulador",
		Short:   "CLI do Sistema Runner para gerenciar o Simulador HubSaúde",
		Version: version,
	}

	root.AddCommand(iniciarCmd())
	root.AddCommand(pararCmd())
	root.AddCommand(statusCmd())

	return root
}
