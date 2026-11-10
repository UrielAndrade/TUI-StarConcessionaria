package model

import "time"

// ContatoResponse mensagem de SAC ou lead.
type ContatoResponse struct {
	ID                    int       `json:"id"`
	Protocolo             string    `json:"protocolo"`
	Tipo                  string    `json:"tipo"`
	ClienteID             int       `json:"clienteId"`
	NomeCliente           string    `json:"nomeCliente"`
	EmailCliente          string    `json:"emailCliente"`
	TelefoneCliente       string    `json:"telefoneCliente"`
	VeiculoID             *int      `json:"veiculoId"`
	VeiculoDescricao      *string   `json:"veiculoDescricao"`
	UnidadeID             int       `json:"unidadeId"`
	UnidadeNome           string    `json:"unidadeNome"`
	UsuarioAtendimentoID  *int      `json:"usuarioAtendimentoId"`
	UsuarioAtendimentoNome *string  `json:"usuarioAtendimentoNome"`
	Respondido            bool      `json:"respondido"`
	Mensagem              string    `json:"mensagem"`
	Data                  time.Time `json:"data"`
}

// NovoContatoRequest payload para envio de mensagem geral ao SAC.
type NovoContatoRequest struct {
	Assunto   string `json:"assunto"`
	Nome      string `json:"nome"`
	Email     string `json:"email"`
	Telefone  string `json:"telefone"`
	Mensagem  string `json:"mensagem"`
	VeiculoID *int   `json:"veiculoId,omitempty"`
	UnidadeID *int   `json:"unidadeId,omitempty"`
}

// ContatoVendedorRequest mensagem direta para o vendedor da loja.
type ContatoVendedorRequest struct {
	VeiculoID int    `json:"veiculoId"`
	UnidadeID int    `json:"unidadeId"`
	Nome      string `json:"nome"`
	Email     string `json:"email"`
	Telefone  string `json:"telefone"`
	Mensagem  string `json:"mensagem"`
}

// ResponderContatoRequest payload para registrar resposta/atendimento.
type ResponderContatoRequest struct {
	NotasResposta string `json:"notasResposta"`
}

// FiltroContatos parâmetros de pesquisa de chamados de atendimento.
type FiltroContatos struct {
	Tipo           string
	Respondido     *bool
	UnidadeID      *int
	Pagina         int
	ItensPorPagina int
}
