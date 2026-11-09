package model

import "fmt"

// PagedResult representa o envelope genérico de paginação da API.
type PagedResult[T any] struct {
	Itens          []T `json:"itens"`
	Pagina         int `json:"pagina"`
	ItensPorPagina int `json:"itensPorPagina"`
	TotalItens     int `json:"totalItens"`
	TotalPaginas   int `json:"totalPaginas"`
}

// ProblemDetails representa a estrutura de erro RFC 7807 da API ASP.NET Core.
type ProblemDetails struct {
	Type     string              `json:"type,omitempty"`
	Title    string              `json:"title"`
	Status   int                 `json:"status"`
	Detail   string              `json:"detail,omitempty"`
	Instance string              `json:"instance,omitempty"`
	Errors   map[string][]string `json:"errors,omitempty"`
	TraceID  string              `json:"traceId,omitempty"`
}

func (p *ProblemDetails) Error() string {
	if len(p.Errors) > 0 {
		var firstErr string
		for campo, msgs := range p.Errors {
			if len(msgs) > 0 {
				firstErr = fmt.Sprintf("%s: %s", campo, msgs[0])
				break
			}
		}
		if firstErr != "" {
			return fmt.Sprintf("%s (%s)", p.Title, firstErr)
		}
	}
	if p.Detail != "" {
		return fmt.Sprintf("%s: %s", p.Title, p.Detail)
	}
	return p.Title
}
