package model

// PecaResponse representação de uma peça automotiva em estoque.
type PecaResponse struct {
	ID                int     `json:"id"`
	Nome              string  `json:"nome"`
	Descricao         string  `json:"descricao"`
	Preco             float64 `json:"preco"`
	QuantidadeEstoque int     `json:"quantidadeEstoque"`
	Disponivel        bool    `json:"disponivel"`
}

// SalvarPecaRequest payload de criação ou edição de peça.
type SalvarPecaRequest struct {
	Nome              string  `json:"nome"`
	Descricao         string  `json:"descricao"`
	Preco             float64 `json:"preco"`
	QuantidadeEstoque int     `json:"quantidadeEstoque"`
}

// RegistrarVendaPecaRequest payload para dar baixa no estoque por venda direta.
type RegistrarVendaPecaRequest struct {
	Quantidade      int    `json:"quantidade"`
	NomeCliente     string `json:"nomeCliente"`
	EmailCliente    string `json:"emailCliente"`
	TelefoneCliente string `json:"telefoneCliente"`
}

// FiltroPecas parâmetros de busca para peças.
type FiltroPecas struct {
	Busca          string
	Disponivel     *bool
	Pagina         int
	ItensPorPagina int
}
