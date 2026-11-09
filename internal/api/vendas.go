package api

import (
	"context"
	"fmt"
	"strconv"

	"github.com/UrielAndrade/TUI-StarConcessionaria/internal/model"
)

// ListVendas obtém histórico de vendas com filtros.
func (c *Client) ListVendas(ctx context.Context, f model.FiltroVendas) (*model.PagedResult[model.VendaResponse], error) {
	req := c.http.R().SetContext(ctx)

	if f.UnidadeID != nil && *f.UnidadeID > 0 {
		req.SetQueryParam("UnidadeId", strconv.Itoa(*f.UnidadeID))
	}
	if f.VendedorID != nil && *f.VendedorID > 0 {
		req.SetQueryParam("VendedorId", strconv.Itoa(*f.VendedorID))
	}
	if f.PeriodoInicio != "" {
		req.SetQueryParam("PeriodoInicio", f.PeriodoInicio)
	}
	if f.PeriodoFim != "" {
		req.SetQueryParam("PeriodoFim", f.PeriodoFim)
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

	var res model.PagedResult[model.VendaResponse]
	resp, err := req.
		SetResult(&res).
		Get("/api/v1/vendas")

	if err != nil {
		return nil, err
	}

	if err := HandleResponseError(resp); err != nil {
		return nil, err
	}

	return &res, nil
}

// RegistrarVenda registra uma nova venda de veículo.
func (c *Client) RegistrarVenda(ctx context.Context, payload model.RegistrarVendaRequest) error {
	resp, err := c.http.R().
		SetContext(ctx).
		SetBody(payload).
		Post("/api/v1/vendas")

	if err != nil {
		return err
	}

	return HandleResponseError(resp)
}

// ListAgendamentos busca agendamentos de visita e test-drives.
func (c *Client) ListAgendamentos(ctx context.Context, f model.FiltroAgendamentos) (*model.PagedResult[model.AgendamentoResponse], error) {
	req := c.http.R().SetContext(ctx)

	if f.UnidadeID != nil && *f.UnidadeID > 0 {
		req.SetQueryParam("UnidadeId", strconv.Itoa(*f.UnidadeID))
	}
	if f.Status != "" {
		req.SetQueryParam("Status", f.Status)
	}
	if f.Tipo != "" {
		req.SetQueryParam("Tipo", f.Tipo)
	}
	if f.Data != "" {
		req.SetQueryParam("Data", f.Data)
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

	var res model.PagedResult[model.AgendamentoResponse]
	resp, err := req.
		SetResult(&res).
		Get("/api/v1/agendamentos")

	if err != nil {
		return nil, err
	}

	if err := HandleResponseError(resp); err != nil {
		return nil, err
	}

	return &res, nil
}

// AtualizarStatusAgendamento altera status de agendamento.
func (c *Client) AtualizarStatusAgendamento(ctx context.Context, id int, status, obs string) error {
	payload := model.AtualizarStatusAgendamentoRequest{
		Status:      status,
		Observacoes: obs,
	}

	resp, err := c.http.R().
		SetContext(ctx).
		SetBody(payload).
		Patch(fmt.Sprintf("/api/v1/agendamentos/%d/status", id))

	if err != nil {
		return err
	}

	return HandleResponseError(resp)
}
