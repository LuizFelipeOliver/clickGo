package clickup

import (
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestAuthURL(t *testing.T) {
	u := AuthURL("cid-123", "http://localhost:8080/callback", "state-abc")
	for _, want := range []string{
		"https://app.clickup.com/api",
		"client_id=cid-123",
		"redirect_uri=http%3A%2F%2Flocalhost%3A8080%2Fcallback",
		"response_type=code",
		"state=state-abc",
	} {
		if !strings.Contains(u, want) {
			t.Errorf("AuthURL() = %q, quer conter %q", u, want)
		}
	}
}

func TestListenForCode(t *testing.T) {
	const (
		port  = 9091
		state = "test-state"
	)
	done := make(chan string, 1)
	go func() {
		code, err := ListenForCode(port, state)
		if err != nil {
			done <- "err: " + err.Error()
			return
		}
		done <- code
	}()

	time.Sleep(300 * time.Millisecond)
	resp, err := http.Get("http://localhost:9091/callback?code=abc123&state=" + state)
	if err != nil {
		t.Fatalf("GET callback: %v", err)
	}
	resp.Body.Close()

	select {
	case got := <-done:
		if got != "abc123" {
			t.Errorf("ListenForCode() = %q, esperado abc123", got)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("timeout: callback não capturado")
	}
}

func TestListenForCodeRejectsWrongState(t *testing.T) {
	const port = 9092
	go func() {
		ListenForCode(port, "esperado") //nolint:errcheck
	}()

	time.Sleep(300 * time.Millisecond)
	resp, err := http.Get("http://localhost:9092/callback?code=abc123&state=errado")
	if err != nil {
		t.Fatalf("GET callback: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("status = %d, esperado 400", resp.StatusCode)
	}
}
