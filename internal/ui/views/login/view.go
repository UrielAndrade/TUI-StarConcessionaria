package login

import (
	"fmt"
	"strings"

	"github.com/UrielAndrade/TUI-StarConcessionaria/internal/ui/theme"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type SubmitLoginMsg struct {
	Email string
	Senha string
}

type LogoutMsg struct{}

type Model struct {
	focusIndex int
	inputs     []textinput.Model
	Carregando bool
	Erro       string
	Logado     bool
	Usuario    string
	Width      int
	Height     int
}

func New() Model {
	m := Model{
		inputs: make([]textinput.Model, 2),
	}

	// Campo de E-mail
	email := textinput.New()
	email.Placeholder = "admin@starconcessionaria.com.br"
	email.Focus()
	email.CharLimit = 64
	email.Width = 36
	m.inputs[0] = email

	// Campo de Senha
	senha := textinput.New()
	senha.Placeholder = "••••••••"
	senha.EchoMode = textinput.EchoPassword
	senha.EchoCharacter = '•'
	senha.CharLimit = 32
	senha.Width = 36
	m.inputs[1] = senha

	m.focusIndex = 0
	return m
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if m.Carregando {
		return m, nil
	}

	// Se já estiver logado, atende a comandos de perfil e logout
	if m.Logado {
		if keyMsg, ok := msg.(tea.KeyMsg); ok {
			switch keyMsg.String() {
			case "l", "L":
				m.Logado = false
				m.Usuario = ""
				m.inputs[0].SetValue("")
				m.inputs[1].SetValue("")
				m.focusIndex = 0
				m.inputs[0].Focus()
				return m, func() tea.Msg {
					return LogoutMsg{}
				}
			}
		}
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "shift+tab", "up", "down":
			s := msg.String()

			if s == "up" || s == "shift+tab" {
				m.focusIndex--
			} else {
				m.focusIndex++
			}

			// 0 = email, 1 = senha, 2 = botão entrar
			if m.focusIndex > 2 {
				m.focusIndex = 0
			} else if m.focusIndex < 0 {
				m.focusIndex = 2
			}

			cmds := make([]tea.Cmd, len(m.inputs))
			for i := 0; i < len(m.inputs); i++ {
				if i == m.focusIndex {
					cmds[i] = m.inputs[i].Focus()
				} else {
					m.inputs[i].Blur()
				}
			}
			return m, tea.Batch(cmds...)

		case "enter":
			if m.focusIndex == 0 {
				m.focusIndex = 1
				m.inputs[0].Blur()
				return m, m.inputs[1].Focus()
			}

			// Submeter login
			emailVal := strings.TrimSpace(m.inputs[0].Value())
			senhaVal := strings.TrimSpace(m.inputs[1].Value())

			if emailVal == "" {
				m.Erro = "O e-mail é obrigatório."
				m.focusIndex = 0
				m.inputs[1].Blur()
				return m, m.inputs[0].Focus()
			}
			if senhaVal == "" {
				m.Erro = "A senha é obrigatória."
				m.focusIndex = 1
				m.inputs[0].Blur()
				return m, m.inputs[1].Focus()
			}

			m.Carregando = true
			m.Erro = ""
			return m, func() tea.Msg {
				return SubmitLoginMsg{
					Email: emailVal,
					Senha: senhaVal,
				}
			}
		}
	}

	// Atualiza os inputs de texto
	cmds := make([]tea.Cmd, len(m.inputs))
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.Logado {
		infoLines := []string{
			lipgloss.NewStyle().Bold(true).Foreground(theme.Success).Render("✓ SESSÃO ATIVA"),
			"",
			fmt.Sprintf("Usuário Conectado: %s", lipgloss.NewStyle().Bold(true).Foreground(theme.FgLight).Render(m.Usuario)),
			"Status: Autenticado com Token JWT Bearer",
			"Permissões: Acesso completo aos recursos gerenciais e métricas",
			"",
			lipgloss.NewStyle().Foreground(theme.FgMuted).Render("Atalhos rápidos para navegação:"),
			lipgloss.NewStyle().Foreground(theme.Accent).Render("  [1] Ir para o Dashboard de Indicadores"),
			lipgloss.NewStyle().Foreground(theme.Accent).Render("  [2] Acessar Catálogo de Veículos"),
			lipgloss.NewStyle().Foreground(theme.Accent).Render("  [3] Acessar Estoque de Peças"),
			lipgloss.NewStyle().Foreground(theme.Danger).Render("  [L] Encerrar Sessão (Logout)"),
			lipgloss.NewStyle().Foreground(theme.FgMuted).Render("  [Esc] Voltar ao Dashboard"),
		}
		box := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(theme.Success).
			Padding(1, 3).
			Width(56).
			Render(strings.Join(infoLines, "\n"))
		return lipgloss.Place(m.Width, m.Height-4, lipgloss.Center, lipgloss.Center, box)
	}

	var statusHeader string
	if m.Carregando {
		statusHeader = lipgloss.NewStyle().Foreground(theme.Accent).Render("⏳ Autenticando com o servidor...")
	} else if m.Erro != "" {
		statusHeader = lipgloss.NewStyle().Foreground(theme.Danger).Bold(true).Render(fmt.Sprintf("✕ %s", m.Erro))
	} else {
		statusHeader = lipgloss.NewStyle().Foreground(theme.FgMuted).Render("Informe suas credenciais para liberar recursos administrativos:")
	}

	emailLabel := lipgloss.NewStyle().Bold(true).Render("E-mail corporativo:")
	senhaLabel := lipgloss.NewStyle().Bold(true).Render("Senha de acesso:")

	btnStyle := lipgloss.NewStyle().
		Padding(0, 3).
		Border(lipgloss.RoundedBorder())

	if m.focusIndex == 2 {
		btnStyle = btnStyle.
			Foreground(theme.FgLight).
			Background(theme.Primary).
			BorderForeground(theme.Primary).
			Bold(true)
	} else {
		btnStyle = btnStyle.
			Foreground(theme.FgMuted).
			BorderForeground(theme.BorderCol)
	}
	btnEntrar := btnStyle.Render("Entrar")

	formContent := lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.NewStyle().Bold(true).Foreground(theme.Primary).Render("★ AUTENTICAÇÃO DE USUÁRIO"),
		statusHeader,
		"",
		emailLabel,
		m.inputs[0].View(),
		"",
		senhaLabel,
		m.inputs[1].View(),
		"",
		btnEntrar,
		"",
		lipgloss.NewStyle().Foreground(theme.FgMuted).Render("[Tab / ↑ / ↓] Alternar campos  •  [Enter] Avançar / Entrar  •  [Esc] Voltar"),
	)

	box := theme.ModalBoxStyle.Width(52).Render(formContent)

	return lipgloss.Place(m.Width, m.Height-4, lipgloss.Center, lipgloss.Center, box)
}
