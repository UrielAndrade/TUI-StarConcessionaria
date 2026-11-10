package api

import (
	"context"
	"fmt"
	"strconv"

	"github.com/UrielAndrade/TUI-StarConcessionaria/internal/model"
)

// ListPropostas consulta propostas comerciais de clientes com filtros.
func (c *Client) ListPropostas(ctx context.Context, f model.FiltroPropostas) (*model.PagedResult[model.PropostaResponse], error) {
	req := c.http.R().SetContext(ctx)

	if f.Respondido != nil {
		req.SetQueryParam("respondido", strconv.FormatBool(*f.Respondido))
	}
	if f.UnidadeID != nil && *f.UnidadeID > 0 {
		req.SetQueryParam("unidadeId", strconv.Itoa(*f.UnidadeID))
	}

	pagina := f.Pagina
	if pagina < 1 {
		pagina = 1
	}
	itensPorPagina := f.ItensPorPagina
	if itensPorPagina < 1 {
		itensPorPagina = 12
	}

	req.SetQueryParam("pagina", strconv.Itoa(pagina))
	req.SetQueryParam("itensPorPagina", strconv.Itoa(itensPorPagina))

	var res model.PagedResult[model.PropostaResponse]
	resp, err := req.
		SetResult(&res).
		Get("/api/v1/propostas")

	if err != nil {
		return nil, err
	}

	if err := HandleResponseError(resp); err != nil {
		return nil, err
	}

	return &res, nil
}

// SolicitarProposta envia proposta formal de compra com veículo de troca ou entrada.
func (c *Client) SolicitarProposta(ctx context.Context, payload model.SolicitarPropostaRequest) error {
	resp, err := c.http.R().
		SetContext(ctx).
		SetBody(payload).
		Post("/api/v1/propostas")

	if err != nil {
		return err
	}

	return HandleResponseError(resp)
}

// AtualizarStatusProposta responde e atualiza o andamento da proposta.
func (c *Client) AtualizarStatusProposta(ctx context.Context, id int, notasResposta string) error {
	payload := model.ResponderContatoRequest{
		NotasResposta: notasResposta,
	}

	resp, err := c.http.R().
		SetContext(ctx).
		SetBody(payload).
		Patch(fmt.Sprintf("/api/v1/propostas/%d/status", id))

	if err != nil {
		return err
	}

	return HandleResponseError(resp)
}
