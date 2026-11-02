package clickup

import (
	"fmt"
	"os/exec"
	"runtime"
)

const TokenSettingsURL = "https://app.clickup.com/settings/apps"

// OpenBrowser abre a URL informada no navegador padrão do sistema operacional.
func OpenBrowser(url string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("abrir navegador: %w", err)
	}
	return nil
}
