package api_test

import (
	"context"
	"testing"
	"time"

	"github.com/UrielAndrade/TUI-StarConcessionaria/internal/api"
	"github.com/UrielAndrade/TUI-StarConcessionaria/internal/config"
	"github.com/UrielAndrade/TUI-StarConcessionaria/internal/model"
)

func TestLiveAPIConnectivity(t *testing.T) {
	cfg := &config.Config{
		BaseURL:          config.DefaultBaseURL,
		Timeout:          15 * time.Second,
		ColdStartTimeout: 45 * time.Second,
	}

	client := api.NewClient(cfg)
	ctx := context.Background()

	t.Run("Ping Root", func(t *testing.T) {
		err := client.Ping(ctx)
		if err != nil {
			t.Fatalf("Ping falhou: %v", err)
		}
	})

	t.Run("List Unidades", func(t *testing.T) {
		unidades, err := client.ListUnidades(ctx)
		if err != nil {
			t.Fatalf("ListUnidades falhou: %v", err)
		}
		if len(unidades) == 0 {
			t.Errorf("Esperava ao menos 1 unidade cadastrada")
		}
		t.Logf("Unidades retornadas: %d", len(unidades))
	})

	t.Run("List Veiculos Public", func(t *testing.T) {
		veiculos, err := client.ListVeiculos(ctx, model.FiltroVeiculos{
			Pagina:         1,
			ItensPorPagina: 5,
		})
		if err != nil {
			t.Fatalf("ListVeiculos falhou: %v", err)
		}
		if veiculos.TotalItens == 0 {
			t.Errorf("Esperava ao menos 1 veículo no catálogo")
		}
		t.Logf("Veículos retornados: %d de um total de %d", len(veiculos.Itens), veiculos.TotalItens)
	})

	t.Run("List Pecas Public", func(t *testing.T) {
		pecas, err := client.ListPecas(ctx, model.FiltroPecas{
			Pagina:         1,
			ItensPorPagina: 5,
		})
		if err != nil {
			t.Fatalf("ListPecas falhou: %v", err)
		}
		if pecas.TotalItens == 0 {
			t.Errorf("Esperava peças cadastradas")
		}
		t.Logf("Peças retornadas: %d de um total de %d", len(pecas.Itens), pecas.TotalItens)
	})
}
