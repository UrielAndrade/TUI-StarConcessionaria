package model

import "time"

// VendaResponse registro de venda de veículo finalizada.
type VendaResponse struct {
	ID                 int       `json:"id"`
	VeiculoID          int       `json:"veiculoId"`
	VeiculoDescricao   string    `json:"veiculoDescricao"`
	ClienteID          int       `json:"clienteId"`
	NomeCliente        string    `json:"nomeCliente"`
	VendedorID         int       `json:"vendedorId"`
	NomeVendedor       string    `json:"nomeVendedor"`
	LojaID             int       `json:"lojaId"`
	NomeLoja           string    `json:"nomeLoja"`
	FormaPagamentoID   int       `json:"formaPagamentoId"`
	FormaPagamentoNome string    `json:"formaPagamentoNome"`
	ValorFinal         float64   `json:"valorFinal"`
	DataVenda          time.Time `json:"dataVenda"`
}

// RegistrarVendaRequest payload para registrar venda de veículo.
type RegistrarVendaRequest struct {
	VeiculoID        int     `json:"veiculoId"`
	NomeCliente      string  `json:"nomeCliente"`
	EmailCliente     string  `json:"emailCliente"`
	TelefoneCliente  string  `json:"telefoneCliente"`
	FormaPagamentoID int     `json:"formaPagamentoId"`
	ValorFinal       float64 `json:"valorFinal"`
}

// FiltroVendas parâmetros de busca de vendas.
type FiltroVendas struct {
	UnidadeID      *int
	VendedorID     *int
	PeriodoInicio  string
	PeriodoFim     string
	Pagina         int
	ItensPorPagina int
}
