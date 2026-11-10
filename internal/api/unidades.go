package api

import (
	"context"
	"fmt"

	"github.com/UrielAndrade/TUI-StarConcessionaria/internal/model"
)

// ListUnidades obtém a lista de todas as concessionárias cadastradas com filtros opcionais.
func (c *Client) ListUnidades(ctx context.Context, cidade, estado string) ([]model.UnidadeResponse, error) {
	req := c.http.R().SetContext(ctx)

	if cidade != "" {
		req.SetQueryParam("cidade", cidade)
	}
	if estado != "" {
		req.SetQueryParam("estado", estado)
	}

	var res []model.UnidadeResponse
	resp, err := req.
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

// GetUnidade obtém uma concessionária específica pelo ID.
func (c *Client) GetUnidade(ctx context.Context, id int) (*model.UnidadeResponse, error) {
	var res model.UnidadeResponse
	resp, err := c.http.R().
		SetContext(ctx).
		SetResult(&res).
		Get(fmt.Sprintf("/api/v1/unidades/%d", id))

	if err != nil {
		return nil, err
	}

	if err := HandleResponseError(resp); err != nil {
		return nil, err
	}

	return &res, nil
}

// CriarUnidade cadastra uma nova unidade/loja da concessionária.
func (c *Client) CriarUnidade(ctx context.Context, payload model.SalvarUnidadeRequest) (*model.UnidadeResponse, error) {
	var res model.UnidadeResponse
	resp, err := c.http.R().
		SetContext(ctx).
		SetBody(payload).
		SetResult(&res).
		Post("/api/v1/unidades")

	if err != nil {
		return nil, err
	}

	if err := HandleResponseError(resp); err != nil {
		return nil, err
	}

	return &res, nil
}

// AtualizarUnidade atualiza dados cadastrais de uma filial existente.
func (c *Client) AtualizarUnidade(ctx context.Context, id int, payload model.SalvarUnidadeRequest) error {
	resp, err := c.http.R().
		SetContext(ctx).
		SetBody(payload).
		Put(fmt.Sprintf("/api/v1/unidades/%d", id))

	if err != nil {
		return err
	}

	return HandleResponseError(resp)
}
