package agendamentos

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
	EstadoAgendarVisita
	EstadoAgendarTestDrive
	EstadoStatus
)

type Model struct {
	Table         table.Model
	Agendamentos  []model.AgendamentoResponse
	TotalItens    int
	PaginaAtual   int
	Carregando    bool
	Erro          string
	Width         int
	Height        int
	Estado        EstadoView
	VisitaForm    *forms.AgendarVisitaForm
	TestDriveForm *forms.AgendarTestDriveForm
	StatusForm    *forms.StatusAgendamentoForm
	IsAutenticado bool
}

func New() Model {
	cols := []table.Column{
		{Title: "ID", Width: 4},
		{Title: "Código", Width: 10},
		{Title: "Tipo", Width: 12},
		{Title: "Data / Horário", Width: 18},
		{Title: "Unidade", Width: 22},
		{Title: "Veículo", Width: 26},
		{Title: "Cliente", Width: 18},
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

func (m *Model) SetDados(resultado *model.PagedResult[model.AgendamentoResponse]) {
	m.Agendamentos = resultado.Itens
	m.TotalItens = resultado.TotalItens
	m.PaginaAtual = resultado.Pagina
	m.Carregando = false
	m.Erro = ""

	var rows []table.Row
	for _, a := range resultado.Itens {
		dataStr := a.DataHoraInicio.Format("02/01 15:04")
		veiculoStr := "Visita Geral"
		if a.VeiculoDescricao != nil && *a.VeiculoDescricao != "" {
			veiculoStr = *a.VeiculoDescricao
		}

		rows = append(rows, table.Row{
			fmt.Sprintf("%d", a.ID),
			a.Codigo,
			a.Tipo,
			dataStr,
			a.UnidadeNome,
			veiculoStr,
			a.NomeCliente,
			a.Situacao,
		})
	}
	m.Table.SetRows(rows)
}

func (m Model) AgendamentoSelecionado() *model.AgendamentoResponse {
	c := m.Table.Cursor()
	if c >= 0 && c < len(m.Agendamentos) {
		return &m.Agendamentos[c]
	}
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch m.Estado {
	case EstadoAgendarVisita:
		if m.VisitaForm != nil {
			var cmd tea.Cmd
			m.VisitaForm, cmd = m.VisitaForm.Update(msg)
			return m, cmd
		}
	case EstadoAgendarTestDrive:
		if m.TestDriveForm != nil {
			var cmd tea.Cmd
			m.TestDriveForm, cmd = m.TestDriveForm.Update(msg)
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
		case "v":
			m.Estado = EstadoAgendarVisita
			m.VisitaForm = forms.NewAgendarVisitaForm()
			return m, nil

		case "t":
			m.Estado = EstadoAgendarTestDrive
			m.TestDriveForm = forms.NewAgendarTestDriveForm(1)
			return m, nil

		case "s":
			if !m.IsAutenticado {
				m.Erro = "Autentique-se na aba [9] Perfil para alterar status de agendamentos."
				return m, nil
			}
			sel := m.AgendamentoSelecionado()
			if sel != nil {
				m.Estado = EstadoStatus
				m.StatusForm = forms.NewStatusAgendamentoForm(sel.ID, sel.Situacao)
				return m, nil
			}
		}
	}

	var cmd tea.Cmd
	m.Table, cmd = m.Table.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if m.Estado == EstadoAgendarVisita && m.VisitaForm != nil {
		formBox := theme.ModalBoxStyle.Width(58).Render(
			lipgloss.JoinVertical(lipgloss.Left,
				lipgloss.NewStyle().Bold(true).Foreground(theme.Primary).Render("📅 AGENDAR VISITA À CONCESSIONÁRIA"),
				lipgloss.NewStyle().Foreground(theme.FgMuted).Render("[Esc] para cancelar."),
				"",
				m.VisitaForm.View(),
			),
		)
		return lipgloss.Place(m.Width, m.Height-4, lipgloss.Center, lipgloss.Center, formBox)
	}

	if m.Estado == EstadoAgendarTestDrive && m.TestDriveForm != nil {
		formBox := theme.ModalBoxStyle.Width(58).Render(
			lipgloss.JoinVertical(lipgloss.Left,
				lipgloss.NewStyle().Bold(true).Foreground(theme.Secondary).Render("🚗 AGENDAR TEST-DRIVE COM CNH"),
				lipgloss.NewStyle().Foreground(theme.FgMuted).Render("[Esc] para cancelar."),
				"",
				m.TestDriveForm.View(),
			),
		)
		return lipgloss.Place(m.Width, m.Height-4, lipgloss.Center, lipgloss.Center, formBox)
	}

	if m.Estado == EstadoStatus && m.StatusForm != nil {
		formBox := theme.ModalBoxStyle.Width(56).Render(
			lipgloss.JoinVertical(lipgloss.Left,
				lipgloss.NewStyle().Bold(true).Foreground(theme.Accent).Render("⚙ ATUALIZAR STATUS DE AGENDAMENTO"),
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
			Render("⏳ Carregando agenda de visitas e test-drives...")
	}

	if m.Erro != "" {
		return lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(theme.Warning).
			Padding(1, 2).
			Render(fmt.Sprintf("Aviso: %s\n\n(Dica: Agendamentos exigem autenticação na aba [9] Perfil)", m.Erro))
	}

	renderedTable := m.Table.View()

	cursor := m.Table.Cursor()
	var previewCard string
	if cursor >= 0 && cursor < len(m.Agendamentos) {
		sel := m.Agendamentos[cursor]
		veicDesc := "Visita institucional à loja"
		if sel.VeiculoDescricao != nil && *sel.VeiculoDescricao != "" {
			veicDesc = *sel.VeiculoDescricao
		}

		info := []string{
			lipgloss.NewStyle().Bold(true).Foreground(theme.Secondary).Render(fmt.Sprintf("DETALHES DO AGENDAMENTO: [%s] Código: %s", sel.Tipo, sel.Codigo)),
			fmt.Sprintf("• Cliente: %s | E-mail: %s | Tel: %s", sel.NomeCliente, sel.EmailCliente, sel.TelefoneCliente),
			fmt.Sprintf("• Concessionária: %s | Horário: %s", sel.UnidadeNome, sel.DataHoraInicio.Format("02/01/2006 às 15:04")),
			fmt.Sprintf("• Veículo: %s | Situação Atual: %s", veicDesc, sel.Situacao),
		}
		previewCard = theme.CardStyle.Width(m.Width - 4).Render(strings.Join(info, "\n"))
	}

	paginacaoInfo := lipgloss.NewStyle().
		Foreground(theme.FgMuted).
		Render(fmt.Sprintf("Exibindo %d de %d agendamento(s) | [v] Agendar Visita  [t] Novo Test-Drive  [s] Alterar Status  [↑/↓] Navegar", len(m.Agendamentos), m.TotalItens))

	return lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.NewStyle().Bold(true).Foreground(theme.FgLight).Padding(1, 0).Render("AGENDA DE TEST-DRIVES E VISITAS"),
		renderedTable,
		paginacaoInfo,
		"",
		previewCard,
	)
}
