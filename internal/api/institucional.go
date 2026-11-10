package api

import (
	"context"

	"github.com/UrielAndrade/TUI-StarConcessionaria/internal/model"
)

// ListFormasPagamento consulta as opções de pagamento aceitas pela concessionária.
func (c *Client) ListFormasPagamento(ctx context.Context) ([]model.FormaPagamentoDetalheResponse, error) {
	var res []model.FormaPagamentoDetalheResponse
	resp, err := c.http.R().
		SetContext(ctx).
		SetResult(&res).
		Get("/api/v1/institucional/formas-pagamento")

	if err != nil {
		return nil, err
	}

	if err := HandleResponseError(resp); err != nil {
		return nil, err
	}

	return res, nil
}

// ListFaqs pesquisa a base de dúvidas frequentes com filtro textual opcional.
func (c *Client) ListFaqs(ctx context.Context, busca string) ([]model.FaqResponse, error) {
	req := c.http.R().SetContext(ctx)

	if busca != "" {
		req.SetQueryParam("busca", busca)
	}

	var res []model.FaqResponse
	resp, err := req.
		SetResult(&res).
		Get("/api/v1/faqs")

	if err != nil {
		return nil, err
	}

	if err := HandleResponseError(resp); err != nil {
		return nil, err
	}

	return res, nil
}
