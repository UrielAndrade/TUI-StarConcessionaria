package unidades

import (
	"fmt"
	"strings"

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
	EstadoNovaUnidade
)

type Model struct {
	Table         table.Model
	Unidades      []model.UnidadeResponse
	Carregando    bool
	Erro          string
	Width         int
	Height        int
	Estado        EstadoView
	NovaForm      *forms.NovaUnidadeForm
	IsAutenticado bool
}

func New() Model {
	cols := []table.Column{
		{Title: "ID", Width: 4},
		{Title: "Nome da Concessionária", Width: 36},
		{Title: "Cidade / UF", Width: 20},
		{Title: "Bairro", Width: 18},
		{Title: "Telefone", Width: 16},
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
		Table:      t,
		Carregando: true,
		Estado:     EstadoLista,
	}
}

func (m *Model) SetDados(unidades []model.UnidadeResponse) {
	m.Unidades = unidades
	m.Carregando = false
	m.Erro = ""

	var rows []table.Row
	for _, u := range unidades {
		cidadeUF := fmt.Sprintf("%s / %s", u.Cidade, u.Estado)
		rows = append(rows, table.Row{
			fmt.Sprintf("%d", u.ID),
			u.Nome,
			cidadeUF,
			u.Bairro,
			u.Telefone,
		})
	}
	m.Table.SetRows(rows)
}

func (m Model) UnidadeSelecionada() *model.UnidadeResponse {
	c := m.Table.Cursor()
	if c >= 0 && c < len(m.Unidades) {
		return &m.Unidades[c]
	}
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if m.Estado == EstadoNovaUnidade && m.NovaForm != nil {
		var cmd tea.Cmd
		m.NovaForm, cmd = m.NovaForm.Update(msg)
		return m, cmd
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "a":
			if !m.IsAutenticado {
				m.Erro = "Autentique-se na aba [0] Perfil para cadastrar novas lojas/filiais."
				return m, nil
			}
			m.Estado = EstadoNovaUnidade
			m.NovaForm = forms.NewNovaUnidadeForm()
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.Table, cmd = m.Table.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if m.Estado == EstadoNovaUnidade && m.NovaForm != nil {
		formBox := theme.ModalBoxStyle.Width(58).Render(
			lipgloss.JoinVertical(lipgloss.Left,
				lipgloss.NewStyle().Bold(true).Foreground(theme.Primary).Render("🏢 CADASTRAR CONCESSIONÁRIA / FILIAL"),
				lipgloss.NewStyle().Foreground(theme.FgMuted).Render("Preencha os dados da loja. [Esc] para cancelar."),
				"",
				m.NovaForm.View(),
			),
		)
		return lipgloss.Place(m.Width, m.Height-4, lipgloss.Center, lipgloss.Center, formBox)
	}

	if m.Carregando {
		return lipgloss.NewStyle().
			Padding(4).
			Foreground(theme.Accent).
			Render("⏳ Carregando unidades e filiais da concessionária...")
	}

	if m.Erro != "" {
		return lipgloss.NewStyle().
			Padding(2).
			Foreground(theme.Danger).
			Render(fmt.Sprintf("Erro ao carregar concessionárias: %s", m.Erro))
	}

	renderedTable := m.Table.View()

	cursor := m.Table.Cursor()
	var previewCard string
	if cursor >= 0 && cursor < len(m.Unidades) {
		sel := m.Unidades[cursor]
		compl := "Nenhum"
		if sel.Complemento != nil && *sel.Complemento != "" {
			compl = *sel.Complemento
		}

		info := []string{
			lipgloss.NewStyle().Bold(true).Foreground(theme.Secondary).Render(fmt.Sprintf("FICHA DA FILIAL: #%d - %s", sel.ID, sel.Nome)),
			fmt.Sprintf("• Endereço: %s, Nº %s (%s) - Bairro: %s", sel.Endereco, sel.Numero, compl, sel.Bairro),
			fmt.Sprintf("• Município: %s - %s | CEP: %s | Telefone: %s", sel.Cidade, sel.Estado, sel.CEP, sel.Telefone),
			fmt.Sprintf("• Geolocalização: Latitude %.4f | Longitude %.4f", sel.Latitude, sel.Longitude),
		}
		previewCard = theme.CardStyle.Width(m.Width - 4).Render(strings.Join(info, "\n"))
	}

	paginacaoInfo := lipgloss.NewStyle().
		Foreground(theme.FgMuted).
		Render(fmt.Sprintf("Exibindo %d filial(is) ativa(s) | [a] Nova Concessionária  [↑/↓] Navegar", len(m.Unidades)))

	return lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.NewStyle().Bold(true).Foreground(theme.FgLight).Padding(1, 0).Render("REDE DE CONCESSIONÁRIAS E LOJAS FÍSICAS"),
		renderedTable,
		paginacaoInfo,
		"",
		previewCard,
	)
}
