package ui

import (
	"context"

	"github.com/UrielAndrade/TUI-StarConcessionaria/internal/api"
	"github.com/UrielAndrade/TUI-StarConcessionaria/internal/config"
	"github.com/UrielAndrade/TUI-StarConcessionaria/internal/model"
	"github.com/UrielAndrade/TUI-StarConcessionaria/internal/ui/components"
	"github.com/UrielAndrade/TUI-StarConcessionaria/internal/ui/forms"
	"github.com/UrielAndrade/TUI-StarConcessionaria/internal/ui/views/agendamentos"
	"github.com/UrielAndrade/TUI-StarConcessionaria/internal/ui/views/contatos"
	"github.com/UrielAndrade/TUI-StarConcessionaria/internal/ui/views/dashboard"
	"github.com/UrielAndrade/TUI-StarConcessionaria/internal/ui/views/login"
	"github.com/UrielAndrade/TUI-StarConcessionaria/internal/ui/views/pecas"
	"github.com/UrielAndrade/TUI-StarConcessionaria/internal/ui/views/reservas"
	"github.com/UrielAndrade/TUI-StarConcessionaria/internal/ui/views/unidades"
	"github.com/UrielAndrade/TUI-StarConcessionaria/internal/ui/views/veiculos"
	"github.com/UrielAndrade/TUI-StarConcessionaria/internal/ui/views/vendas"
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

type VendasLoadedMsg struct {
	Data *model.PagedResult[model.VendaResponse]
	Err  error
}

type AgendamentosLoadedMsg struct {
	Data *model.PagedResult[model.AgendamentoResponse]
	Err  error
}

type ReservasLoadedMsg struct {
	Data *model.PagedResult[model.ReservaResponse]
	Err  error
}

type ContatosLoadedMsg struct {
	Data *model.PagedResult[model.ContatoResponse]
	Err  error
}

type UnidadesLoadedMsg struct {
	Data []model.UnidadeResponse
	Err  error
}

type LoginResultMsg struct {
	Resp *model.LoginResponse
	Err  error
}

type PerfilLoadedMsg struct {
	Usuario *model.UsuarioResponse
	Err     error
}

type VeiculoMutacaoResultMsg struct {
	Mensagem string
	Err      error
}

type PecaMutacaoResultMsg struct {
	Mensagem string
	Err      error
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
	dashboardView   dashboard.Model
	veiculosView    veiculos.Model
	pecasView       pecas.Model
	vendasView      vendas.Model
	agendaView      agendamentos.Model
	reservasView    reservas.Model
	contatosView    contatos.Model
	unidadesView    unidades.Model
	loginView       login.Model
}

func NewAppModel(client *api.Client, cfg *config.Config) AppModel {
	return AppModel{
		client:        client,
		cfg:           cfg,
		activeTab:     0, // Inicia no Dashboard
		statusMsg:     "Pronto",
		dashboardView: dashboard.New(),
		veiculosView:  veiculos.New(),
		pecasView:     pecas.New(),
		vendasView:    vendas.New(),
		agendaView:    agendamentos.New(),
		reservasView:  reservas.New(),
		contatosView:  contatos.New(),
		unidadesView:  unidades.New(),
		loginView:     login.New(),
	}
}

func (m AppModel) Init() tea.Cmd {
	cmds := []tea.Cmd{
		m.loadDashboardCmd(),
		m.loadVeiculosCmd(),
		m.loadPecasCmd(),
		m.loadUnidadesCmd(),
		m.loginView.Init(),
	}

	if m.cfg.Token != "" {
		cmds = append(cmds,
			m.loadPerfilCmd(),
			m.loadVendasCmd(),
			m.loadAgendamentosCmd(),
			m.loadReservasCmd(),
			m.loadContatosCmd(),
		)
	}

	return tea.Batch(cmds...)
}

func (m AppModel) loadPerfilCmd() tea.Cmd {
	return func() tea.Msg {
		usuario, err := m.client.GetPerfil(context.Background())
		return PerfilLoadedMsg{Usuario: usuario, Err: err}
	}
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

func (m AppModel) loadVendasCmd() tea.Cmd {
	return func() tea.Msg {
		data, err := m.client.ListVendas(context.Background(), model.FiltroVendas{
			Pagina:         1,
			ItensPorPagina: 12,
		})
		return VendasLoadedMsg{Data: data, Err: err}
	}
}

func (m AppModel) loadAgendamentosCmd() tea.Cmd {
	return func() tea.Msg {
		data, err := m.client.ListAgendamentos(context.Background(), model.FiltroAgendamentos{
			Pagina:         1,
			ItensPorPagina: 12,
		})
		return AgendamentosLoadedMsg{Data: data, Err: err}
	}
}

func (m AppModel) loadReservasCmd() tea.Cmd {
	return func() tea.Msg {
		data, err := m.client.ListReservas(context.Background(), model.FiltroReservas{
			Pagina:         1,
			ItensPorPagina: 12,
		})
		return ReservasLoadedMsg{Data: data, Err: err}
	}
}

func (m AppModel) loadContatosCmd() tea.Cmd {
	return func() tea.Msg {
		data, err := m.client.ListContatos(context.Background(), model.FiltroContatos{
			Pagina:         1,
			ItensPorPagina: 12,
		})
		return ContatosLoadedMsg{Data: data, Err: err}
	}
}

func (m AppModel) loadUnidadesCmd() tea.Cmd {
	return func() tea.Msg {
		data, err := m.client.ListUnidades(context.Background(), "", "")
		return UnidadesLoadedMsg{Data: data, Err: err}
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
		m.vendasView.Width = msg.Width
		m.vendasView.Height = msg.Height
		m.agendaView.Width = msg.Width
		m.agendaView.Height = msg.Height
		m.reservasView.Width = msg.Width
		m.reservasView.Height = msg.Height
		m.contatosView.Width = msg.Width
		m.contatosView.Height = msg.Height
		m.unidadesView.Width = msg.Width
		m.unidadesView.Height = msg.Height
		m.loginView.Width = msg.Width
		m.loginView.Height = msg.Height

	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

		// Se estiver na aba de veículos com modal interativo aberto:
		if m.activeTab == 1 && m.veiculosView.Estado != veiculos.EstadoLista {
			var cmd tea.Cmd
			m.veiculosView, cmd = m.veiculosView.Update(msg)
			return m, cmd
		}

		// Se estiver na aba de peças com modal interativo aberto:
		if m.activeTab == 2 && m.pecasView.Estado != pecas.EstadoLista {
			var cmd tea.Cmd
			m.pecasView, cmd = m.pecasView.Update(msg)
			return m, cmd
		}

		// Se estiver na tela de login e NÃO estiver autenticado:
		if m.activeTab == 8 && !m.loginView.Logado {
			if msg.String() == "esc" {
				m.activeTab = 0 // Volta para o Dashboard
				return m, nil
			}
			var cmd tea.Cmd
			m.loginView, cmd = m.loginView.Update(msg)
			return m, cmd
		}

		// Navegação global entre as 9 abas:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "1":
			m.activeTab = 0
		case "2":
			m.activeTab = 1
		case "3":
			m.activeTab = 2
		case "4":
			m.activeTab = 3
			if len(m.vendasView.Vendas) == 0 {
				m.vendasView.Carregando = true
				cmds = append(cmds, m.loadVendasCmd())
			}
		case "5":
			m.activeTab = 4
			if len(m.agendaView.Agendamentos) == 0 {
				m.agendaView.Carregando = true
				cmds = append(cmds, m.loadAgendamentosCmd())
			}
		case "6":
			m.activeTab = 5
			if len(m.reservasView.Reservas) == 0 {
				m.reservasView.Carregando = true
				cmds = append(cmds, m.loadReservasCmd())
			}
		case "7":
			m.activeTab = 6
			if len(m.contatosView.Contatos) == 0 {
				m.contatosView.Carregando = true
				cmds = append(cmds, m.loadContatosCmd())
			}
		case "8":
			m.activeTab = 7
			if len(m.unidadesView.Unidades) == 0 {
				m.unidadesView.Carregando = true
				cmds = append(cmds, m.loadUnidadesCmd())
			}
		case "9":
			m.activeTab = 8
		case "r":
			// Recarregar dados da tela ativa
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
			case 3:
				m.vendasView.Carregando = true
				cmds = append(cmds, m.loadVendasCmd())
			case 4:
				m.agendaView.Carregando = true
				cmds = append(cmds, m.loadAgendamentosCmd())
			case 5:
				m.reservasView.Carregando = true
				cmds = append(cmds, m.loadReservasCmd())
			case 6:
				m.contatosView.Carregando = true
				cmds = append(cmds, m.loadContatosCmd())
			case 7:
				m.unidadesView.Carregando = true
				cmds = append(cmds, m.loadUnidadesCmd())
			}
		case "tab":
			m.activeTab = (m.activeTab + 1) % 9
		case "esc":
			if m.activeTab != 0 {
				m.activeTab = 0
			}
		}

		// Se estiver na aba de perfil, permite logout via 'L'
		if m.activeTab == 8 && m.loginView.Logado {
			var cmd tea.Cmd
			m.loginView, cmd = m.loginView.Update(msg)
			cmds = append(cmds, cmd)
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

	case VendasLoadedMsg:
		m.vendasView.Carregando = false
		if msg.Err != nil {
			m.vendasView.Erro = msg.Err.Error()
		} else {
			m.vendasView.SetDados(msg.Data)
			m.isConnected = true
		}

	case AgendamentosLoadedMsg:
		m.agendaView.Carregando = false
		if msg.Err != nil {
			m.agendaView.Erro = msg.Err.Error()
		} else {
			m.agendaView.SetDados(msg.Data)
			m.isConnected = true
		}

	case ReservasLoadedMsg:
		m.reservasView.Carregando = false
		if msg.Err != nil {
			m.reservasView.Erro = msg.Err.Error()
		} else {
			m.reservasView.SetDados(msg.Data)
			m.isConnected = true
		}

	case ContatosLoadedMsg:
		m.contatosView.Carregando = false
		if msg.Err != nil {
			m.contatosView.Erro = msg.Err.Error()
		} else {
			m.contatosView.SetDados(msg.Data)
			m.isConnected = true
		}

	case UnidadesLoadedMsg:
		m.unidadesView.Carregando = false
		if msg.Err != nil {
			m.unidadesView.Erro = msg.Err.Error()
		} else {
			m.unidadesView.SetDados(msg.Data)
			m.isConnected = true
		}

	case forms.NovoVeiculoResultMsg:
		m.veiculosView.Estado = veiculos.EstadoLista
		m.veiculosView.NovoForm = nil
		if !msg.Cancel && msg.Payload != nil {
			m.statusMsg = "Cadastrando veículo no estoque..."
			cmds = append(cmds, func() tea.Msg {
				_, err := m.client.CriarVeiculo(context.Background(), *msg.Payload)
				return VeiculoMutacaoResultMsg{
					Mensagem: "Veículo cadastrado com sucesso!",
					Err:      err,
				}
			})
		}

	case forms.AlterarSituacaoResultMsg:
		m.veiculosView.Estado = veiculos.EstadoLista
		m.veiculosView.SituacaoForm = nil
		if !msg.Cancel {
			m.statusMsg = "Atualizando situação do veículo..."
			cmds = append(cmds, func() tea.Msg {
				err := m.client.AlterarSituacaoVeiculo(context.Background(), msg.VeiculoID, msg.NovaSituacao, msg.Justificativa)
				return VeiculoMutacaoResultMsg{
					Mensagem: "Situação do veículo alterada para " + msg.NovaSituacao + "!",
					Err:      err,
				}
			})
		}

	case VeiculoMutacaoResultMsg:
		if msg.Err != nil {
			m.veiculosView.Erro = msg.Err.Error()
			m.statusMsg = "Falha: " + msg.Err.Error()
		} else {
			m.veiculosView.Erro = ""
			m.statusMsg = msg.Mensagem
			m.veiculosView.Carregando = true
			cmds = append(cmds, m.loadVeiculosCmd(), m.loadDashboardCmd())
		}

	case forms.NovaPecaResultMsg:
		m.pecasView.Estado = pecas.EstadoLista
		m.pecasView.NovaForm = nil
		if !msg.Cancel && msg.Payload != nil {
			m.statusMsg = "Cadastrando peça no estoque..."
			cmds = append(cmds, func() tea.Msg {
				_, err := m.client.CriarPeca(context.Background(), *msg.Payload)
				return PecaMutacaoResultMsg{
					Mensagem: "Peça cadastrada com sucesso!",
					Err:      err,
				}
			})
		}

	case forms.VenderPecaResultMsg:
		m.pecasView.Estado = pecas.EstadoLista
		m.pecasView.VenderForm = nil
		if !msg.Cancel && msg.Payload != nil {
			m.statusMsg = "Registrando venda de peça..."
			cmds = append(cmds, func() tea.Msg {
				err := m.client.VenderPeca(context.Background(), msg.PecaID, *msg.Payload)
				return PecaMutacaoResultMsg{
					Mensagem: "Venda de peça registrada com sucesso!",
					Err:      err,
				}
			})
		}

	case PecaMutacaoResultMsg:
		if msg.Err != nil {
			m.pecasView.Erro = msg.Err.Error()
			m.statusMsg = "Falha: " + msg.Err.Error()
		} else {
			m.pecasView.Erro = ""
			m.statusMsg = msg.Mensagem
			m.pecasView.Carregando = true
			cmds = append(cmds, m.loadPecasCmd(), m.loadDashboardCmd())
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
			m.veiculosView.IsAutenticado = true
			m.pecasView.IsAutenticado = true
			m.statusMsg = "Usuário: " + msg.Resp.Usuario.Nome
			m.cfg.Token = msg.Resp.Token
			_ = m.cfg.Save()

			// Transita automaticamente para o Dashboard com os dados liberados!
			m.activeTab = 0
			m.dashboardView.Carregando = true
			cmds = append(cmds,
				m.loadDashboardCmd(),
				m.loadVendasCmd(),
				m.loadAgendamentosCmd(),
				m.loadReservasCmd(),
				m.loadContatosCmd(),
			)
		}

	case PerfilLoadedMsg:
		if msg.Err == nil && msg.Usuario != nil {
			m.loginView.Logado = true
			m.loginView.Usuario = msg.Usuario.Nome
			m.veiculosView.IsAutenticado = true
			m.pecasView.IsAutenticado = true
			m.statusMsg = "Usuário: " + msg.Usuario.Nome
			cmds = append(cmds,
				m.loadDashboardCmd(),
				m.loadVendasCmd(),
				m.loadAgendamentosCmd(),
				m.loadReservasCmd(),
				m.loadContatosCmd(),
			)
		}

	case login.LogoutMsg:
		m.cfg.Token = ""
		_ = m.cfg.Save()
		m.client.SetToken("")
		m.statusMsg = "Pronto (Desconectado)"
		m.activeTab = 0
		m.veiculosView.IsAutenticado = false
		m.pecasView.IsAutenticado = false
		m.dashboardView.Indicadores = nil
		m.dashboardView.Erro = "Efetue login para visualizar dados restritos"
		cmds = append(cmds, m.loadDashboardCmd())
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
		m.vendasView, cmd = m.vendasView.Update(msg)
		cmds = append(cmds, cmd)
	case 4:
		var cmd tea.Cmd
		m.agendaView, cmd = m.agendaView.Update(msg)
		cmds = append(cmds, cmd)
	case 5:
		var cmd tea.Cmd
		m.reservasView, cmd = m.reservasView.Update(msg)
		cmds = append(cmds, cmd)
	case 6:
		var cmd tea.Cmd
		m.contatosView, cmd = m.contatosView.Update(msg)
		cmds = append(cmds, cmd)
	case 7:
		var cmd tea.Cmd
		m.unidadesView, cmd = m.unidadesView.Update(msg)
		cmds = append(cmds, cmd)
	case 8:
		if !m.loginView.Logado {
			var cmd tea.Cmd
			m.loginView, cmd = m.loginView.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m AppModel) View() string {
	loginTabTitle := "Login"
	if m.loginView.Logado {
		loginTabTitle = "Perfil"
	}

	tabs := []components.TabItem{
		{Key: "1", Title: "Dashboard"},
		{Key: "2", Title: "Veículos"},
		{Key: "3", Title: "Peças"},
		{Key: "4", Title: "Vendas"},
		{Key: "5", Title: "Agenda"},
		{Key: "6", Title: "Reservas"},
		{Key: "7", Title: "SAC"},
		{Key: "8", Title: "Lojas"},
		{Key: "9", Title: loginTabTitle},
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
		body = m.vendasView.View()
	case 4:
		body = m.agendaView.View()
	case 5:
		body = m.reservasView.View()
	case 6:
		body = m.contatosView.View()
	case 7:
		body = m.unidadesView.View()
	case 8:
		body = m.loginView.View()
	}

	shortcuts := []components.Shortcut{
		{Key: "1-9", Desc: "Mudar Aba"},
		{Key: "Tab", Desc: "Próxima Aba"},
		{Key: "r", Desc: "Recarregar"},
		{Key: "↑/↓", Desc: "Navegar"},
		{Key: "q", Desc: "Sair"},
	}

	footer := components.RenderFooter(shortcuts, m.statusMsg, m.isConnected, m.width)

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
