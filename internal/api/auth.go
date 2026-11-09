package api

import (
	"context"

	"github.com/UrielAndrade/TUI-StarConcessionaria/internal/model"
)

// Login autentica as credenciais e armazena o token recebido.
func (c *Client) Login(ctx context.Context, email, senha string) (*model.LoginResponse, error) {
	reqBody := model.LoginRequest{
		Email: email,
		Senha: senha,
	}

	var res model.LoginResponse
	resp, err := c.http.R().
		SetContext(ctx).
		SetBody(reqBody).
		SetResult(&res).
		Post("/api/v1/auth/login")

	if err != nil {
		return nil, err
	}

	if err := HandleResponseError(resp); err != nil {
		return nil, err
	}

	c.SetToken(res.Token)
	return &res, nil
}

// GetPerfil busca o perfil do usuário logado.
func (c *Client) GetPerfil(ctx context.Context) (*model.UsuarioResponse, error) {
	var res model.UsuarioResponse
	resp, err := c.http.R().
		SetContext(ctx).
		SetResult(&res).
		Get("/api/v1/auth/me")

	if err != nil {
		return nil, err
	}

	if err := HandleResponseError(resp); err != nil {
		return nil, err
	}

	return &res, nil
}
