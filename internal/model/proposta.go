package model

import "time"

// PropostaResponse proposta de aquisição/financiamento.
type PropostaResponse struct {
	ID              int       `json:"id"`
	Protocolo       string    `json:"protocolo"`
	VeiculoID       int       `json:"veiculoId"`
	UnidadeID       int       `json:"unidadeId"`
	NomeCliente     string    `json:"nomeCliente"`
	EmailCliente    string    `json:"emailCliente"`
	TelefoneCliente string    `json:"telefoneCliente"`
	FormaPagamento  string    `json:"formaPagamento"`
	PossuiTroca     bool      `json:"possuiTroca"`
	ValorEntrada    float64   `json:"valorEntrada"`
	DataSolicitacao time.Time `json:"dataSolicitacao"`
	Mensagem        string    `json:"mensagem"`
}

// SolicitarPropostaRequest payload para envio de proposta.
type SolicitarPropostaRequest struct {
	VeiculoID       int     `json:"veiculoId"`
	UnidadeID       int     `json:"unidadeId"`
	NomeCliente     string  `json:"nomeCliente"`
	EmailCliente    string  `json:"emailCliente"`
	TelefoneCliente string  `json:"telefoneCliente"`
	FormaPagamento  string  `json:"formaPagamento"`
	PossuiTroca     bool    `json:"possuiTroca"`
	DadosTroca      string  `json:"dadosTroca"`
	ValorEntrada    float64 `json:"valorEntrada"`
}

// FiltroPropostas parâmetros de busca para propostas.
type FiltroPropostas struct {
	Respondido     *bool
	UnidadeID      *int
	Pagina         int
	ItensPorPagina int
}
