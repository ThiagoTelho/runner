package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func criarCmd() *cobra.Command {
	var bundlePath, provenancePath, materialPath string
	var modoLocal bool

	cmd := &cobra.Command{
		Use:   "criar",
		Short: "Cria uma assinatura digital simulada",
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: implement local and HTTP invocation of assinador.jar
			fmt.Println("criar: ainda não implementado")
			fmt.Printf("  bundle=%s provenance=%s material=%s local=%v\n",
				bundlePath, provenancePath, materialPath, modoLocal)
			return nil
		},
	}

	cmd.Flags().StringVarP(&bundlePath, "bundle", "b", "", "Caminho do Bundle FHIR R4 (obrigatório)")
	cmd.Flags().StringVarP(&provenancePath, "provenance", "p", "", "Caminho do Provenance FHIR R4 (obrigatório)")
	cmd.Flags().StringVarP(&materialPath, "material", "m", "", "Caminho do material criptográfico (obrigatório)")
	cmd.Flags().BoolVar(&modoLocal, "local", false, "Forçar invocação local (cold start)")
	_ = cmd.MarkFlagRequired("bundle")
	_ = cmd.MarkFlagRequired("provenance")
	_ = cmd.MarkFlagRequired("material")

	return cmd
}
