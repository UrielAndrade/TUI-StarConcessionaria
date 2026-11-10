package model

// FormaPagamentoDetalheResponse representa uma condição de pagamento institucional.
type FormaPagamentoDetalheResponse struct {
	ID           int    `json:"id"`
	Nome         string `json:"nome"`
	Descricao    string `json:"descricao"`
	AceitaOnline bool   `json:"aceitaOnline"`
}

// FaqResponse representa uma pergunta frequente da concessionária.
type FaqResponse struct {
	ID       int    `json:"id"`
	Pergunta string `json:"pergunta"`
	Resposta string `json:"resposta"`
}
