package model

// UnidadeResponse concessionária ou filial.
type UnidadeResponse struct {
	ID          int      `json:"id"`
	Nome        string   `json:"nome"`
	Telefone    string   `json:"telefone"`
	CEP         string   `json:"cep"`
	Endereco    string   `json:"endereco"`
	Numero      string   `json:"numero"`
	Complemento *string  `json:"complemento"`
	Bairro      string   `json:"bairro"`
	Cidade      string   `json:"cidade"`
	Estado      string   `json:"estado"`
	Latitude    float64  `json:"latitude"`
	Longitude   float64  `json:"longitude"`
}

// SalvarUnidadeRequest payload para cadastrar filial.
type SalvarUnidadeRequest struct {
	Nome        string  `json:"nome"`
	Telefone    string  `json:"telefone"`
	CEP         string  `json:"cep"`
	Endereco    string  `json:"endereco"`
	Numero      string  `json:"numero"`
	Complemento string  `json:"complemento"`
	Bairro      string  `json:"bairro"`
	Cidade      string  `json:"cidade"`
	Estado      string  `json:"estado"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
}
