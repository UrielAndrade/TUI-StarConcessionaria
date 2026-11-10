package propostas

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
	EstadoNovaProposta
	EstadoStatus
)

type Model struct {
	Table         table.Model
	Propostas     []model.PropostaResponse
	TotalItens    int
	PaginaAtual   int
	Carregando    bool
	Erro          string
	Width         int
	Height        int
	Estado        EstadoView
	NovaForm      *forms.NovaPropostaForm
	StatusForm    *forms.StatusPropostaForm
	IsAutenticado bool
}

func New() Model {
	cols := []table.Column{
		{Title: "ID", Width: 4},
		{Title: "Protocolo", Width: 12},
		{Title: "Veículo", Width: 10},
		{Title: "Filial", Width: 8},
		{Title: "Cliente", Width: 20},
		{Title: "Telefone", Width: 16},
		{Title: "Pagamento", Width: 15},
		{Title: "Entrada (R$)", Width: 14},
		{Title: "Troca?", Width: 8},
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

func (m *Model) SetDados(resultado *model.PagedResult[model.PropostaResponse]) {
	m.Propostas = resultado.Itens
	m.TotalItens = resultado.TotalItens
	m.PaginaAtual = resultado.Pagina
	m.Carregando = false
	m.Erro = ""

	var rows []table.Row
	for _, p := range resultado.Itens {
		trocaStr := "Não"
		if p.PossuiTroca {
			trocaStr = "Sim"
		}

		rows = append(rows, table.Row{
			fmt.Sprintf("%d", p.ID),
			p.Protocolo,
			fmt.Sprintf("#%d", p.VeiculoID),
			fmt.Sprintf("#%d", p.UnidadeID),
			p.NomeCliente,
			p.TelefoneCliente,
			p.FormaPagamento,
			fmt.Sprintf("%.2f", p.ValorEntrada),
			trocaStr,
		})
	}
	m.Table.SetRows(rows)
}

func (m Model) PropostaSelecionada() *model.PropostaResponse {
	c := m.Table.Cursor()
	if c >= 0 && c < len(m.Propostas) {
		return &m.Propostas[c]
	}
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch m.Estado {
	case EstadoNovaProposta:
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
			m.Estado = EstadoNovaProposta
			m.NovaForm = forms.NewNovaPropostaForm(1, 1)
			return m, nil

		case "s":
			if !m.IsAutenticado {
				m.Erro = "Autentique-se na aba [0] Perfil para responder propostas."
				return m, nil
			}
			sel := m.PropostaSelecionada()
			if sel != nil {
				m.Estado = EstadoStatus
				m.StatusForm = forms.NewStatusPropostaForm(sel.ID, sel.Protocolo)
				return m, nil
			}
		}
	}

	var cmd tea.Cmd
	m.Table, cmd = m.Table.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if m.Estado == EstadoNovaProposta && m.NovaForm != nil {
		formBox := theme.ModalBoxStyle.Width(58).Render(
			lipgloss.JoinVertical(lipgloss.Left,
				lipgloss.NewStyle().Bold(true).Foreground(theme.Primary).Render("📝 ENVIAR PROPOSTA COMERCIAL DE COMPRA"),
				lipgloss.NewStyle().Foreground(theme.FgMuted).Render("Preencha a proposta comercial. [Esc] para cancelar."),
				"",
				m.NovaForm.View(),
			),
		)
		return lipgloss.Place(m.Width, m.Height-4, lipgloss.Center, lipgloss.Center, formBox)
	}

	if m.Estado == EstadoStatus && m.StatusForm != nil {
		formBox := theme.ModalBoxStyle.Width(56).Render(
			lipgloss.JoinVertical(lipgloss.Left,
				lipgloss.NewStyle().Bold(true).Foreground(theme.Secondary).Render("💼 PARECER SOBRE A PROPOSTA"),
				lipgloss.NewStyle().Foreground(theme.FgMuted).Render("[Esc] para voltar sem salvar parecer."),
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
			Render("⏳ Carregando propostas de compra e financiamento...")
	}

	if m.Erro != "" {
		return lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(theme.Warning).
			Padding(1, 2).
			Render(fmt.Sprintf("Aviso: %s\n\n(Dica: Propostas exigem login de vendedor na aba [0] Perfil)", m.Erro))
	}

	renderedTable := m.Table.View()

	cursor := m.Table.Cursor()
	var previewCard string
	if cursor >= 0 && cursor < len(m.Propostas) {
		sel := m.Propostas[cursor]
		trocaDesc := "Não informada"
		if sel.PossuiTroca {
			trocaDesc = "Sim (Veículo ofertado na negociação)"
		}

		info := []string{
			lipgloss.NewStyle().Bold(true).Foreground(theme.Secondary).Render(fmt.Sprintf("PROPOSTA #%d - Protocolo: %s | Veículo: #%d | Concessionária: #%d", sel.ID, sel.Protocolo, sel.VeiculoID, sel.UnidadeID)),
			fmt.Sprintf("• Comprador: %s | E-mail: %s | Tel: %s", sel.NomeCliente, sel.EmailCliente, sel.TelefoneCliente),
			fmt.Sprintf("• Pagamento: %s | Entrada: R$ %.2f | Veículo de Troca: %s", sel.FormaPagamento, sel.ValorEntrada, trocaDesc),
			fmt.Sprintf("• Data de Envio: %s", sel.DataSolicitacao.Format("02/01/2006 às 15:04")),
		}
		previewCard = theme.CardStyle.Width(m.Width - 4).Render(strings.Join(info, "\n"))
	}

	paginacaoInfo := lipgloss.NewStyle().
		Foreground(theme.FgMuted).
		Render(fmt.Sprintf("Exibindo %d de %d proposta(s) | [a] Nova Proposta  [s] Dar Parecer  [↑/↓] Navegar", len(m.Propostas), m.TotalItens))

	return lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.NewStyle().Bold(true).Foreground(theme.FgLight).Padding(1, 0).Render("PROPOSTAS COMERCIAIS E NEGOCIAÇÕES"),
		renderedTable,
		paginacaoInfo,
		"",
		previewCard,
	)
}
