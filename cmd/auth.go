package cmd

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"charm.land/lipgloss/v2"
	"github.com/spf13/cobra"

	"clickGo/internal/clickup"
)

var (
	authClientID     string
	authClientSecret string
	authPort         int
)

func randomState() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return fmt.Sprintf("%d", time.Now().UnixNano())
	}
	return hex.EncodeToString(b)
}

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Autentica com o ClickUp via OAuth",
	Long: `Autentica com o ClickUp via OAuth 2.0.

Requer um app OAuth criado em ClickUp > Settings > Apps, com
redirect_uri = http://localhost:8080/callback. Informe o client_id e o
client_secret via flags ou via variáveis de ambiente:

  CLICKUP_CLIENT_ID / CLICKUP_CLIENT_SECRET`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := clickup.LoadConfig()
		if err != nil {
			return err
		}

		cfg.ClientID = firstNonEmpty(authClientID, getenv("CLICKUP_CLIENT_ID"), cfg.ClientID)
		cfg.ClientSecret = firstNonEmpty(authClientSecret, getenv("CLICKUP_CLIENT_SECRET"), cfg.ClientSecret)

		if cfg.ClientID == "" || cfg.ClientSecret == "" {
			return fmt.Errorf("client_id e client_secret ausentes. Crie um app OAuth em ClickUp > Settings > Apps (redirect_uri http://localhost:%d/callback) e informe com --client-id / --client-secret ou as envs CLICKUP_CLIENT_ID / CLICKUP_CLIENT_SECRET", authPort)
		}

		state := randomState()
		redirectURI := fmt.Sprintf("http://localhost:%d/callback", authPort)
		authURL := clickup.AuthURL(cfg.ClientID, redirectURI, state)

		lipgloss.Println(titleStyle.Render("Autorização do ClickUp"))
		lipgloss.Println(mutedStyle.Render("Abra esta URL no navegador para autorizar:"))
		lipgloss.Println("")
		lipgloss.Println(boxStyle.Render(authURL))
		lipgloss.Println("")
		if err := clickup.OpenBrowser(authURL); err != nil {
			lipgloss.Println(mutedStyle.Render("Não foi possível abrir o navegador automaticamente — cole a URL acima manualmente."))
		}

		code, err := clickup.ListenForCode(authPort, state)
		if err != nil {
			return err
		}

		tok, err := clickup.ExchangeToken(cfg.ClientID, cfg.ClientSecret, code, redirectURI)
		if err != nil {
			return err
		}

		cfg.AccessToken = tok.AccessToken
		cfg.TeamID = tok.TeamID
		cfg.TeamName = tok.TeamName

		if err := cfg.Save(); err != nil {
			return err
		}

		lipgloss.Println(successStyle.Render("Autenticado!"))
		if tok.TeamName != "" {
			lipgloss.Println(lipgloss.JoinHorizontal(lipgloss.Top,
				keyStyle.Render("Workspace:"),
				valStyle.Render(fmt.Sprintf("%s (%s)", tok.TeamName, tok.TeamID)),
			))
		}
		return nil
	},
}

var authStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Mostra o estado da autenticação",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := clickup.LoadConfig()
		if err != nil {
			return err
		}

		if cfg.HasToken() {
			path, _ := clickup.ConfigPath()
			lipgloss.Println(successStyle.Render("Autenticado"))
			if cfg.TeamName != "" {
				lipgloss.Println(lipgloss.JoinHorizontal(lipgloss.Top,
					keyStyle.Render("Workspace:"),
					valStyle.Render(fmt.Sprintf("%s (%s)", cfg.TeamName, cfg.TeamID)),
				))
			}
			lipgloss.Println(mutedStyle.Render("Token salvo em: " + path))
		} else {
			lipgloss.Println(errStyle.Render("Não autenticado"))
			lipgloss.Println(mutedStyle.Render("Rode: clickGo auth"))
		}
		return nil
	},
}

var meCmd = &cobra.Command{
	Use:   "me",
	Short: "Exibe o usuário autenticado (testa o token)",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := clickup.LoadConfig()
		if err != nil {
			return err
		}
		if !cfg.HasToken() {
			return fmt.Errorf("sem token. Rode: clickGo auth")
		}

		client := clickup.New(cfg.AccessToken)
		user, err := client.GetUser()
		if err != nil {
			return err
		}

		rows := lipgloss.JoinVertical(lipgloss.Left,
			titleStyle.Render("Usuário autenticado"),
			"",
			lipgloss.JoinHorizontal(lipgloss.Top, keyStyle.Render("ID:"), valStyle.Render(fmt.Sprintf("%d", user.ID))),
			lipgloss.JoinHorizontal(lipgloss.Top, keyStyle.Render("Usuário:"), valStyle.Render(user.Username)),
			lipgloss.JoinHorizontal(lipgloss.Top, keyStyle.Render("E-mail:"), valStyle.Render(user.Email)),
			lipgloss.JoinHorizontal(lipgloss.Top, keyStyle.Render("Iniciais:"), valStyle.Render(user.Initials)),
		)
		lipgloss.Println(boxStyle.Render(rows))
		return nil
	},
}

func init() {
	authCmd.Flags().StringVar(&authClientID, "client-id", "", "Client ID do app OAuth do ClickUp")
	authCmd.Flags().StringVar(&authClientSecret, "client-secret", "", "Client Secret do app OAuth do ClickUp")
	authCmd.Flags().IntVar(&authPort, "port", 8080, "Porta do servidor local de callback")

	authCmd.AddCommand(authStatusCmd)
	rootCmd.AddCommand(authCmd, meCmd)
}
