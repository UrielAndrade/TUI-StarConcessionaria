package veiculos

import (
	"fmt"
	"strings"

	"github.com/UrielAndrade/TUI-StarConcessionaria/internal/model"
	"github.com/UrielAndrade/TUI-StarConcessionaria/internal/ui/theme"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	Table       table.Model
	Veiculos    []model.VeiculoCatalogoResponse
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
		{Title: "Tipo", Width: 8},
		{Title: "Marca", Width: 12},
		{Title: "Modelo", Width: 26},
		{Title: "Ano", Width: 6},
		{Title: "Preço (R$)", Width: 14},
		{Title: "Km", Width: 10},
		{Title: "Cor", Width: 14},
		{Title: "Situação", Width: 12},
	}

	t := table.New(
		table.WithColumns(cols),
		table.WithFocused(true),
		table.WithHeight(8),
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

func (m *Model) SetDados(resultado *model.PagedResult[model.VeiculoCatalogoResponse]) {
	m.Veiculos = resultado.Itens
	m.TotalItens = resultado.TotalItens
	m.PaginaAtual = resultado.Pagina
	m.Carregando = false
	m.Erro = ""

	var rows []table.Row
	for _, v := range resultado.Itens {
		rows = append(rows, table.Row{
			fmt.Sprintf("%d", v.ID),
			v.Tipo,
			v.Marca,
			v.Modelo,
			fmt.Sprintf("%d", v.Ano),
			fmt.Sprintf("%.2f", v.Preco),
			fmt.Sprintf("%d km", v.Quilometragem),
			v.Cor,
			v.Situacao,
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
			Render("⏳ Carregando catálogo de veículos...")
	}

	if m.Erro != "" {
		return lipgloss.NewStyle().
			Padding(2).
			Foreground(theme.Danger).
			Render(fmt.Sprintf("Erro ao carregar veículos: %s", m.Erro))
	}

	// Tabela de veículos
	renderedTable := m.Table.View()

	// Card do Veículo Selecionado
	cursor := m.Table.Cursor()
	var previewCard string
	if cursor >= 0 && cursor < len(m.Veiculos) {
		sel := m.Veiculos[cursor]
		info := []string{
			lipgloss.NewStyle().Bold(true).Foreground(theme.Secondary).Render(fmt.Sprintf("DETALHES: #%d - %s %s (%d)", sel.ID, sel.Marca, sel.Modelo, sel.Ano)),
			fmt.Sprintf("• Loja: %s", sel.LojaNome),
			fmt.Sprintf("• Valor de Venda: R$ %.2f | Categoria: %s | Cor: %s", sel.Preco, sel.Tipo, sel.Cor),
			fmt.Sprintf("• Quilometragem: %d km | Situação: %s", sel.Quilometragem, sel.Situacao),
			fmt.Sprintf("• Foto de Capa: %s", sel.FotoCapaURL),
		}
		previewCard = theme.CardStyle.Width(m.Width - 4).Render(strings.Join(info, "\n"))
	}

	paginacaoInfo := lipgloss.NewStyle().
		Foreground(theme.FgMuted).
		Render(fmt.Sprintf("Exibindo %d de %d veículo(s) cadastrado(s) | Pressione [↑/↓] para navegar na tabela", len(m.Veiculos), m.TotalItens))

	return lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.NewStyle().Bold(true).Foreground(theme.FgLight).Padding(1, 0).Render("CATÁLOGO DE VEÍCULOS NO ESTOQUE"),
		renderedTable,
		paginacaoInfo,
		"",
		previewCard,
	)
}
