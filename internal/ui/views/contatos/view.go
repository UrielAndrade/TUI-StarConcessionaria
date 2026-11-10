package contatos

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
	EstadoNovoContato
	EstadoResponder
)

type Model struct {
	Table         table.Model
	Contatos      []model.ContatoResponse
	TotalItens    int
	PaginaAtual   int
	Carregando    bool
	Erro          string
	Width         int
	Height        int
	Estado        EstadoView
	NovoForm      *forms.NovoContatoForm
	ResponderForm *forms.ResponderContatoForm
	IsAutenticado bool
}

func New() Model {
	cols := []table.Column{
		{Title: "ID", Width: 4},
		{Title: "Protocolo", Width: 12},
		{Title: "Tipo", Width: 14},
		{Title: "Data", Width: 12},
		{Title: "Cliente", Width: 20},
		{Title: "Telefone", Width: 16},
		{Title: "Unidade", Width: 20},
		{Title: "Respondido", Width: 12},
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

func (m *Model) SetDados(resultado *model.PagedResult[model.ContatoResponse]) {
	m.Contatos = resultado.Itens
	m.TotalItens = resultado.TotalItens
	m.PaginaAtual = resultado.Pagina
	m.Carregando = false
	m.Erro = ""

	var rows []table.Row
	for _, c := range resultado.Itens {
		dataStr := c.Data.Format("02/01/2006")
		respStr := "Não"
		if c.Respondido {
			respStr = "Sim"
		}

		rows = append(rows, table.Row{
			fmt.Sprintf("%d", c.ID),
			c.Protocolo,
			c.Tipo,
			dataStr,
			c.NomeCliente,
			c.TelefoneCliente,
			c.UnidadeNome,
			respStr,
		})
	}
	m.Table.SetRows(rows)
}

func (m Model) ContatoSelecionado() *model.ContatoResponse {
	c := m.Table.Cursor()
	if c >= 0 && c < len(m.Contatos) {
		return &m.Contatos[c]
	}
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch m.Estado {
	case EstadoNovoContato:
		if m.NovoForm != nil {
			var cmd tea.Cmd
			m.NovoForm, cmd = m.NovoForm.Update(msg)
			return m, cmd
		}
	case EstadoResponder:
		if m.ResponderForm != nil {
			var cmd tea.Cmd
			m.ResponderForm, cmd = m.ResponderForm.Update(msg)
			return m, cmd
		}
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "a":
			m.Estado = EstadoNovoContato
			m.NovoForm = forms.NewNovoContatoForm()
			return m, nil

		case "r":
			if !m.IsAutenticado {
				m.Erro = "Autentique-se na aba [0] Perfil para responder chamados do SAC."
				return m, nil
			}
			sel := m.ContatoSelecionado()
			if sel != nil {
				m.Estado = EstadoResponder
				m.ResponderForm = forms.NewResponderContatoForm(sel.ID, sel.Protocolo)
				return m, nil
			}
		}
	}

	var cmd tea.Cmd
	m.Table, cmd = m.Table.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if m.Estado == EstadoNovoContato && m.NovoForm != nil {
		formBox := theme.ModalBoxStyle.Width(58).Render(
			lipgloss.JoinVertical(lipgloss.Left,
				lipgloss.NewStyle().Bold(true).Foreground(theme.Primary).Render("✉ NOVO ATENDIMENTO / CHAMADO SAC"),
				lipgloss.NewStyle().Foreground(theme.FgMuted).Render("Preencha a mensagem. [Esc] para cancelar."),
				"",
				m.NovoForm.View(),
			),
		)
		return lipgloss.Place(m.Width, m.Height-4, lipgloss.Center, lipgloss.Center, formBox)
	}

	if m.Estado == EstadoResponder && m.ResponderForm != nil {
		formBox := theme.ModalBoxStyle.Width(56).Render(
			lipgloss.JoinVertical(lipgloss.Left,
				lipgloss.NewStyle().Bold(true).Foreground(theme.Secondary).Render("📝 RESPONDER CHAMADO SAC"),
				lipgloss.NewStyle().Foreground(theme.FgMuted).Render("[Esc] para voltar sem responder."),
				"",
				m.ResponderForm.View(),
			),
		)
		return lipgloss.Place(m.Width, m.Height-4, lipgloss.Center, lipgloss.Center, formBox)
	}

	if m.Carregando {
		return lipgloss.NewStyle().
			Padding(4).
			Foreground(theme.Accent).
			Render("⏳ Carregando mensagens e chamados do SAC...")
	}

	if m.Erro != "" {
		return lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(theme.Warning).
			Padding(1, 2).
			Render(fmt.Sprintf("Aviso: %s\n\n(Dica: Mensagens de clientes exigem login na aba [0] Perfil)", m.Erro))
	}

	renderedTable := m.Table.View()

	cursor := m.Table.Cursor()
	var previewCard string
	if cursor >= 0 && cursor < len(m.Contatos) {
		sel := m.Contatos[cursor]
		info := []string{
			lipgloss.NewStyle().Bold(true).Foreground(theme.Secondary).Render(fmt.Sprintf("CHAMADO #%d - Protocolo: %s (Tipo: %s)", sel.ID, sel.Protocolo, sel.Tipo)),
			fmt.Sprintf("• Cliente: %s | E-mail: %s | Tel: %s", sel.NomeCliente, sel.EmailCliente, sel.TelefoneCliente),
			fmt.Sprintf("• Concessionária: %s | Data: %s", sel.UnidadeNome, sel.Data.Format("02/01/2006 às 15:04")),
			fmt.Sprintf("• Mensagem do Lead: \"%s\"", sel.Mensagem),
		}
		previewCard = theme.CardStyle.Width(m.Width - 4).Render(strings.Join(info, "\n"))
	}

	paginacaoInfo := lipgloss.NewStyle().
		Foreground(theme.FgMuted).
		Render(fmt.Sprintf("Exibindo %d de %d contato(s) | [a] Novo Atendimento  [r] Responder Chamado  [↑/↓] Navegar", len(m.Contatos), m.TotalItens))

	return lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.NewStyle().Bold(true).Foreground(theme.FgLight).Padding(1, 0).Render("ATENDIMENTO AO CLIENTE & LEADS (SAC)"),
		renderedTable,
		paginacaoInfo,
		"",
		previewCard,
	)
}
