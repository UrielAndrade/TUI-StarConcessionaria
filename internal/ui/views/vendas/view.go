package vendas

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
	EstadoNovaVenda
)

type Model struct {
	Table         table.Model
	Vendas        []model.VendaResponse
	TotalItens    int
	PaginaAtual   int
	Carregando    bool
	Erro          string
	Width         int
	Height        int
	Estado        EstadoView
	NovaForm      *forms.NovaVendaForm
	IsAutenticado bool
}

func New() Model {
	cols := []table.Column{
		{Title: "ID", Width: 4},
		{Title: "Data", Width: 12},
		{Title: "Veículo", Width: 28},
		{Title: "Cliente", Width: 18},
		{Title: "Vendedor", Width: 16},
		{Title: "Unidade / Loja", Width: 22},
		{Title: "Pagamento", Width: 14},
		{Title: "Valor Final (R$)", Width: 16},
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

func (m *Model) SetDados(resultado *model.PagedResult[model.VendaResponse]) {
	m.Vendas = resultado.Itens
	m.TotalItens = resultado.TotalItens
	m.PaginaAtual = resultado.Pagina
	m.Carregando = false
	m.Erro = ""

	var rows []table.Row
	for _, v := range resultado.Itens {
		dataStr := v.DataVenda.Format("02/01/2006")
		rows = append(rows, table.Row{
			fmt.Sprintf("%d", v.ID),
			dataStr,
			v.VeiculoDescricao,
			v.NomeCliente,
			v.NomeVendedor,
			v.NomeLoja,
			v.FormaPagamentoNome,
			fmt.Sprintf("%.2f", v.ValorFinal),
		})
	}
	m.Table.SetRows(rows)
}

func (m Model) VendaSelecionada() *model.VendaResponse {
	c := m.Table.Cursor()
	if c >= 0 && c < len(m.Vendas) {
		return &m.Vendas[c]
	}
	return nil
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if m.Estado == EstadoNovaVenda && m.NovaForm != nil {
		var cmd tea.Cmd
		m.NovaForm, cmd = m.NovaForm.Update(msg)
		return m, cmd
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "a":
			if !m.IsAutenticado {
				m.Erro = "Autentique-se na aba [9] Perfil para registrar vendas."
				return m, nil
			}
			m.Estado = EstadoNovaVenda
			m.NovaForm = forms.NewNovaVendaForm(1, 0)
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.Table, cmd = m.Table.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if m.Estado == EstadoNovaVenda && m.NovaForm != nil {
		formBox := theme.ModalBoxStyle.Width(58).Render(
			lipgloss.JoinVertical(lipgloss.Left,
				lipgloss.NewStyle().Bold(true).Foreground(theme.Success).Render("💰 REGISTRAR VENDA DE VEÍCULO"),
				lipgloss.NewStyle().Foreground(theme.FgMuted).Render("Preencha os dados da negociação. [Esc] para cancelar."),
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
			Render("⏳ Carregando histórico de vendas de veículos...")
	}

	if m.Erro != "" {
		return lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(theme.Warning).
			Padding(1, 2).
			Render(fmt.Sprintf("Aviso: %s\n\n(Dica: As vendas consolidadas exigem login de vendedor/gerente na aba [9] Perfil)", m.Erro))
	}

	renderedTable := m.Table.View()

	cursor := m.Table.Cursor()
	var previewCard string
	if cursor >= 0 && cursor < len(m.Vendas) {
		sel := m.Vendas[cursor]
		info := []string{
			lipgloss.NewStyle().Bold(true).Foreground(theme.Success).Render(fmt.Sprintf("COMPROVANTE DE VENDA: #%d - %s", sel.ID, sel.VeiculoDescricao)),
			fmt.Sprintf("• Comprador: %s | Vendedor Responsável: %s", sel.NomeCliente, sel.NomeVendedor),
			fmt.Sprintf("• Filial: %s | Data da Operação: %s", sel.NomeLoja, sel.DataVenda.Format("02/01/2006 às 15:04")),
			fmt.Sprintf("• Forma de Pagamento: %s | Valor Total: R$ %.2f", sel.FormaPagamentoNome, sel.ValorFinal),
		}
		previewCard = theme.CardStyle.Width(m.Width - 4).Render(strings.Join(info, "\n"))
	}

	paginacaoInfo := lipgloss.NewStyle().
		Foreground(theme.FgMuted).
		Render(fmt.Sprintf("Exibindo %d de %d venda(s) realizada(s) | [a] Nova Venda  [↑/↓] Navegar", len(m.Vendas), m.TotalItens))

	return lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.NewStyle().Bold(true).Foreground(theme.FgLight).Padding(1, 0).Render("RELATÓRIO DE VENDAS DE VEÍCULOS"),
		renderedTable,
		paginacaoInfo,
		"",
		previewCard,
	)
}
