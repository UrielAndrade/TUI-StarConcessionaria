package api

import (
	"context"
	"strconv"

	"github.com/UrielAndrade/TUI-StarConcessionaria/internal/model"
)

// GetIndicadores recupera métricas e consolidados gerenciais.
func (c *Client) GetIndicadores(ctx context.Context, unidadeID *int) (*model.IndicadoresGerenciaisResponse, error) {
	req := c.http.R().SetContext(ctx)

	if unidadeID != nil && *unidadeID > 0 {
		req.SetQueryParam("unidadeId", strconv.Itoa(*unidadeID))
	}

	var res model.IndicadoresGerenciaisResponse
	resp, err := req.
		SetResult(&res).
		Get("/api/v1/dashboard/indicadores")

	if err != nil {
		return nil, err
	}

	if err := HandleResponseError(resp); err != nil {
		return nil, err
	}

	return &res, nil
}
