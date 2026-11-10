package api

import (
	"context"
	"fmt"
	"strconv"

	"github.com/UrielAndrade/TUI-StarConcessionaria/internal/model"
)

// ListContatos lista mensagens de atendimento e contatos com paginação e filtros.
func (c *Client) ListContatos(ctx context.Context, f model.FiltroContatos) (*model.PagedResult[model.ContatoResponse], error) {
	req := c.http.R().SetContext(ctx)

	if f.Tipo != "" {
		req.SetQueryParam("Tipo", f.Tipo)
	}
	if f.Respondido != nil {
		req.SetQueryParam("Respondido", strconv.FormatBool(*f.Respondido))
	}
	if f.UnidadeID != nil && *f.UnidadeID > 0 {
		req.SetQueryParam("UnidadeId", strconv.Itoa(*f.UnidadeID))
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

	var res model.PagedResult[model.ContatoResponse]
	resp, err := req.
		SetResult(&res).
		Get("/api/v1/contatos")

	if err != nil {
		return nil, err
	}

	if err := HandleResponseError(resp); err != nil {
		return nil, err
	}

	return &res, nil
}

// NovoContato envia uma mensagem geral para o SAC da concessionária.
func (c *Client) NovoContato(ctx context.Context, payload model.NovoContatoRequest) error {
	resp, err := c.http.R().
		SetContext(ctx).
		SetBody(payload).
		Post("/api/v1/contatos")

	if err != nil {
		return err
	}

	return HandleResponseError(resp)
}

// ContatoVendedor envia mensagem direcionada para a equipe de vendas referente a um veículo.
func (c *Client) ContatoVendedor(ctx context.Context, payload model.ContatoVendedorRequest) error {
	resp, err := c.http.R().
		SetContext(ctx).
		SetBody(payload).
		Post("/api/v1/contatos/vendedor")

	if err != nil {
		return err
	}

	return HandleResponseError(resp)
}

// ResponderContato registra a resposta ou anotação de atendimento do chamado.
func (c *Client) ResponderContato(ctx context.Context, id int, notasResposta string) error {
	payload := model.ResponderContatoRequest{
		NotasResposta: notasResposta,
	}

	resp, err := c.http.R().
		SetContext(ctx).
		SetBody(payload).
		Patch(fmt.Sprintf("/api/v1/contatos/%d/responder", id))

	if err != nil {
		return err
	}

	return HandleResponseError(resp)
}
