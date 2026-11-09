package api

import (
	"context"
	"sync"
	"time"

	"github.com/UrielAndrade/TUI-StarConcessionaria/internal/config"
	"github.com/go-resty/resty/v2"
)

// Client gerencia a comunicação HTTP com a API da Star Concessionária.
type Client struct {
	http       *resty.Client
	baseURL    string
	tokenMu    sync.RWMutex
	token      string
	timeout    time.Duration
	coldWait   time.Duration
}

// NewClient instancia um cliente configurado.
func NewClient(cfg *config.Config) *Client {
	r := resty.New()
	r.SetBaseURL(cfg.BaseURL)
	r.SetTimeout(cfg.Timeout)
	r.SetHeader("Accept", "application/json")
	r.SetHeader("User-Agent", "TUI-StarConcessionaria/1.0.0")

	c := &Client{
		http:     r,
		baseURL:  cfg.BaseURL,
		token:    cfg.Token,
		timeout:  cfg.Timeout,
		coldWait: cfg.ColdStartTimeout,
	}

	// Middleware para inclusão dinâmica do Bearer Token
	r.OnBeforeRequest(func(client *resty.Client, req *resty.Request) error {
		c.tokenMu.RLock()
		tok := c.token
		c.tokenMu.RUnlock()

		if tok != "" && req.Header.Get("Authorization") == "" {
			req.SetAuthToken(tok)
		}
		return nil
	})

	return c
}

// SetToken atualiza o token JWT da sessão.
func (c *Client) SetToken(token string) {
	c.tokenMu.Lock()
	defer c.tokenMu.Unlock()
	c.token = token
}

// GetToken obtém o token JWT atual.
func (c *Client) GetToken() string {
	c.tokenMu.RLock()
	defer c.tokenMu.RUnlock()
	return c.token
}

// Ping verifica a conectividade com o servidor, permitindo despertar o Render.
func (c *Client) Ping(ctx context.Context) error {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, c.coldWait)
	defer cancel()

	resp, err := c.http.R().
		SetContext(ctxWithTimeout).
		Get("/")
	if err != nil {
		return err
	}
	return HandleResponseError(resp)
}
