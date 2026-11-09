package api

import (
	"context"
	"fmt"
	"strconv"

	"github.com/UrielAndrade/TUI-StarConcessionaria/internal/model"
)

// ListVeiculos retorna página de veículos segundo os filtros informados.
func (c *Client) ListVeiculos(ctx context.Context, f model.FiltroVeiculos) (*model.PagedResult[model.VeiculoCatalogoResponse], error) {
	req := c.http.R().SetContext(ctx)

	if f.Marca != "" {
		req.SetQueryParam("Marca", f.Marca)
	}
	if f.Modelo != "" {
		req.SetQueryParam("Modelo", f.Modelo)
	}
	if f.Tipo != "" {
		req.SetQueryParam("Tipo", f.Tipo)
	}
	if f.PrecoMin != nil {
		req.SetQueryParam("PrecoMin", fmt.Sprintf("%.2f", *f.PrecoMin))
	}
	if f.PrecoMax != nil {
		req.SetQueryParam("PrecoMax", fmt.Sprintf("%.2f", *f.PrecoMax))
	}
	if f.AnoMin != nil {
		req.SetQueryParam("AnoMin", strconv.Itoa(*f.AnoMin))
	}
	if f.AnoMax != nil {
		req.SetQueryParam("AnoMax", strconv.Itoa(*f.AnoMax))
	}
	if f.UnidadeID != nil && *f.UnidadeID > 0 {
		req.SetQueryParam("UnidadeId", strconv.Itoa(*f.UnidadeID))
	}
	if f.SomenteDisponivel != nil {
		req.SetQueryParam("SomenteDisponiveis", strconv.FormatBool(*f.SomenteDisponivel))
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

	var res model.PagedResult[model.VeiculoCatalogoResponse]
	resp, err := req.
		SetResult(&res).
		Get("/api/v1/veiculos")

	if err != nil {
		return nil, err
	}

	if err := HandleResponseError(resp); err != nil {
		return nil, err
	}

	return &res, nil
}

// GetVeiculo retorna a ficha completa de um veículo específico.
func (c *Client) GetVeiculo(ctx context.Context, id int) (*model.VeiculoDetalheResponse, error) {
	var res model.VeiculoDetalheResponse
	resp, err := c.http.R().
		SetContext(ctx).
		SetResult(&res).
		Get(fmt.Sprintf("/api/v1/veiculos/%d", id))

	if err != nil {
		return nil, err
	}

	if err := HandleResponseError(resp); err != nil {
		return nil, err
	}

	return &res, nil
}

// CriarVeiculo cadastra um novo veículo no estoque.
func (c *Client) CriarVeiculo(ctx context.Context, payload model.SalvarVeiculoRequest) (*model.VeiculoDetalheResponse, error) {
	var res model.VeiculoDetalheResponse
	resp, err := c.http.R().
		SetContext(ctx).
		SetBody(payload).
		SetResult(&res).
		Post("/api/v1/veiculos")

	if err != nil {
		return nil, err
	}

	if err := HandleResponseError(resp); err != nil {
		return nil, err
	}

	return &res, nil
}

// AtualizarVeiculo altera dados de um veículo existente.
func (c *Client) AtualizarVeiculo(ctx context.Context, id int, payload model.SalvarVeiculoRequest) error {
	resp, err := c.http.R().
		SetContext(ctx).
		SetBody(payload).
		Put(fmt.Sprintf("/api/v1/veiculos/%d", id))

	if err != nil {
		return err
	}

	return HandleResponseError(resp)
}

// AlterarSituacaoVeiculo atualiza o status operacional do veículo.
func (c *Client) AlterarSituacaoVeiculo(ctx context.Context, id int, situacao, justificativa string) error {
	payload := model.AlterarSituacaoRequest{
		Situacao:      situacao,
		Justificativa: justificativa,
	}

	resp, err := c.http.R().
		SetContext(ctx).
		SetBody(payload).
		Patch(fmt.Sprintf("/api/v1/veiculos/%d/situacao", id))

	if err != nil {
		return err
	}

	return HandleResponseError(resp)
}
