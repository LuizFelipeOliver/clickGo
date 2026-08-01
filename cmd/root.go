/*
clickGo é um CLI de integração com a API do ClickUp.
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "clickGo",
	Short: "CLI de integração com o ClickUp",
	Long: `clickGo é um CLI para integrar com a API do ClickUp
(tarefas, listas, espaços e time tracking) usando OAuth 2.0.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
