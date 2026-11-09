package model

import "time"

// LoginRequest payload para autenticação.
type LoginRequest struct {
	Email string `json:"email"`
	Senha string `json:"senha"`
}

// UsuarioResponse dados do usuário autenticado.
type UsuarioResponse struct {
	ID     int    `json:"id"`
	Nome   string `json:"nome"`
	Email  string `json:"email"`
	Perfil string `json:"perfil"`
}

// LoginResponse resposta com JWT e perfil.
type LoginResponse struct {
	Token    string          `json:"token"`
	ExpiraEm time.Time       `json:"expiraEm"`
	Usuario  UsuarioResponse `json:"usuario"`
}
