package dashboard

import (
	"fmt"
	"strings"

	"github.com/UrielAndrade/TUI-StarConcessionaria/internal/model"
	"github.com/UrielAndrade/TUI-StarConcessionaria/internal/ui/theme"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	Indicadores *model.IndicadoresGerenciaisResponse
	Carregando  bool
	Erro        string
	Width       int
	Height      int
}

func New() Model {
	return Model{
		Carregando: true,
	}
}

func (m Model) View() string {
	if m.Carregando {
		return lipgloss.NewStyle().
			Padding(4).
			Foreground(theme.Accent).
			Render("⏳ Carregando indicadores gerenciais da API...")
	}

	if m.Erro != "" {
		box := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(theme.Warning).
			Padding(1, 2).
			Render(fmt.Sprintf("Aviso: %s\n\n(Dica: Pressione [4] para efetuar login caso as métricas sejam restritas à gerência)", m.Erro))
		return box
	}

	if m.Indicadores == nil {
		return "Nenhum indicador disponível."
	}

	// Linha de Cards Principais
	card1 := theme.CardStyle.Width(22).Render(
		lipgloss.JoinVertical(lipgloss.Left,
			theme.CardTitleStyle.Render("DISPONÍVEIS"),
			theme.CardValueStyle.Foreground(theme.Success).Render(fmt.Sprintf("%d veíc.", m.Indicadores.TotalVeiculosDisponiveis)),
			lipgloss.NewStyle().Foreground(theme.FgMuted).Render("Prontos para venda"),
		),
	)

	card2 := theme.CardStyle.Width(22).Render(
		lipgloss.JoinVertical(lipgloss.Left,
			theme.CardTitleStyle.Render("RESERVADOS"),
			theme.CardValueStyle.Foreground(theme.Warning).Render(fmt.Sprintf("%d veíc.", m.Indicadores.TotalVeiculosReservados)),
			lipgloss.NewStyle().Foreground(theme.FgMuted).Render(fmt.Sprintf("%d ativas", m.Indicadores.TotalReservasAtivas)),
		),
	)

	card3 := theme.CardStyle.Width(22).Render(
		lipgloss.JoinVertical(lipgloss.Left,
			theme.CardTitleStyle.Render("VENDIDOS"),
			theme.CardValueStyle.Foreground(theme.Primary).Render(fmt.Sprintf("%d veíc.", m.Indicadores.TotalVeiculosVendidos)),
			lipgloss.NewStyle().Foreground(theme.FgMuted).Render(fmt.Sprintf("%d vendas total", m.Indicadores.TotalVendasVeiculos)),
		),
	)

	card4 := theme.CardStyle.Width(26).Render(
		lipgloss.JoinVertical(lipgloss.Left,
			theme.CardTitleStyle.Render("FATURAMENTO"),
			theme.CardValueStyle.Foreground(theme.Accent).Render(fmt.Sprintf("R$ %.2f", m.Indicadores.ValorTotalVendasVeiculos)),
			lipgloss.NewStyle().Foreground(theme.FgMuted).Render("Vendas de veículos"),
		),
	)

	cardsRow := lipgloss.JoinHorizontal(lipgloss.Top, card1, card2, card3, card4)

	// Painel de Veículos por Loja
	var linhasLojas []string
	linhasLojas = append(linhasLojas, lipgloss.NewStyle().Bold(true).Foreground(theme.Secondary).Render("ESTOQUE POR UNIDADE"))
	linhasLojas = append(linhasLojas, strings.Repeat("─", 40))
	for _, l := range m.Indicadores.VeiculosDisponiveisPorLoja {
		linhasLojas = append(linhasLojas, fmt.Sprintf("• %-28s %3d veíc.", l.LojaNome, l.Quantidade))
	}
	if len(m.Indicadores.VeiculosDisponiveisPorLoja) == 0 {
		linhasLojas = append(linhasLojas, "Nenhum dado por filial registrado.")
	}
	painelLojas := theme.CardStyle.Width(45).Render(strings.Join(linhasLojas, "\n"))

	// Painel de Peças Mais Vendidas
	var linhasPecas []string
	linhasPecas = append(linhasPecas, lipgloss.NewStyle().Bold(true).Foreground(theme.Secondary).Render("TOP PEÇAS MAIS VENDIDAS"))
	linhasPecas = append(linhasPecas, strings.Repeat("─", 40))
	for _, p := range m.Indicadores.TopPecasMaisVendidas {
		linhasPecas = append(linhasPecas, fmt.Sprintf("• %-28s %3d un.", p.PecaNome, p.TotalQuantidade))
	}
	if len(m.Indicadores.TopPecasMaisVendidas) == 0 {
		linhasPecas = append(linhasPecas, "Nenhuma venda de peça registrada.")
	}
	painelPecas := theme.CardStyle.Width(45).Render(strings.Join(linhasPecas, "\n"))

	paineisRow := lipgloss.JoinHorizontal(lipgloss.Top, painelLojas, painelPecas)

	return lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.NewStyle().Bold(true).Foreground(theme.FgLight).Padding(1, 0).Render("PANORAMA GERAL DA OPERAÇÃO"),
		cardsRow,
		"",
		paineisRow,
	)
}
