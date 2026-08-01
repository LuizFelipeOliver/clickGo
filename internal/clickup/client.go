package clickup

import (
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
)

const baseURL = "https://api.clickup.com/api/v2"

// Client encapsula os chamados à API do ClickUp.
type Client struct {
	http  *resty.Client
	token string
}

// New cria um Client autenticado com o token informado.
func New(token string) *Client {
	return &Client{
		token: token,
		http: resty.New().
			SetBaseURL(baseURL).
			SetHeader("Authorization", "Bearer "+token),
	}
}

// do executa uma requisição, decodifica a resposta JSON em out e trata erros.
func (c *Client) do(method, path string, query map[string]string, body, out any) error {
	req := c.http.R()
	if query != nil {
		req.SetQueryParams(query)
	}
	if body != nil {
		req.SetBody(body)
	}
	if out != nil {
		req.SetResult(out)
	}

	resp, err := req.Execute(method, path)
	if err != nil {
		return fmt.Errorf("%s %s: %w", method, path, err)
	}
	if resp.IsError() {
		return fmt.Errorf("%s %s: HTTP %d: %s", method, path, resp.StatusCode(), string(resp.Body()))
	}
	return nil
}

// User é a representação do usuário autenticado.
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Initials string `json:"initials"`
}

// GetUser retorna o usuário autenticado (GET /user).
func (c *Client) GetUser() (*User, error) {
	var out User
	if err := c.do(http.MethodGet, "/user", nil, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
