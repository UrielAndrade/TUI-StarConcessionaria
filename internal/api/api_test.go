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
		unidades, err := client.ListUnidades(ctx, "", "")
		if err != nil {
			t.Fatalf("ListUnidades falhou: %v", err)
		}
		if len(unidades) == 0 {
			t.Errorf("Esperava ao menos 1 unidade cadastrada")
		}
		t.Logf("Unidades retornadas: %d", len(unidades))
	})

	t.Run("Get Unidade 1", func(t *testing.T) {
		unidade, err := client.GetUnidade(ctx, 1)
		if err != nil {
			t.Fatalf("GetUnidade falhou: %v", err)
		}
		if unidade.ID != 1 {
			t.Errorf("Esperava Unidade com ID 1, obteve %d", unidade.ID)
		}
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

	t.Run("Get Veiculo 1 Detalhes", func(t *testing.T) {
		veiculo, err := client.GetVeiculo(ctx, 1)
		if err != nil {
			t.Fatalf("GetVeiculo falhou: %v", err)
		}
		if veiculo.ID != 1 {
			t.Errorf("Esperava Veículo ID 1, obteve %d", veiculo.ID)
		}
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

	t.Run("Get Peca 1", func(t *testing.T) {
		peca, err := client.GetPeca(ctx, 1)
		if err != nil {
			t.Fatalf("GetPeca falhou: %v", err)
		}
		if peca.ID != 1 {
			t.Errorf("Esperava Peça ID 1, obteve %d", peca.ID)
		}
	})

	t.Run("List Formas Pagamento", func(t *testing.T) {
		formas, err := client.ListFormasPagamento(ctx)
		if err != nil {
			t.Fatalf("ListFormasPagamento falhou: %v", err)
		}
		if len(formas) == 0 {
			t.Errorf("Esperava formas de pagamento cadastradas")
		}
		t.Logf("Formas de pagamento retornadas: %d", len(formas))
	})

	t.Run("List Faqs", func(t *testing.T) {
		faqs, err := client.ListFaqs(ctx, "")
		if err != nil {
			t.Fatalf("ListFaqs falhou: %v", err)
		}
		if len(faqs) == 0 {
			t.Errorf("Esperava FAQs cadastradas")
		}
		t.Logf("FAQs retornadas: %d", len(faqs))
	})
}
