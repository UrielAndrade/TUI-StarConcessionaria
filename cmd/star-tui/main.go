package main

import (
	"fmt"
	"os"

	"github.com/UrielAndrade/TUI-StarConcessionaria/internal/api"
	"github.com/UrielAndrade/TUI-StarConcessionaria/internal/config"
	"github.com/UrielAndrade/TUI-StarConcessionaria/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Inicialização específica de plataforma (Windows UTF-8 / Virtual Terminal)
	initPlatform()

	// Carrega configurações
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao carregar configurações: %v\n", err)
		os.Exit(1)
	}

	// Instancia cliente da API
	apiClient := api.NewClient(cfg)

	// Instancia o modelo Bubble Tea
	appModel := ui.NewAppModel(apiClient, cfg)

	// Inicia o programa TUI em tela cheia (AltScreen) com captura de mouse
	p := tea.NewProgram(
		appModel,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Erro na execução da TUI: %v\n", err)
		os.Exit(1)
	}
}
