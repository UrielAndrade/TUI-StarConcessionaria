package pecas

import (
	"fmt"

	"github.com/UrielAndrade/TUI-StarConcessionaria/internal/model"
	"github.com/UrielAndrade/TUI-StarConcessionaria/internal/ui/theme"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	Table       table.Model
	Pecas       []model.PecaResponse
	TotalItens  int
	PaginaAtual int
	Carregando  bool
	Erro        string
	Width       int
	Height      int
}

func New() Model {
	cols := []table.Column{
		{Title: "ID", Width: 4},
		{Title: "Nome da Peça", Width: 32},
		{Title: "Descrição", Width: 40},
		{Title: "Preço (R$)", Width: 12},
		{Title: "Estoque", Width: 10},
		{Title: "Status", Width: 12},
	}

	t := table.New(
		table.WithColumns(cols),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(theme.BorderCol).
		BorderBottom(true).
		Bold(true).
		Foreground(theme.FgLight).
		Background(theme.Primary)

	s.Selected = s.Selected.
		Foreground(theme.FgLight).
		Background(lipgloss.Color("#2D3748")).
		Bold(true)

	t.SetStyles(s)

	return Model{
		Table:       t,
		Carregando:  true,
		PaginaAtual: 1,
	}
}

func (m *Model) SetDados(resultado *model.PagedResult[model.PecaResponse]) {
	m.Pecas = resultado.Itens
	m.TotalItens = resultado.TotalItens
	m.PaginaAtual = resultado.Pagina
	m.Carregando = false
	m.Erro = ""

	var rows []table.Row
	for _, p := range resultado.Itens {
		status := "Indisponível"
		if p.Disponivel && p.QuantidadeEstoque > 0 {
			status = "Disponível"
		}

		rows = append(rows, table.Row{
			fmt.Sprintf("%d", p.ID),
			p.Nome,
			p.Descricao,
			fmt.Sprintf("%.2f", p.Preco),
			fmt.Sprintf("%d un.", p.QuantidadeEstoque),
			status,
		})
	}
	m.Table.SetRows(rows)
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	m.Table, cmd = m.Table.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if m.Carregando {
		return lipgloss.NewStyle().
			Padding(4).
			Foreground(theme.Accent).
			Render("⏳ Carregando catálogo de peças automotivas...")
	}

	if m.Erro != "" {
		return lipgloss.NewStyle().
			Padding(2).
			Foreground(theme.Danger).
			Render(fmt.Sprintf("Erro ao carregar peças: %s", m.Erro))
	}

	paginacaoInfo := lipgloss.NewStyle().
		Foreground(theme.FgMuted).
		Render(fmt.Sprintf("Exibindo %d de %d peça(s) no estoque | Pressione [↑/↓] para navegar", len(m.Pecas), m.TotalItens))

	return lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.NewStyle().Bold(true).Foreground(theme.FgLight).Padding(1, 0).Render("ESTOQUE E PEÇAS AUTOMOTIVAS"),
		m.Table.View(),
		paginacaoInfo,
	)
}
