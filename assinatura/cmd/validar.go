package cmd

import (
	"github.com/spf13/cobra"
	"github.com/thiagotelho/runner/assinatura/internal/assinador"
)

func validarCmd() *cobra.Command {
	var jwsPath, politicaRevogacao, politicaAssinatura, bundlePath string
	var timestamp int64
	var modoLocal bool
	var porta int

	cmd := &cobra.Command{
		Use:   "validar",
		Short: "Valida uma assinatura digital simulada (JWS)",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmdContext(cmd.Context())
			java, err := assinador.FindJava()
			if err != nil {
				return err
			}
			jar, err := assinador.FindAssinadorJar()
			if err != nil {
				return err
			}
			p := assinador.ValidarParams{
				JwsPath:            jwsPath,
				PoliticaRevogacao:  politicaRevogacao,
				TimestampUnixUTC:   timestamp,
				PoliticaAssinatura: politicaAssinatura,
				BundlePathOpcional: bundlePath,
			}
			return assinador.RunValidar(ctx, java, jar, porta, modoLocal, p,
				cmd.OutOrStdout(), cmd.ErrOrStderr())
		},
	}

	cmd.Flags().StringVarP(&jwsPath, "jws", "j", "", "Caminho do arquivo JWS (obrigatório)")
	cmd.Flags().StringVarP(&politicaRevogacao, "politica-revogacao", "r", "", "warn | soft-fail | strict (obrigatório)")
	cmd.Flags().Int64VarP(&timestamp, "timestamp", "t", 0, "Timestamp de referência Unix UTC (obrigatório)")
	cmd.Flags().StringVarP(&politicaAssinatura, "politica-assinatura", "a", "", "URI da política de assinatura (obrigatório)")
	cmd.Flags().StringVarP(&bundlePath, "bundle", "b", "", "Bundle original (opcional)")
	cmd.Flags().BoolVar(&modoLocal, "local", false, "Forçar invocação local (cold start), ignorando modo servidor")
	cmd.Flags().IntVar(&porta, "porta", assinador.DefaultPort, "Porta do servidor HTTP")
	_ = cmd.MarkFlagRequired("jws")
	_ = cmd.MarkFlagRequired("politica-revogacao")
	_ = cmd.MarkFlagRequired("timestamp")
	_ = cmd.MarkFlagRequired("politica-assinatura")

	return cmd
}
