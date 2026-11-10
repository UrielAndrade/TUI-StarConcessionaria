package veiculos

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
	EstadoNovoVeiculo
	EstadoAlterarSituacao
)

type Model struct {
	Table        table.Model
	Veiculos     []model.VeiculoCatalogoResponse
	TotalItens   int
	PaginaAtual  int
	Carregando   bool
	Erro         string
	Width        int
	Height       int
	Estado       EstadoView
	NovoForm     *forms.NovoVeiculoForm
	SituacaoForm *forms.AlterarSituacaoForm
	IsAutenticado bool
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
		Estado:      EstadoLista,
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

func (m Model) VeiculoSelecionado() *model.VeiculoCatalogoResponse {
	c := m.Table.Cursor()
	if c >= 0 && c < len(m.Veiculos) {
		return &m.Veiculos[c]
	}
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch m.Estado {
	case EstadoNovoVeiculo:
		if m.NovoForm != nil {
			var cmd tea.Cmd
			m.NovoForm, cmd = m.NovoForm.Update(msg)
			return m, cmd
		}
	case EstadoAlterarSituacao:
		if m.SituacaoForm != nil {
			var cmd tea.Cmd
			m.SituacaoForm, cmd = m.SituacaoForm.Update(msg)
			return m, cmd
		}
	}

	// Modo Lista padrão
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "a":
			if !m.IsAutenticado {
				m.Erro = "Autentique-se na aba [4] Login para cadastrar novos veículos."
				return m, nil
			}
			m.Estado = EstadoNovoVeiculo
			m.NovoForm = forms.NewNovoVeiculoForm()
			return m, nil

		case "s":
			if !m.IsAutenticado {
				m.Erro = "Autentique-se na aba [4] Login para alterar situação de veículos."
				return m, nil
			}
			sel := m.VeiculoSelecionado()
			if sel != nil {
				m.Estado = EstadoAlterarSituacao
				m.SituacaoForm = forms.NewAlterarSituacaoForm(sel.ID, fmt.Sprintf("%s %s", sel.Marca, sel.Modelo), sel.Situacao)
				return m, nil
			}
		}
	}

	var cmd tea.Cmd
	m.Table, cmd = m.Table.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if m.Estado == EstadoNovoVeiculo && m.NovoForm != nil {
		formBox := theme.ModalBoxStyle.Width(58).Render(
			lipgloss.JoinVertical(lipgloss.Left,
				lipgloss.NewStyle().Bold(true).Foreground(theme.Primary).Render("★ CADASTRAR NOVO VEÍCULO"),
				lipgloss.NewStyle().Foreground(theme.FgMuted).Render("Preencha os campos abaixo. [Esc] para cancelar."),
				"",
				m.NovoForm.View(),
			),
		)
		return lipgloss.Place(m.Width, m.Height-4, lipgloss.Center, lipgloss.Center, formBox)
	}

	if m.Estado == EstadoAlterarSituacao && m.SituacaoForm != nil {
		formBox := theme.ModalBoxStyle.Width(54).Render(
			lipgloss.JoinVertical(lipgloss.Left,
				lipgloss.NewStyle().Bold(true).Foreground(theme.Secondary).Render("⚙ ATUALIZAR STATUS OPERACIONAL"),
				lipgloss.NewStyle().Foreground(theme.FgMuted).Render("[Esc] para voltar sem salvar."),
				"",
				m.SituacaoForm.View(),
			),
		)
		return lipgloss.Place(m.Width, m.Height-4, lipgloss.Center, lipgloss.Center, formBox)
	}

	if m.Carregando {
		return lipgloss.NewStyle().
			Padding(4).
			Foreground(theme.Accent).
			Render("⏳ Carregando catálogo de veículos...")
	}

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

	var statusAviso string
	if m.Erro != "" {
		statusAviso = lipgloss.NewStyle().Foreground(theme.Warning).Bold(true).Render(" ⚠ " + m.Erro)
	}

	paginacaoInfo := lipgloss.NewStyle().
		Foreground(theme.FgMuted).
		Render(fmt.Sprintf("Exibindo %d de %d veículo(s) | [a] Novo Veículo  [s] Alterar Status  [↑/↓] Navegar%s", len(m.Veiculos), m.TotalItens, statusAviso))

	return lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.NewStyle().Bold(true).Foreground(theme.FgLight).Padding(1, 0).Render("CATÁLOGO DE VEÍCULOS NO ESTOQUE"),
		renderedTable,
		paginacaoInfo,
		"",
		previewCard,
	)
}
