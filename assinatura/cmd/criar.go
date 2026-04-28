package cmd

import (
	"github.com/spf13/cobra"
	"github.com/thiagotelho/runner/assinatura/internal/assinador"
)

func criarCmd() *cobra.Command {
	var bundlePath, provenancePath, materialPath string
	var modoLocal bool
	var porta int

	cmd := &cobra.Command{
		Use:   "criar",
		Short: "Cria uma assinatura digital simulada",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmdContext(cmd.Context())
			java, err := assinador.FindOrProvisionJava(ctx)
			if err != nil {
				return err
			}
			jar, err := assinador.FindAssinadorJar()
			if err != nil {
				return err
			}
			return assinador.RunCriar(ctx, java, jar, porta, modoLocal,
				bundlePath, provenancePath, materialPath,
				cmd.OutOrStdout(), cmd.ErrOrStderr())
		},
	}

	cmd.Flags().StringVarP(&bundlePath, "bundle", "b", "", "Caminho do Bundle FHIR R4 (obrigatório)")
	cmd.Flags().StringVarP(&provenancePath, "provenance", "p", "", "Caminho do Provenance FHIR R4 (obrigatório)")
	cmd.Flags().StringVarP(&materialPath, "material", "m", "", "Caminho do material criptográfico (obrigatório)")
	cmd.Flags().BoolVar(&modoLocal, "local", false, "Forçar invocação local (cold start), ignorando modo servidor")
	cmd.Flags().IntVar(&porta, "porta", assinador.DefaultPort, "Porta do servidor HTTP")
	_ = cmd.MarkFlagRequired("bundle")
	_ = cmd.MarkFlagRequired("provenance")
	_ = cmd.MarkFlagRequired("material")

	return cmd
}
