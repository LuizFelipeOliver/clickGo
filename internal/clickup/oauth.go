package clickup

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os/exec"
	"time"

	"github.com/go-resty/resty/v2"
)

const (
	oauthAuthURL  = "https://app.clickup.com/api"
	oauthTokenURL = "https://api.clickup.com/api/v2/oauth/token"
)

// AuthURL monta a URL de autorização para abrir no navegador.
func AuthURL(clientID, redirectURI, state string) string {
	q := url.Values{}
	q.Set("client_id", clientID)
	q.Set("redirect_uri", redirectURI)
	q.Set("response_type", "code")
	q.Set("state", state)
	return oauthAuthURL + "?" + q.Encode()
}

// tokenResponse é a resposta da troca do code por access token.
type tokenResponse struct {
	AccessToken string `json:"access_token"`
	TeamID      string `json:"team_id"`
	TeamName    string `json:"team_name"`
	Team        *struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"team"`
}

// ExchangeToken troca o code de autorização por um access token.
func ExchangeToken(clientID, clientSecret, code, redirectURI string) (*tokenResponse, error) {
	req := resty.New().R().SetQueryParams(map[string]string{
		"client_id":     clientID,
		"client_secret": clientSecret,
		"code":          code,
	})
	if redirectURI != "" {
		req.SetQueryParam("redirect_uri", redirectURI)
	}

	var out tokenResponse
	resp, err := req.SetResult(&out).Post(oauthTokenURL)
	if err != nil {
		return nil, fmt.Errorf("troca de token: %w", err)
	}
	if resp.IsError() {
		return nil, fmt.Errorf("troca de token: HTTP %d: %s", resp.StatusCode(), string(resp.Body()))
	}

	if out.Team != nil && out.TeamID == "" {
		out.TeamID = out.Team.ID
		out.TeamName = out.Team.Name
	}
	return &out, nil
}

// ListenForCode sobe um servidor local em localhost:port/callback e aguarda o
// redirecionamento do ClickUp com o code de autorização.
func ListenForCode(port int, expectedState string) (string, error) {
	codeCh := make(chan string, 1)

	ln, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		return "", fmt.Errorf("servidor local: %w", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("state") != expectedState {
			http.Error(w, "state inválido", http.StatusBadRequest)
			return
		}
		code := q.Get("code")
		if code == "" {
			http.Error(w, "code ausente", http.StatusBadRequest)
			return
		}
		codeCh <- code
		fmt.Fprintln(w, "Autenticação concluída! Você já pode fechar esta aba.")
	})

	srv := &http.Server{Handler: mux}
	go srv.Serve(ln)
	defer srv.Close()

	select {
	case code := <-codeCh:
		return code, nil
	case <-time.After(5 * time.Minute):
		return "", errors.New("timeout aguardando a autorização no navegador")
	}
}

// OpenBrowser abre a URL no navegador padrão (xdg-open no Linux).
func OpenBrowser(url string) error {
	if err := exec.Command("xdg-open", url).Start(); err != nil {
		return fmt.Errorf("abrir navegador: %w", err)
	}
	return nil
}
