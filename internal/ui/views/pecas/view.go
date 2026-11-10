package pecas

import (
	"fmt"

	"github.com/UrielAndrade/TUI-StarConcessionaria/internal/model"
	"github.com/UrielAndrade/TUI-StarConcessionaria/internal/ui/forms"
	"github.com/UrielAndrade/TUI-StarConcessionaria/internal/ui/theme"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type EstadoView int

const (
	EstadoLista EstadoView = iota
	EstadoNovaPeca
	EstadoVenderPeca
)

type Model struct {
	Table         table.Model
	Pecas         []model.PecaResponse
	TotalItens    int
	PaginaAtual   int
	Carregando    bool
	Erro          string
	Width         int
	Height        int
	Estado        EstadoView
	NovaForm      *forms.NovaPecaForm
	VenderForm    *forms.VenderPecaForm
	IsAutenticado bool
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
		Estado:      EstadoLista,
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

func (m Model) PecaSelecionada() *model.PecaResponse {
	c := m.Table.Cursor()
	if c >= 0 && c < len(m.Pecas) {
		return &m.Pecas[c]
	}
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch m.Estado {
	case EstadoNovaPeca:
		if m.NovaForm != nil {
			var cmd tea.Cmd
			m.NovaForm, cmd = m.NovaForm.Update(msg)
			return m, cmd
		}
	case EstadoVenderPeca:
		if m.VenderForm != nil {
			var cmd tea.Cmd
			m.VenderForm, cmd = m.VenderForm.Update(msg)
			return m, cmd
		}
	}

	// Modo Lista padrão
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "a":
			if !m.IsAutenticado {
				m.Erro = "Autentique-se na aba [4] Login para cadastrar novas peças."
				return m, nil
			}
			m.Estado = EstadoNovaPeca
			m.NovaForm = forms.NewNovaPecaForm()
			return m, nil

		case "v":
			if !m.IsAutenticado {
				m.Erro = "Autentique-se na aba [4] Login para dar baixa em peças."
				return m, nil
			}
			sel := m.PecaSelecionada()
			if sel != nil {
				if sel.QuantidadeEstoque <= 0 {
					m.Erro = fmt.Sprintf("A peça '%s' não possui estoque disponível para venda.", sel.Nome)
					return m, nil
				}
				m.Estado = EstadoVenderPeca
				m.VenderForm = forms.NewVenderPecaForm(sel.ID, sel.Nome, sel.QuantidadeEstoque)
				return m, nil
			}
		}
	}

	var cmd tea.Cmd
	m.Table, cmd = m.Table.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if m.Estado == EstadoNovaPeca && m.NovaForm != nil {
		formBox := theme.ModalBoxStyle.Width(58).Render(
			lipgloss.JoinVertical(lipgloss.Left,
				lipgloss.NewStyle().Bold(true).Foreground(theme.Primary).Render("★ CADASTRAR NOVA PEÇA NO ESTOQUE"),
				lipgloss.NewStyle().Foreground(theme.FgMuted).Render("[Esc] para cancelar."),
				"",
				m.NovaForm.View(),
			),
		)
		return lipgloss.Place(m.Width, m.Height-4, lipgloss.Center, lipgloss.Center, formBox)
	}

	if m.Estado == EstadoVenderPeca && m.VenderForm != nil {
		formBox := theme.ModalBoxStyle.Width(58).Render(
			lipgloss.JoinVertical(lipgloss.Left,
				lipgloss.NewStyle().Bold(true).Foreground(theme.Secondary).Render("💰 REGISTRAR VENDA / BAIXA DE PEÇA"),
				lipgloss.NewStyle().Foreground(theme.FgMuted).Render("[Esc] para voltar sem salvar."),
				"",
				m.VenderForm.View(),
			),
		)
		return lipgloss.Place(m.Width, m.Height-4, lipgloss.Center, lipgloss.Center, formBox)
	}

	if m.Carregando {
		return lipgloss.NewStyle().
			Padding(4).
			Foreground(theme.Accent).
			Render("⏳ Carregando catálogo de peças automotivas...")
	}

	var statusAviso string
	if m.Erro != "" {
		statusAviso = lipgloss.NewStyle().Foreground(theme.Warning).Bold(true).Render(" ⚠ " + m.Erro)
	}

	paginacaoInfo := lipgloss.NewStyle().
		Foreground(theme.FgMuted).
		Render(fmt.Sprintf("Exibindo %d de %d peça(s) no estoque | [a] Nova Peça  [v] Registrar Venda  [↑/↓] Navegar%s", len(m.Pecas), m.TotalItens, statusAviso))

	return lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.NewStyle().Bold(true).Foreground(theme.FgLight).Padding(1, 0).Render("ESTOQUE E PEÇAS AUTOMOTIVAS"),
		m.Table.View(),
		paginacaoInfo,
	)
}
