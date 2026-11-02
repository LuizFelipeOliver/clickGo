package clickup

import (
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
)

const baseURL = "https://api.clickup.com/api/v2"

type Client struct {
	http  *resty.Client
	token string
}

func NewClient(token string) *Client {
	return &Client{
		token: token,
		http: resty.New().
			SetBaseURL(baseURL).
			SetHeader("Authorization", token),
	}
}

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
		return fmt.Errorf("erro na API (%d): %s", resp.StatusCode(), string(resp.Body()))
	}
	return nil
}

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Initials string `json:"initials"`
}

type userResponse struct {
	User User `json:"user"`
}

// GetUser retorna os dados do usuário dono do token (GET /user).
func (c *Client) GetUser() (*User, error) {
	var res userResponse
	if err := c.do(http.MethodGet, "/user", nil, nil, &res); err != nil {
		return nil, err
	}
	return &res.User, nil
}
