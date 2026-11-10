package api

import (
	"context"
	"fmt"
	"strconv"

	"github.com/UrielAndrade/TUI-StarConcessionaria/internal/model"
)

// ListReservas lista as reservas efetuadas com paginação e filtros.
func (c *Client) ListReservas(ctx context.Context, f model.FiltroReservas) (*model.PagedResult[model.ReservaResponse], error) {
	req := c.http.R().SetContext(ctx)

	if f.Status != "" {
		req.SetQueryParam("Status", f.Status)
	}
	if f.VeiculoID != nil && *f.VeiculoID > 0 {
		req.SetQueryParam("VeiculoId", strconv.Itoa(*f.VeiculoID))
	}
	if f.DataInicio != "" {
		req.SetQueryParam("DataInicio", f.DataInicio)
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

	var res model.PagedResult[model.ReservaResponse]
	resp, err := req.
		SetResult(&res).
		Get("/api/v1/reservas")

	if err != nil {
		return nil, err
	}

	if err := HandleResponseError(resp); err != nil {
		return nil, err
	}

	return &res, nil
}

// GetReserva consulta os detalhes de uma reserva específica.
func (c *Client) GetReserva(ctx context.Context, id int) (*model.ReservaResponse, error) {
	var res model.ReservaResponse
	resp, err := c.http.R().
		SetContext(ctx).
		SetResult(&res).
		Get(fmt.Sprintf("/api/v1/reservas/%d", id))

	if err != nil {
		return nil, err
	}

	if err := HandleResponseError(resp); err != nil {
		return nil, err
	}

	return &res, nil
}

// SolicitarReserva bloqueia um veículo para um interessado por período determinado.
func (c *Client) SolicitarReserva(ctx context.Context, payload model.SolicitarReservaRequest) error {
	resp, err := c.http.R().
		SetContext(ctx).
		SetBody(payload).
		Post("/api/v1/reservas")

	if err != nil {
		return err
	}

	return HandleResponseError(resp)
}

// AtualizarStatusReserva atualiza situação (Ativa, Cancelada, Concluída) com justificativa.
func (c *Client) AtualizarStatusReserva(ctx context.Context, id int, novoStatus, motivo string) error {
	payload := model.AtualizarStatusReservaRequest{
		NovoStatus: novoStatus,
		Motivo:     motivo,
	}

	resp, err := c.http.R().
		SetContext(ctx).
		SetBody(payload).
		Patch(fmt.Sprintf("/api/v1/reservas/%d/status", id))

	if err != nil {
		return err
	}

	return HandleResponseError(resp)
}
