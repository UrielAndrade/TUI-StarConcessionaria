package ui

import (
	"context"

	"github.com/UrielAndrade/TUI-StarConcessionaria/internal/api"
	"github.com/UrielAndrade/TUI-StarConcessionaria/internal/config"
	"github.com/UrielAndrade/TUI-StarConcessionaria/internal/model"
	"github.com/UrielAndrade/TUI-StarConcessionaria/internal/ui/components"
	"github.com/UrielAndrade/TUI-StarConcessionaria/internal/ui/views/dashboard"
	"github.com/UrielAndrade/TUI-StarConcessionaria/internal/ui/views/login"
	"github.com/UrielAndrade/TUI-StarConcessionaria/internal/ui/views/pecas"
	"github.com/UrielAndrade/TUI-StarConcessionaria/internal/ui/views/veiculos"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Mensagens assíncronas do Bubble Tea
type IndicadoresLoadedMsg struct {
	Data *model.IndicadoresGerenciaisResponse
	Err  error
}

type VeiculosLoadedMsg struct {
	Data *model.PagedResult[model.VeiculoCatalogoResponse]
	Err  error
}

type PecasLoadedMsg struct {
	Data *model.PagedResult[model.PecaResponse]
	Err  error
}

type LoginResultMsg struct {
	Resp *model.LoginResponse
	Err  error
}

type AppModel struct {
	client      *api.Client
	cfg         *config.Config
	activeTab   int
	width       int
	height      int
	isConnected bool
	statusMsg   string

	// Telas
	dashboardView dashboard.Model
	veiculosView  veiculos.Model
	pecasView     pecas.Model
	loginView     login.Model
}

func NewAppModel(client *api.Client, cfg *config.Config) AppModel {
	return AppModel{
		client:        client,
		cfg:           cfg,
		activeTab:     1, // Inicia na aba de Veículos (pública)
		statusMsg:     "Pronto",
		dashboardView: dashboard.New(),
		veiculosView:  veiculos.New(),
		pecasView:     pecas.New(),
		loginView:     login.New(),
	}
}

func (m AppModel) Init() tea.Cmd {
	return tea.Batch(
		m.loadVeiculosCmd(),
		m.loadPecasCmd(),
		m.loadDashboardCmd(),
	)
}

func (m AppModel) loadDashboardCmd() tea.Cmd {
	return func() tea.Msg {
		data, err := m.client.GetIndicadores(context.Background(), nil)
		return IndicadoresLoadedMsg{Data: data, Err: err}
	}
}

func (m AppModel) loadVeiculosCmd() tea.Cmd {
	return func() tea.Msg {
		data, err := m.client.ListVeiculos(context.Background(), model.FiltroVeiculos{
			Pagina:         1,
			ItensPorPagina: 12,
		})
		return VeiculosLoadedMsg{Data: data, Err: err}
	}
}

func (m AppModel) loadPecasCmd() tea.Cmd {
	return func() tea.Msg {
		data, err := m.client.ListPecas(context.Background(), model.FiltroPecas{
			Pagina:         1,
			ItensPorPagina: 12,
		})
		return PecasLoadedMsg{Data: data, Err: err}
	}
}

func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.dashboardView.Width = msg.Width
		m.dashboardView.Height = msg.Height
		m.veiculosView.Width = msg.Width
		m.veiculosView.Height = msg.Height
		m.pecasView.Width = msg.Width
		m.pecasView.Height = msg.Height
		m.loginView.Width = msg.Width
		m.loginView.Height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			// Se estiver editando formulário de login, não sai com 'q'
			if m.activeTab != 3 || msg.String() == "ctrl+c" {
				return m, tea.Quit
			}
		case "1":
			m.activeTab = 0
		case "2":
			m.activeTab = 1
		case "3":
			m.activeTab = 2
		case "4":
			m.activeTab = 3
		case "r":
			// Recarregar dados da aba atual
			switch m.activeTab {
			case 0:
				m.dashboardView.Carregando = true
				cmds = append(cmds, m.loadDashboardCmd())
			case 1:
				m.veiculosView.Carregando = true
				cmds = append(cmds, m.loadVeiculosCmd())
			case 2:
				m.pecasView.Carregando = true
				cmds = append(cmds, m.loadPecasCmd())
			}
		case "tab":
			m.activeTab = (m.activeTab + 1) % 4
		}

	case IndicadoresLoadedMsg:
		m.dashboardView.Carregando = false
		if msg.Err != nil {
			m.dashboardView.Erro = msg.Err.Error()
		} else {
			m.dashboardView.Indicadores = msg.Data
			m.dashboardView.Erro = ""
			m.isConnected = true
		}

	case VeiculosLoadedMsg:
		m.veiculosView.Carregando = false
		if msg.Err != nil {
			m.veiculosView.Erro = msg.Err.Error()
		} else {
			m.veiculosView.SetDados(msg.Data)
			m.isConnected = true
		}

	case PecasLoadedMsg:
		m.pecasView.Carregando = false
		if msg.Err != nil {
			m.pecasView.Erro = msg.Err.Error()
		} else {
			m.pecasView.SetDados(msg.Data)
			m.isConnected = true
		}

	case login.SubmitLoginMsg:
		m.loginView.Carregando = true
		cmds = append(cmds, func() tea.Msg {
			resp, err := m.client.Login(context.Background(), msg.Email, msg.Senha)
			return LoginResultMsg{Resp: resp, Err: err}
		})

	case LoginResultMsg:
		m.loginView.Carregando = false
		if msg.Err != nil {
			m.loginView.Erro = msg.Err.Error()
		} else {
			m.loginView.Logado = true
			m.loginView.Usuario = msg.Resp.Usuario.Nome
			m.statusMsg = "Usuário: " + msg.Resp.Usuario.Nome
			// Atualiza configuração e recarrega métricas restritas
			m.cfg.Token = msg.Resp.Token
			_ = m.cfg.Save()
			cmds = append(cmds, m.loadDashboardCmd())
		}
	}

	// Encaminha eventos para a visão ativa
	switch m.activeTab {
	case 1:
		var cmd tea.Cmd
		m.veiculosView, cmd = m.veiculosView.Update(msg)
		cmds = append(cmds, cmd)
	case 2:
		var cmd tea.Cmd
		m.pecasView, cmd = m.pecasView.Update(msg)
		cmds = append(cmds, cmd)
	case 3:
		var cmd tea.Cmd
		m.loginView, cmd = m.loginView.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m AppModel) View() string {
	tabs := []components.TabItem{
		{Key: "1", Title: "Dashboard"},
		{Key: "2", Title: "Veículos"},
		{Key: "3", Title: "Peças"},
		{Key: "4", Title: "Login"},
	}

	header := components.RenderHeader(m.activeTab, tabs, m.width)

	var body string
	switch m.activeTab {
	case 0:
		body = m.dashboardView.View()
	case 1:
		body = m.veiculosView.View()
	case 2:
		body = m.pecasView.View()
	case 3:
		body = m.loginView.View()
	}

	shortcuts := []components.Shortcut{
		{Key: "1-4", Desc: "Mudar Aba"},
		{Key: "Tab", Desc: "Próxima Aba"},
		{Key: "r", Desc: "Recarregar"},
		{Key: "↑/↓", Desc: "Navegar"},
		{Key: "q", Desc: "Sair"},
	}

	footer := components.RenderFooter(shortcuts, m.statusMsg, m.isConnected, m.width)

	// Layout vertical completo
	contentHeight := m.height - lipgloss.Height(header) - lipgloss.Height(footer)
	if contentHeight < 5 {
		contentHeight = 5
	}

	styledBody := lipgloss.NewStyle().
		Height(contentHeight).
		Width(m.width).
		Render(body)

	return lipgloss.JoinVertical(lipgloss.Left, header, styledBody, footer)
}
