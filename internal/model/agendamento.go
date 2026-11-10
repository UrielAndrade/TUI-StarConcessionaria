package model

import "time"

// AgendamentoResponse agendamento de visita ou test-drive.
type AgendamentoResponse struct {
	ID               int       `json:"id"`
	Codigo           string    `json:"codigo"`
	Tipo             string    `json:"tipo"`
	UnidadeID        int       `json:"unidadeId"`
	UnidadeNome      string    `json:"unidadeNome"`
	VeiculoID        *int      `json:"veiculoId"`
	VeiculoDescricao *string   `json:"veiculoDescricao"`
	ClienteID        int       `json:"clienteId"`
	NomeCliente      string    `json:"nomeCliente"`
	EmailCliente     string    `json:"emailCliente"`
	TelefoneCliente  string    `json:"telefoneCliente"`
	DataHoraInicio   time.Time `json:"dataHoraInicio"`
	DataHoraFim      time.Time `json:"dataHoraFim"`
	Situacao         string    `json:"situacao"`
}

// ReservaResponse reserva de veículo solicitada por cliente.
type ReservaResponse struct {
	ID               int       `json:"id"`
	Codigo           string    `json:"codigo"`
	VeiculoID        int       `json:"veiculoId"`
	VeiculoDescricao string    `json:"veiculoDescricao"`
	ClienteID        int       `json:"clienteId"`
	NomeCliente      string    `json:"nomeCliente"`
	EmailCliente     string    `json:"emailCliente"`
	TelefoneCliente  string    `json:"telefoneCliente"`
	DataSolicitacao  time.Time `json:"dataSolicitacao"`
	Situacao         string    `json:"situacao"`
}

// AtualizarStatusAgendamentoRequest payload para mudar status de visita/test-drive.
type AtualizarStatusAgendamentoRequest struct {
	Status      string `json:"status"`
	Observacoes string `json:"observacoes"`
}

// AtualizarStatusReservaRequest payload para aprovar/cancelar reserva.
type AtualizarStatusReservaRequest struct {
	NovoStatus string `json:"novoStatus"`
	Motivo     string `json:"motivo"`
}

// FiltroAgendamentos parâmetros de busca de agendamentos.
type FiltroAgendamentos struct {
	UnidadeID      *int
	Status         string
	Tipo           string
	Data           string
	Pagina         int
	ItensPorPagina int
}

// FiltroReservas parâmetros de busca de reservas.
type FiltroReservas struct {
	Status         string
	VeiculoID      *int
	DataInicio     string
	Pagina         int
	ItensPorPagina int
}

// AgendarVisitaRequest payload para agendar visita na concessionária.
type AgendarVisitaRequest struct {
	UnidadeID       int       `json:"unidadeId"`
	NomeCliente     string    `json:"nomeCliente"`
	EmailCliente    string    `json:"emailCliente"`
	TelefoneCliente string    `json:"telefoneCliente"`
	DataHoraInicio  time.Time `json:"dataHoraInicio"`
}

// AgendarTestDriveRequest payload para agendar test-drive de veículo.
type AgendarTestDriveRequest struct {
	VeiculoID       int       `json:"veiculoId"`
	UnidadeID       int       `json:"unidadeId"`
	NomeCliente     string    `json:"nomeCliente"`
	EmailCliente    string    `json:"emailCliente"`
	TelefoneCliente string    `json:"telefoneCliente"`
	CNHCliente      string    `json:"cnhCliente"`
	DataHoraInicio  time.Time `json:"dataHoraInicio"`
}

// SolicitarReservaRequest payload para reservar um veículo.
type SolicitarReservaRequest struct {
	VeiculoID       int    `json:"veiculoId"`
	NomeCliente     string `json:"nomeCliente"`
	EmailCliente    string `json:"emailCliente"`
	TelefoneCliente string `json:"telefoneCliente"`
	Observacoes     string `json:"observacoes"`
}

