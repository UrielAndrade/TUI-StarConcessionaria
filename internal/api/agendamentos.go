package api

import (
	"context"
	"fmt"
	"strconv"

	"github.com/UrielAndrade/TUI-StarConcessionaria/internal/model"
)

// ListAgendamentos busca agendamentos de visita e test-drives com paginação e filtros.
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

// GetHorariosDisponiveis consulta a grade de horários livres para agendamentos.
func (c *Client) GetHorariosDisponiveis(ctx context.Context, unidadeID int, data string, veiculoID *int) ([]string, error) {
	req := c.http.R().SetContext(ctx)
	req.SetQueryParam("unidadeId", strconv.Itoa(unidadeID))
	req.SetQueryParam("data", data)

	if veiculoID != nil && *veiculoID > 0 {
		req.SetQueryParam("veiculoId", strconv.Itoa(*veiculoID))
	}

	var res []string
	resp, err := req.
		SetResult(&res).
		Get("/api/v1/agendamentos/horarios-disponiveis")

	if err != nil {
		return nil, err
	}

	if err := HandleResponseError(resp); err != nil {
		return nil, err
	}

	return res, nil
}

// AgendarVisita registra um novo comparecimento de cliente à loja.
func (c *Client) AgendarVisita(ctx context.Context, payload model.AgendarVisitaRequest) error {
	resp, err := c.http.R().
		SetContext(ctx).
		SetBody(payload).
		Post("/api/v1/agendamentos/visita")

	if err != nil {
		return err
	}

	return HandleResponseError(resp)
}

// AgendarTestDrive registra solicitação de condução de teste de veículo com CNH.
func (c *Client) AgendarTestDrive(ctx context.Context, payload model.AgendarTestDriveRequest) error {
	resp, err := c.http.R().
		SetContext(ctx).
		SetBody(payload).
		Post("/api/v1/agendamentos/test-drive")

	if err != nil {
		return err
	}

	return HandleResponseError(resp)
}

// AtualizarStatusAgendamento altera status de agendamento (Confirmado, Concluído, Cancelado).
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
