package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func validarCmd() *cobra.Command {
	var jwsPath, politicaRevogacao, politicaAssinatura, bundlePath string
	var timestamp int64
	var modoLocal bool

	cmd := &cobra.Command{
		Use:   "validar",
		Short: "Valida uma assinatura digital simulada (JWS)",
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: implement local and HTTP invocation of assinador.jar
			fmt.Println("validar: ainda não implementado")
			fmt.Printf("  jws=%s politicaRevogacao=%s timestamp=%d politicaAssinatura=%s local=%v\n",
				jwsPath, politicaRevogacao, timestamp, politicaAssinatura, modoLocal)
			return nil
		},
	}

	cmd.Flags().StringVarP(&jwsPath, "jws", "j", "", "Caminho do arquivo JWS (obrigatório)")
	cmd.Flags().StringVarP(&politicaRevogacao, "politica-revogacao", "r", "", "warn | soft-fail | strict (obrigatório)")
	cmd.Flags().Int64VarP(&timestamp, "timestamp", "t", 0, "Timestamp de referência Unix UTC (obrigatório)")
	cmd.Flags().StringVarP(&politicaAssinatura, "politica-assinatura", "a", "", "URI da política de assinatura (obrigatório)")
	cmd.Flags().StringVarP(&bundlePath, "bundle", "b", "", "Bundle original (opcional)")
	cmd.Flags().BoolVar(&modoLocal, "local", false, "Forçar invocação local (cold start)")
	_ = cmd.MarkFlagRequired("jws")
	_ = cmd.MarkFlagRequired("politica-revogacao")
	_ = cmd.MarkFlagRequired("timestamp")
	_ = cmd.MarkFlagRequired("politica-assinatura")

	return cmd
}
