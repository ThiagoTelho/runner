package cmd

import (
	"context"

	"github.com/spf13/cobra"
)

func cmdContext(c context.Context) context.Context {
	if c != nil {
		return c
	}
	return context.Background()
}

const version = "0.1.0"

func Root() *cobra.Command {
	root := &cobra.Command{
		Use:     "assinatura",
		Short:   "CLI do Sistema Runner para assinatura digital simulada (HubSaúde)",
		Version: version,
	}

	root.AddCommand(criarCmd())
	root.AddCommand(validarCmd())
	root.AddCommand(servidorCmd())

	return root
}
