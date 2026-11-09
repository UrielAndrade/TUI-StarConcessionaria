package api

import (
	"context"

	"github.com/UrielAndrade/TUI-StarConcessionaria/internal/model"
)

// ListUnidades obtém a lista de todas as concessionárias cadastradas.
func (c *Client) ListUnidades(ctx context.Context) ([]model.UnidadeResponse, error) {
	var res []model.UnidadeResponse
	resp, err := c.http.R().
		SetContext(ctx).
		SetResult(&res).
		Get("/api/v1/unidades")

	if err != nil {
		return nil, err
	}

	if err := HandleResponseError(resp); err != nil {
		return nil, err
	}

	return res, nil
}
