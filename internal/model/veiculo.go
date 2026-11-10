package model

import "time"

// VeiculoFotoResponse representa uma imagem associada ao veículo.
type VeiculoFotoResponse struct {
	ID        int    `json:"id"`
	URL       string `json:"url"`
	Principal bool   `json:"principal"`
	Ordem     int    `json:"ordem"`
}

// VeiculoCatalogoResponse item resumido retornado na listagem de veículos.
type VeiculoCatalogoResponse struct {
	ID            int     `json:"id"`
	LojaID        int     `json:"lojaId"`
	LojaNome      string  `json:"lojaNome"`
	Tipo          string  `json:"tipo"`
	Marca         string  `json:"marca"`
	Modelo        string  `json:"modelo"`
	Ano           int     `json:"ano"`
	Preco         float64 `json:"preco"`
	Quilometragem int     `json:"quilometragem"`
	Cor           string  `json:"cor"`
	Situacao      string  `json:"situacao"`
	FotoCapaURL   string  `json:"fotoCapaUrl"`
}

// VeiculoDetalheResponse dados completos do veículo para visualização detalhada.
type VeiculoDetalheResponse struct {
	ID            int                   `json:"id"`
	LojaID        int                   `json:"lojaId"`
	LojaNome      string                `json:"lojaNome"`
	Tipo          string                `json:"tipo"`
	Marca         string                `json:"marca"`
	Modelo        string                `json:"modelo"`
	Ano           int     				`json:"ano"`
	Preco         float64 				`json:"preco"`
	Quilometragem int     				`json:"quilometragem"`
	Cor           string  				`json:"cor"`
	Descricao     string  				`json:"descricao"`
	Situacao      string  				`json:"situacao"`
	Fotos         []VeiculoFotoResponse `json:"fotos"`
}

// SalvarVeiculoRequest payload para criação/atualização de veículo.
type SalvarVeiculoRequest struct {
	LojaID        int     `json:"lojaId"`
	Tipo          string  `json:"tipo"`
	Marca         string  `json:"marca"`
	Modelo        string  `json:"modelo"`
	Ano           int     `json:"ano"`
	Preco         float64 `json:"preco"`
	Quilometragem int     `json:"quilometragem"`
	Cor           string  `json:"cor"`
	Descricao     string  `json:"descricao"`
}

// AlterarSituacaoRequest payload para mudar status (Disponível, Reservado, Vendido).
type AlterarSituacaoRequest struct {
	Situacao      string `json:"situacao"`
	Justificativa string `json:"justificativa"`
}

// VeiculoHistoricoResponse histórico de alterações do veículo.
type VeiculoHistoricoResponse struct {
	ID            int       `json:"id"`
	VeiculoID     int       `json:"veiculoId"`
	SituacaoAnterior string `json:"situacaoAnterior"`
	NovaSituacao  string    `json:"novaSituacao"`
	Justificativa string    `json:"justificativa"`
	DataAlteracao time.Time `json:"dataAlteracao"`
	UsuarioNome   string    `json:"usuarioNome"`
}

// FiltroVeiculos parâmetros de consulta para catálogo.
type FiltroVeiculos struct {
	Marca             string
	Modelo            string
	Tipo              string
	PrecoMin          *float64
	PrecoMax          *float64
	AnoMin            *int
	AnoMax            *int
	UnidadeID         *int
	SomenteDisponivel *bool
	Pagina            int
	ItensPorPagina    int
}

// AdicionarFotoRequest payload para adicionar foto ao veículo.
type AdicionarFotoRequest struct {
	URL       string `json:"url"`
	Principal bool   `json:"principal"`
	Ordem     int    `json:"ordem"`
}

