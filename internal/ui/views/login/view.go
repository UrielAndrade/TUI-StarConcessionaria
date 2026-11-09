package login

import (
	"fmt"

	"github.com/UrielAndrade/TUI-StarConcessionaria/internal/ui/theme"
	"github.com/charmbracelet/huh"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type SubmitLoginMsg struct {
	Email string
	Senha string
}

type Model struct {
	Form       *huh.Form
	Email      string
	Senha      string
	Carregando bool
	Erro       string
	Logado     bool
	Usuario    string
	Width      int
	Height     int
}

func New() Model {
	m := Model{}
	m.initForm()
	return m
}

func (m *Model) initForm() {
	m.Form = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("E-mail corporativo").
				Placeholder("usuario@starconcessionaria.com.br").
				Value(&m.Email).
				Validate(func(str string) error {
					if str == "" {
						return fmt.Errorf("o e-mail é obrigatório")
					}
					return nil
				}),

			huh.NewInput().
				Title("Senha de acesso").
				EchoMode(huh.EchoModePassword).
				Value(&m.Senha).
				Validate(func(str string) error {
					if str == "" {
						return fmt.Errorf("a senha é obrigatória")
					}
					return nil
				}),
		),
	).WithTheme(huh.ThemeCharm())
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if m.Carregando {
		return m, nil
	}

	form, cmd := m.Form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.Form = f
	}

	if m.Form.State == huh.StateCompleted {
		m.Carregando = true
		return m, func() tea.Msg {
			return SubmitLoginMsg{
				Email: m.Email,
				Senha: m.Senha,
			}
		}
	}

	return m, cmd
}

func (m Model) View() string {
	if m.Logado {
		box := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(theme.Success).
			Padding(1, 2).
			Render(fmt.Sprintf("✓ Sessão ativa como: %s\n\nVocê possui acesso completo aos recursos administrativos e métricas.", m.Usuario))
		return box
	}

	var statusHeader string
	if m.Carregando {
		statusHeader = lipgloss.NewStyle().Foreground(theme.Accent).Render("⏳ Autenticando com o servidor...")
	} else if m.Erro != "" {
		statusHeader = lipgloss.NewStyle().Foreground(theme.Danger).Render(fmt.Sprintf("Falha: %s", m.Erro))
	} else {
		statusHeader = lipgloss.NewStyle().Foreground(theme.FgMuted).Render("Informe suas credenciais para liberar recursos administrativos:")
	}

	formView := m.Form.View()

	box := theme.ModalBoxStyle.Width(50).Render(
		lipgloss.JoinVertical(lipgloss.Left,
			lipgloss.NewStyle().Bold(true).Foreground(theme.Primary).Render("AUTENTICAÇÃO DE USUÁRIO"),
			statusHeader,
			"",
			formView,
		),
	)

	return lipgloss.Place(m.Width, m.Height-4, lipgloss.Center, lipgloss.Center, box)
}
