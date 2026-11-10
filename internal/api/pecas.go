package api

import (
	"context"
	"fmt"
	"strconv"

	"github.com/UrielAndrade/TUI-StarConcessionaria/internal/model"
)

// ListPecas pesquisa catálogo de peças com paginação.
func (c *Client) ListPecas(ctx context.Context, f model.FiltroPecas) (*model.PagedResult[model.PecaResponse], error) {
	req := c.http.R().SetContext(ctx)

	if f.Busca != "" {
		req.SetQueryParam("Busca", f.Busca)
	}
	if f.Disponivel != nil {
		req.SetQueryParam("Disponivel", strconv.FormatBool(*f.Disponivel))
	}

	pagina := f.Pagina
	if pagina < 1 {
		pagina = 1
	}
	itensPorPagina := f.ItensPorPagina
	if itensPorPagina < 1 {
		itensPorPagina = 12
	}

	req.SetQueryParam("Pagina", strconv.Itoa(pagina))
	req.SetQueryParam("ItensPorPagina", strconv.Itoa(itensPorPagina))

	var res model.PagedResult[model.PecaResponse]
	resp, err := req.
		SetResult(&res).
		Get("/api/v1/pecas")

	if err != nil {
		return nil, err
	}

	if err := HandleResponseError(resp); err != nil {
		return nil, err
	}

	return &res, nil
}

// GetPeca busca os detalhes de uma peça específica pelo seu identificador.
func (c *Client) GetPeca(ctx context.Context, id int) (*model.PecaResponse, error) {
	var res model.PecaResponse
	resp, err := c.http.R().
		SetContext(ctx).
		SetResult(&res).
		Get(fmt.Sprintf("/api/v1/pecas/%d", id))

	if err != nil {
		return nil, err
	}

	if err := HandleResponseError(resp); err != nil {
		return nil, err
	}

	return &res, nil
}

// CriarPeca cadastra uma nova peça automotiva.
func (c *Client) CriarPeca(ctx context.Context, payload model.SalvarPecaRequest) (*model.PecaResponse, error) {
	var res model.PecaResponse
	resp, err := c.http.R().
		SetContext(ctx).
		SetBody(payload).
		SetResult(&res).
		Post("/api/v1/pecas")

	if err != nil {
		return nil, err
	}

	if err := HandleResponseError(resp); err != nil {
		return nil, err
	}

	return &res, nil
}

// AtualizarPeca altera os dados cadastrais ou estoque de uma peça existente.
func (c *Client) AtualizarPeca(ctx context.Context, id int, payload model.SalvarPecaRequest) error {
	resp, err := c.http.R().
		SetContext(ctx).
		SetBody(payload).
		Put(fmt.Sprintf("/api/v1/pecas/%d", id))

	if err != nil {
		return err
	}

	return HandleResponseError(resp)
}

// VenderPeca registra a baixa por venda de estoque da peça.
func (c *Client) VenderPeca(ctx context.Context, pecaID int, payload model.RegistrarVendaPecaRequest) error {
	resp, err := c.http.R().
		SetContext(ctx).
		SetBody(payload).
		Post(fmt.Sprintf("/api/v1/pecas/%d/vendas", pecaID))

	if err != nil {
		return err
	}

	return HandleResponseError(resp)
}
