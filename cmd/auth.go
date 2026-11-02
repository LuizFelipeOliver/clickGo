package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/spf13/cobra"

	"clickGo/internal/clickup"
)

var (
	authTokenFlag     string
	authNoBrowserFlag bool
)

// loadSession carrega a configuração, garante que UserID e TeamID estejam preenchidos e retorna o Client.
func loadSession() (*clickup.Client, *clickup.Config, error) {
	cfg, err := clickup.LoadConfig()
	if err != nil {
		return nil, nil, err
	}
	if !cfg.HasToken() {
		return nil, nil, fmt.Errorf("não autenticado. Execute: clickGo auth")
	}

	client := clickup.NewClient(cfg.Token)
	updated := false

	if cfg.UserID == 0 {
		user, err := client.GetUser()
		if err != nil {
			return nil, nil, fmt.Errorf("falha ao obter usuário: %w", err)
		}
		cfg.UserID = user.ID
		cfg.Username = user.Username
		updated = true
	}

	if cfg.TeamID == "" {
		teams, err := client.GetTeams()
		if err != nil {
			return nil, nil, fmt.Errorf("falha ao obter workspaces: %w", err)
		}
		if len(teams) == 0 {
			return nil, nil, fmt.Errorf("nenhum Workspace encontrado nesta conta do ClickUp")
		}
		cfg.TeamID = teams[0].ID
		cfg.TeamName = teams[0].Name
		updated = true
	}

	if updated {
		_ = cfg.Save()
	}

	return client, cfg, nil
}

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Autentica na API do ClickUp usando seu Personal API Token",
	Long: `Abre a página de configurações do ClickUp no navegador para você copiar
seu Personal API Token (pk_...) e salvá-lo localmente na CLI.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		token := strings.TrimSpace(authTokenFlag)

		if token == "" {
			lipgloss.Println(titleStyle.Render("Autenticação no ClickUp"))
			lipgloss.Println(mutedStyle.Render("Copie ou gere seu Personal API Token (pk_...) no link abaixo:"))
			lipgloss.Println("")
			lipgloss.Println(boxStyle.Render(clickup.TokenSettingsURL))
			lipgloss.Println("")

			if !authNoBrowserFlag {
				if err := clickup.OpenBrowser(clickup.TokenSettingsURL); err != nil {
					lipgloss.Println(mutedStyle.Render("Não foi possível abrir o navegador automaticamente. Abra o link acima manualmente."))
				} else {
					lipgloss.Println(mutedStyle.Render("Navegador aberto!"))
				}
			}

			fmt.Print("\nCole seu token (pk_...): ")
			reader := bufio.NewReader(os.Stdin)
			input, err := reader.ReadString('\n')
			if err != nil {
				return fmt.Errorf("erro ao ler token: %w", err)
			}
			token = strings.TrimSpace(input)
		}

		if token == "" {
			return fmt.Errorf("nenhum token informado")
		}

		lipgloss.Println(mutedStyle.Render("Validando token no ClickUp..."))
		client := clickup.NewClient(token)
		user, err := client.GetUser()
		if err != nil {
			return fmt.Errorf("token inválido ou falha ao conectar: %w", err)
		}

		teams, err := client.GetTeams()
		if err != nil {
			return fmt.Errorf("erro ao buscar workspaces: %w", err)
		}

		cfg, err := clickup.LoadConfig()
		if err != nil {
			return err
		}
		cfg.Token = token
		cfg.UserID = user.ID
		cfg.Username = user.Username
		if len(teams) > 0 {
			cfg.TeamID = teams[0].ID
			cfg.TeamName = teams[0].Name
		}

		if err := cfg.Save(); err != nil {
			return fmt.Errorf("erro ao salvar configuração: %w", err)
		}

		lipgloss.Println("")
		lipgloss.Println(successStyle.Render("✓ Autenticado com sucesso!"))
		lipgloss.Println(lipgloss.JoinHorizontal(lipgloss.Top,
			keyStyle.Render("Usuário:"),
			valStyle.Render(fmt.Sprintf("%s (%s)", user.Username, user.Email)),
		))
		if cfg.TeamName != "" {
			lipgloss.Println(lipgloss.JoinHorizontal(lipgloss.Top,
				keyStyle.Render("Workspace:"),
				valStyle.Render(fmt.Sprintf("%s (%s)", cfg.TeamName, cfg.TeamID)),
			))
		}
		return nil
	},
}

var authStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Verifica o status da autenticação atual",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, cfg, err := loadSession()
		if err != nil {
			lipgloss.Println(errStyle.Render("✗ " + err.Error()))
			return nil
		}

		user, err := client.GetUser()
		if err != nil {
			lipgloss.Println(errStyle.Render("✗ Token salvo é inválido ou expirou"))
			return err
		}

		path, _ := clickup.ConfigPath()
		rows := lipgloss.JoinVertical(lipgloss.Left,
			successStyle.Render("✓ Autenticado"),
			"",
			lipgloss.JoinHorizontal(lipgloss.Top, keyStyle.Render("ID:"), valStyle.Render(fmt.Sprintf("%d", user.ID))),
			lipgloss.JoinHorizontal(lipgloss.Top, keyStyle.Render("Usuário:"), valStyle.Render(user.Username)),
			lipgloss.JoinHorizontal(lipgloss.Top, keyStyle.Render("E-mail:"), valStyle.Render(user.Email)),
			lipgloss.JoinHorizontal(lipgloss.Top, keyStyle.Render("Workspace:"), valStyle.Render(fmt.Sprintf("%s (%s)", cfg.TeamName, cfg.TeamID))),
			lipgloss.JoinHorizontal(lipgloss.Top, keyStyle.Render("Config:"), mutedStyle.Render(path)),
		)
		lipgloss.Println(boxStyle.Render(rows))
		return nil
	},
}

var authLogoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Remove o token salvo localmente",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := &clickup.Config{}
		if err := cfg.Save(); err != nil {
			return err
		}
		lipgloss.Println(successStyle.Render("✓ Token removido com sucesso."))
		return nil
	},
}

func init() {
	authCmd.Flags().StringVarP(&authTokenFlag, "token", "t", "", "Personal API Token (pk_...)")
	authCmd.Flags().BoolVar(&authNoBrowserFlag, "no-browser", false, "Não abrir o navegador automaticamente")

	authCmd.AddCommand(authStatusCmd, authLogoutCmd)
	rootCmd.AddCommand(authCmd)
}
