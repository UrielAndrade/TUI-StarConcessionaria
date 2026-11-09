package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/UrielAndrade/TUI-StarConcessionaria/internal/model"
	"github.com/go-resty/resty/v2"
)

// HandleResponseError analisa a resposta HTTP e extrai ProblemDetails se presente.
func HandleResponseError(resp *resty.Response) error {
	if resp.IsSuccess() {
		return nil
	}

	body := resp.Body()
	if len(body) > 0 {
		var pd model.ProblemDetails
		if err := json.Unmarshal(body, &pd); err == nil && pd.Title != "" {
			return &pd
		}
	}

	switch resp.StatusCode() {
	case http.StatusUnauthorized:
		return fmt.Errorf("não autorizado (401): efetue login para continuar")
	case http.StatusForbidden:
		return fmt.Errorf("acesso negado (403): privilégios insuficientes")
	case http.StatusNotFound:
		return fmt.Errorf("recurso não encontrado (404)")
	case http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout:
		return fmt.Errorf("servidor indisponível (%d): provável cold start ou manutenção", resp.StatusCode())
	default:
		return fmt.Errorf("erro HTTP %d: %s", resp.StatusCode(), resp.Status())
	}
}
