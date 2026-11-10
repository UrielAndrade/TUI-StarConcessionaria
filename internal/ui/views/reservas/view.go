package reservas

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
	EstadoNovaReserva
	EstadoStatus
)

type Model struct {
	Table         table.Model
	Reservas      []model.ReservaResponse
	TotalItens    int
	PaginaAtual   int
	Carregando    bool
	Erro          string
	Width         int
	Height        int
	Estado        EstadoView
	NovaForm      *forms.NovaReservaForm
	StatusForm    *forms.StatusReservaForm
	IsAutenticado bool
}

func New() Model {
	cols := []table.Column{
		{Title: "ID", Width: 4},
		{Title: "Código", Width: 12},
		{Title: "Veículo", Width: 30},
		{Title: "Cliente", Width: 20},
		{Title: "Telefone", Width: 16},
		{Title: "Data Solicitação", Width: 18},
		{Title: "Situação", Width: 14},
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

func (m *Model) SetDados(resultado *model.PagedResult[model.ReservaResponse]) {
	m.Reservas = resultado.Itens
	m.TotalItens = resultado.TotalItens
	m.PaginaAtual = resultado.Pagina
	m.Carregando = false
	m.Erro = ""

	var rows []table.Row
	for _, r := range resultado.Itens {
		dataStr := r.DataSolicitacao.Format("02/01/2006 15:04")
		rows = append(rows, table.Row{
			fmt.Sprintf("%d", r.ID),
			r.Codigo,
			r.VeiculoDescricao,
			r.NomeCliente,
			r.TelefoneCliente,
			dataStr,
			r.Situacao,
		})
	}
	m.Table.SetRows(rows)
}

func (m Model) ReservaSelecionada() *model.ReservaResponse {
	c := m.Table.Cursor()
	if c >= 0 && c < len(m.Reservas) {
		return &m.Reservas[c]
	}
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch m.Estado {
	case EstadoNovaReserva:
		if m.NovaForm != nil {
			var cmd tea.Cmd
			m.NovaForm, cmd = m.NovaForm.Update(msg)
			return m, cmd
		}
	case EstadoStatus:
		if m.StatusForm != nil {
			var cmd tea.Cmd
			m.StatusForm, cmd = m.StatusForm.Update(msg)
			return m, cmd
		}
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "a":
			m.Estado = EstadoNovaReserva
			m.NovaForm = forms.NewNovaReservaForm(1)
			return m, nil

		case "s":
			if !m.IsAutenticado {
				m.Erro = "Autentique-se na aba [9] Perfil para alterar status de reservas."
				return m, nil
			}
			sel := m.ReservaSelecionada()
			if sel != nil {
				m.Estado = EstadoStatus
				m.StatusForm = forms.NewStatusReservaForm(sel.ID, sel.Situacao)
				return m, nil
			}
		}
	}

	var cmd tea.Cmd
	m.Table, cmd = m.Table.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if m.Estado == EstadoNovaReserva && m.NovaForm != nil {
		formBox := theme.ModalBoxStyle.Width(58).Render(
			lipgloss.JoinVertical(lipgloss.Left,
				lipgloss.NewStyle().Bold(true).Foreground(theme.Accent).Render("🔒 SOLICITAR RESERVA DE VEÍCULO"),
				lipgloss.NewStyle().Foreground(theme.FgMuted).Render("[Esc] para cancelar."),
				"",
				m.NovaForm.View(),
			),
		)
		return lipgloss.Place(m.Width, m.Height-4, lipgloss.Center, lipgloss.Center, formBox)
	}

	if m.Estado == EstadoStatus && m.StatusForm != nil {
		formBox := theme.ModalBoxStyle.Width(56).Render(
			lipgloss.JoinVertical(lipgloss.Left,
				lipgloss.NewStyle().Bold(true).Foreground(theme.Primary).Render("⚙ ALTERAR SITUAÇÃO DA RESERVA"),
				lipgloss.NewStyle().Foreground(theme.FgMuted).Render("[Esc] para voltar sem salvar."),
				"",
				m.StatusForm.View(),
			),
		)
		return lipgloss.Place(m.Width, m.Height-4, lipgloss.Center, lipgloss.Center, formBox)
	}

	if m.Carregando {
		return lipgloss.NewStyle().
			Padding(4).
			Foreground(theme.Accent).
			Render("⏳ Carregando lista de reservas de veículos...")
	}

	if m.Erro != "" {
		return lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(theme.Warning).
			Padding(1, 2).
			Render(fmt.Sprintf("Aviso: %s\n\n(Dica: Reservas exigem autenticação na aba [9] Perfil)", m.Erro))
	}

	renderedTable := m.Table.View()

	cursor := m.Table.Cursor()
	var previewCard string
	if cursor >= 0 && cursor < len(m.Reservas) {
		sel := m.Reservas[cursor]
		info := []string{
			lipgloss.NewStyle().Bold(true).Foreground(theme.Accent).Render(fmt.Sprintf("RESERVA ATIVA: %s - %s", sel.Codigo, sel.VeiculoDescricao)),
			fmt.Sprintf("• Cliente Solicitante: %s | E-mail: %s | Telefone: %s", sel.NomeCliente, sel.EmailCliente, sel.TelefoneCliente),
			fmt.Sprintf("• Data de Abertura: %s | Situação Atual: %s", sel.DataSolicitacao.Format("02/01/2006 às 15:04"), sel.Situacao),
		}
		previewCard = theme.CardStyle.Width(m.Width - 4).Render(strings.Join(info, "\n"))
	}

	paginacaoInfo := lipgloss.NewStyle().
		Foreground(theme.FgMuted).
		Render(fmt.Sprintf("Exibindo %d de %d reserva(s) | [a] Solicitar Reserva  [s] Alterar Status  [↑/↓] Navegar", len(m.Reservas), m.TotalItens))

	return lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.NewStyle().Bold(true).Foreground(theme.FgLight).Padding(1, 0).Render("RESERVAS E BLOQUEIOS DE ESTOQUE"),
		renderedTable,
		paginacaoInfo,
		"",
		previewCard,
	)
}
