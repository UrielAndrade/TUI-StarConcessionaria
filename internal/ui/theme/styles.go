package theme

import "github.com/charmbracelet/lipgloss"

var (
	// Estilos Globais
	AppStyle = lipgloss.NewStyle().
			Padding(0, 1)

	// Cabeçalho
	HeaderTitleStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(FgLight).
				Background(Primary).
				Padding(0, 1)

	HeaderSubStyle = lipgloss.NewStyle().
			Foreground(FgMuted).
			Padding(0, 1)

	TabActiveStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(FgLight).
			Background(Primary).
			Padding(0, 1)

	TabInactiveStyle = lipgloss.NewStyle().
				Foreground(FgMuted).
				Background(BgCard).
				Padding(0, 1)

	// Rodapé e Barra de Status
	FooterStyle = lipgloss.NewStyle().
			Foreground(FgMuted).
			Background(BgCard).
			Padding(0, 1)

	KeyStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Accent)

	// Cards do Dashboard
	CardStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(BorderCol).
			Padding(1, 2).
			MarginRight(1).
			Background(BgCard)

	CardTitleStyle = lipgloss.NewStyle().
			Foreground(FgMuted).
			Bold(true)

	CardValueStyle = lipgloss.NewStyle().
			Foreground(FgLight).
			Bold(true).
			MarginTop(1)

	// Tabela
	TableHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(FgLight).
				Background(Primary).
				Padding(0, 1)

	TableSelectedRowStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(FgLight).
				Background(lipgloss.Color("#2E3440"))

	// Modais e Diálogos
	ModalBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(Primary).
			Padding(1, 2).
			Background(BgCard)
)
