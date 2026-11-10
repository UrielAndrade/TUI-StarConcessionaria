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
	"github.com/UrielAndrade/TUI-StarConcessionaria/internal/ui/views/propostas"
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

type PropostasLoadedMsg struct {
	Data *model.PagedResult[model.PropostaResponse]
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

// Mensagens de retorno das mutações
type OperacaoResultMsg struct {
	Modulo   string
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

	// 10 Telas Completas
	dashboardView dashboard.Model
	veiculosView  veiculos.Model
	pecasView     pecas.Model
	vendasView    vendas.Model
	agendaView    agendamentos.Model
	reservasView  reservas.Model
	propostasView propostas.Model
	contatosView  contatos.Model
	unidadesView  unidades.Model
	loginView     login.Model
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
		propostasView: propostas.New(),
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
			m.loadPropostasCmd(),
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
			ItensPorPagina: 15,
		})
		return VeiculosLoadedMsg{Data: data, Err: err}
	}
}

func (m AppModel) loadPecasCmd() tea.Cmd {
	return func() tea.Msg {
		data, err := m.client.ListPecas(context.Background(), model.FiltroPecas{
			Pagina:         1,
			ItensPorPagina: 15,
		})
		return PecasLoadedMsg{Data: data, Err: err}
	}
}

func (m AppModel) loadVendasCmd() tea.Cmd {
	return func() tea.Msg {
		data, err := m.client.ListVendas(context.Background(), model.FiltroVendas{
			Pagina:         1,
			ItensPorPagina: 15,
		})
		return VendasLoadedMsg{Data: data, Err: err}
	}
}

func (m AppModel) loadAgendamentosCmd() tea.Cmd {
	return func() tea.Msg {
		data, err := m.client.ListAgendamentos(context.Background(), model.FiltroAgendamentos{
			Pagina:         1,
			ItensPorPagina: 15,
		})
		return AgendamentosLoadedMsg{Data: data, Err: err}
	}
}

func (m AppModel) loadReservasCmd() tea.Cmd {
	return func() tea.Msg {
		data, err := m.client.ListReservas(context.Background(), model.FiltroReservas{
			Pagina:         1,
			ItensPorPagina: 15,
		})
		return ReservasLoadedMsg{Data: data, Err: err}
	}
}

func (m AppModel) loadPropostasCmd() tea.Cmd {
	return func() tea.Msg {
		data, err := m.client.ListPropostas(context.Background(), model.FiltroPropostas{
			Pagina:         1,
			ItensPorPagina: 15,
		})
		return PropostasLoadedMsg{Data: data, Err: err}
	}
}

func (m AppModel) loadContatosCmd() tea.Cmd {
	return func() tea.Msg {
		data, err := m.client.ListContatos(context.Background(), model.FiltroContatos{
			Pagina:         1,
			ItensPorPagina: 15,
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
		m.propostasView.Width = msg.Width
		m.propostasView.Height = msg.Height
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

		// Se qualquer modal interativo estiver aberto, repassa a tecla EXCLUSIVAMENTE para o formulário
		if m.activeTab == 1 && m.veiculosView.Estado != veiculos.EstadoLista {
			var cmd tea.Cmd
			m.veiculosView, cmd = m.veiculosView.Update(msg)
			return m, cmd
		}

		if m.activeTab == 2 && m.pecasView.Estado != pecas.EstadoLista {
			var cmd tea.Cmd
			m.pecasView, cmd = m.pecasView.Update(msg)
			return m, cmd
		}

		if m.activeTab == 3 && m.vendasView.Estado != vendas.EstadoLista {
			var cmd tea.Cmd
			m.vendasView, cmd = m.vendasView.Update(msg)
			return m, cmd
		}

		if m.activeTab == 4 && m.agendaView.Estado != agendamentos.EstadoLista {
			var cmd tea.Cmd
			m.agendaView, cmd = m.agendaView.Update(msg)
			return m, cmd
		}

		if m.activeTab == 5 && m.reservasView.Estado != reservas.EstadoLista {
			var cmd tea.Cmd
			m.reservasView, cmd = m.reservasView.Update(msg)
			return m, cmd
		}

		if m.activeTab == 6 && m.propostasView.Estado != propostas.EstadoLista {
			var cmd tea.Cmd
			m.propostasView, cmd = m.propostasView.Update(msg)
			return m, cmd
		}

		if m.activeTab == 7 && m.contatosView.Estado != contatos.EstadoLista {
			var cmd tea.Cmd
			m.contatosView, cmd = m.contatosView.Update(msg)
			return m, cmd
		}

		if m.activeTab == 8 && m.unidadesView.Estado != unidades.EstadoLista {
			var cmd tea.Cmd
			m.unidadesView, cmd = m.unidadesView.Update(msg)
			return m, cmd
		}

		// Se estiver na tela de login e não estiver autenticado
		if m.activeTab == 9 && !m.loginView.Logado {
			if msg.String() == "esc" {
				m.activeTab = 0 // Volta para o Dashboard
				return m, nil
			}
			var cmd tea.Cmd
			m.loginView, cmd = m.loginView.Update(msg)
			return m, cmd
		}

		// Navegação global entre as 10 telas:
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
			if len(m.propostasView.Propostas) == 0 {
				m.propostasView.Carregando = true
				cmds = append(cmds, m.loadPropostasCmd())
			}
		case "8":
			m.activeTab = 7
			if len(m.contatosView.Contatos) == 0 {
				m.contatosView.Carregando = true
				cmds = append(cmds, m.loadContatosCmd())
			}
		case "9":
			m.activeTab = 8
			if len(m.unidadesView.Unidades) == 0 {
				m.unidadesView.Carregando = true
				cmds = append(cmds, m.loadUnidadesCmd())
			}
		case "0":
			m.activeTab = 9
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
				m.propostasView.Carregando = true
				cmds = append(cmds, m.loadPropostasCmd())
			case 7:
				m.contatosView.Carregando = true
				cmds = append(cmds, m.loadContatosCmd())
			case 8:
				m.unidadesView.Carregando = true
				cmds = append(cmds, m.loadUnidadesCmd())
			}
		case "tab":
			m.activeTab = (m.activeTab + 1) % 10
		case "shift+tab":
			m.activeTab = (m.activeTab + 9) % 10
		case "esc":
			if m.activeTab != 0 {
				m.activeTab = 0
			}
		}

		// Se estiver na aba de perfil, permite logout via 'L'
		if m.activeTab == 9 && m.loginView.Logado {
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

	case PropostasLoadedMsg:
		m.propostasView.Carregando = false
		if msg.Err != nil {
			m.propostasView.Erro = msg.Err.Error()
		} else {
			m.propostasView.SetDados(msg.Data)
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

	// === PROCESSAMENTO DOS FORMULÁRIOS / CADASTROS ===

	case forms.NovoVeiculoResultMsg:
		m.veiculosView.Estado = veiculos.EstadoLista
		m.veiculosView.NovoForm = nil
		if !msg.Cancel && msg.Payload != nil {
			m.statusMsg = "Cadastrando veículo no estoque..."
			cmds = append(cmds, func() tea.Msg {
				_, err := m.client.CriarVeiculo(context.Background(), *msg.Payload)
				return OperacaoResultMsg{
					Modulo:   "veiculos",
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
				return OperacaoResultMsg{
					Modulo:   "veiculos",
					Mensagem: "Situação do veículo alterada para " + msg.NovaSituacao + "!",
					Err:      err,
				}
			})
		}

	case forms.AdicionarFotoResultMsg:
		m.veiculosView.Estado = veiculos.EstadoLista
		m.veiculosView.FotoForm = nil
		if !msg.Cancel && msg.Payload != nil {
			m.statusMsg = "Adicionando foto à galeria do veículo..."
			cmds = append(cmds, func() tea.Msg {
				_, err := m.client.AdicionarFotoVeiculo(context.Background(), msg.VeiculoID, *msg.Payload)
				return OperacaoResultMsg{
					Modulo:   "veiculos",
					Mensagem: "Foto adicionada à galeria com sucesso!",
					Err:      err,
				}
			})
		}

	case forms.NovaPecaResultMsg:
		m.pecasView.Estado = pecas.EstadoLista
		m.pecasView.NovaForm = nil
		if !msg.Cancel && msg.Payload != nil {
			m.statusMsg = "Cadastrando peça no almoxarifado..."
			cmds = append(cmds, func() tea.Msg {
				_, err := m.client.CriarPeca(context.Background(), *msg.Payload)
				return OperacaoResultMsg{
					Modulo:   "pecas",
					Mensagem: "Peça cadastrada com sucesso!",
					Err:      err,
				}
			})
		}

	case forms.VenderPecaResultMsg:
		m.pecasView.Estado = pecas.EstadoLista
		m.pecasView.VenderForm = nil
		if !msg.Cancel && msg.Payload != nil {
			m.statusMsg = "Registrando baixa/venda de peça..."
			cmds = append(cmds, func() tea.Msg {
				err := m.client.VenderPeca(context.Background(), msg.PecaID, *msg.Payload)
				return OperacaoResultMsg{
					Modulo:   "pecas",
					Mensagem: "Venda de peça registrada com sucesso!",
					Err:      err,
				}
			})
		}

	case forms.NovaVendaResultMsg:
		m.vendasView.Estado = vendas.EstadoLista
		m.vendasView.NovaForm = nil
		if !msg.Cancel && msg.Payload != nil {
			m.statusMsg = "Registrando venda de veículo..."
			cmds = append(cmds, func() tea.Msg {
				err := m.client.RegistrarVenda(context.Background(), *msg.Payload)
				return OperacaoResultMsg{
					Modulo:   "vendas",
					Mensagem: "Venda de veículo registrada com sucesso!",
					Err:      err,
				}
			})
		}

	case forms.AgendarVisitaResultMsg:
		m.agendaView.Estado = agendamentos.EstadoLista
		m.agendaView.VisitaForm = nil
		if !msg.Cancel && msg.Payload != nil {
			m.statusMsg = "Agendando visita presencial à loja..."
			cmds = append(cmds, func() tea.Msg {
				err := m.client.AgendarVisita(context.Background(), *msg.Payload)
				return OperacaoResultMsg{
					Modulo:   "agenda",
					Mensagem: "Visita presencial agendada com sucesso!",
					Err:      err,
				}
			})
		}

	case forms.AgendarTestDriveResultMsg:
		m.agendaView.Estado = agendamentos.EstadoLista
		m.agendaView.TestDriveForm = nil
		if !msg.Cancel && msg.Payload != nil {
			m.statusMsg = "Agendando test-drive veicular..."
			cmds = append(cmds, func() tea.Msg {
				err := m.client.AgendarTestDrive(context.Background(), *msg.Payload)
				return OperacaoResultMsg{
					Modulo:   "agenda",
					Mensagem: "Test-drive agendado com sucesso!",
					Err:      err,
				}
			})
		}

	case forms.StatusAgendamentoResultMsg:
		m.agendaView.Estado = agendamentos.EstadoLista
		m.agendaView.StatusForm = nil
		if !msg.Cancel {
			m.statusMsg = "Atualizando status do agendamento..."
			cmds = append(cmds, func() tea.Msg {
				err := m.client.AtualizarStatusAgendamento(context.Background(), msg.AgendamentoID, msg.Status, msg.Observacoes)
				return OperacaoResultMsg{
					Modulo:   "agenda",
					Mensagem: "Status do agendamento atualizado para " + msg.Status + "!",
					Err:      err,
				}
			})
		}

	case forms.NovaReservaResultMsg:
		m.reservasView.Estado = reservas.EstadoLista
		m.reservasView.NovaForm = nil
		if !msg.Cancel && msg.Payload != nil {
			m.statusMsg = "Solicitando reserva de veículo..."
			cmds = append(cmds, func() tea.Msg {
				err := m.client.SolicitarReserva(context.Background(), *msg.Payload)
				return OperacaoResultMsg{
					Modulo:   "reservas",
					Mensagem: "Reserva de veículo solicitada com sucesso!",
					Err:      err,
				}
			})
		}

	case forms.StatusReservaResultMsg:
		m.reservasView.Estado = reservas.EstadoLista
		m.reservasView.StatusForm = nil
		if !msg.Cancel {
			m.statusMsg = "Atualizando status da reserva..."
			cmds = append(cmds, func() tea.Msg {
				err := m.client.AtualizarStatusReserva(context.Background(), msg.ReservaID, msg.NovoStatus, msg.Motivo)
				return OperacaoResultMsg{
					Modulo:   "reservas",
					Mensagem: "Status da reserva atualizado para " + msg.NovoStatus + "!",
					Err:      err,
				}
			})
		}

	case forms.NovaPropostaResultMsg:
		m.propostasView.Estado = propostas.EstadoLista
		m.propostasView.NovaForm = nil
		if !msg.Cancel && msg.Payload != nil {
			m.statusMsg = "Enviando proposta comercial..."
			cmds = append(cmds, func() tea.Msg {
				err := m.client.SolicitarProposta(context.Background(), *msg.Payload)
				return OperacaoResultMsg{
					Modulo:   "propostas",
					Mensagem: "Proposta comercial enviada com sucesso!",
					Err:      err,
				}
			})
		}

	case forms.StatusPropostaResultMsg:
		m.propostasView.Estado = propostas.EstadoLista
		m.propostasView.StatusForm = nil
		if !msg.Cancel {
			m.statusMsg = "Registrando parecer sobre a proposta..."
			cmds = append(cmds, func() tea.Msg {
				err := m.client.AtualizarStatusProposta(context.Background(), msg.PropostaID, msg.NotasResposta)
				return OperacaoResultMsg{
					Modulo:   "propostas",
					Mensagem: "Parecer da proposta registrado com sucesso!",
					Err:      err,
				}
			})
		}

	case forms.NovoContatoResultMsg:
		m.contatosView.Estado = contatos.EstadoLista
		m.contatosView.NovoForm = nil
		if !msg.Cancel && msg.Payload != nil {
			m.statusMsg = "Enviando chamado ao SAC..."
			cmds = append(cmds, func() tea.Msg {
				err := m.client.NovoContato(context.Background(), *msg.Payload)
				return OperacaoResultMsg{
					Modulo:   "contatos",
					Mensagem: "Chamado do SAC registrado com sucesso!",
					Err:      err,
				}
			})
		}

	case forms.ResponderContatoResultMsg:
		m.contatosView.Estado = contatos.EstadoLista
		m.contatosView.ResponderForm = nil
		if !msg.Cancel {
			m.statusMsg = "Registrando atendimento ao chamado..."
			cmds = append(cmds, func() tea.Msg {
				err := m.client.ResponderContato(context.Background(), msg.ContatoID, msg.Notas)
				return OperacaoResultMsg{
					Modulo:   "contatos",
					Mensagem: "Atendimento do chamado SAC registrado!",
					Err:      err,
				}
			})
		}

	case forms.NovaUnidadeResultMsg:
		m.unidadesView.Estado = unidades.EstadoLista
		m.unidadesView.NovaForm = nil
		if !msg.Cancel && msg.Payload != nil {
			m.statusMsg = "Cadastrando nova concessionária..."
			cmds = append(cmds, func() tea.Msg {
				_, err := m.client.CriarUnidade(context.Background(), *msg.Payload)
				return OperacaoResultMsg{
					Modulo:   "unidades",
					Mensagem: "Nova concessionária cadastrada com sucesso!",
					Err:      err,
				}
			})
		}

	// Feedback Unificado de Mutação
	case OperacaoResultMsg:
		if msg.Err != nil {
			m.statusMsg = "Falha: " + msg.Err.Error()
			switch msg.Modulo {
			case "veiculos":
				m.veiculosView.Erro = msg.Err.Error()
			case "pecas":
				m.pecasView.Erro = msg.Err.Error()
			case "vendas":
				m.vendasView.Erro = msg.Err.Error()
			case "agenda":
				m.agendaView.Erro = msg.Err.Error()
			case "reservas":
				m.reservasView.Erro = msg.Err.Error()
			case "propostas":
				m.propostasView.Erro = msg.Err.Error()
			case "contatos":
				m.contatosView.Erro = msg.Err.Error()
			case "unidades":
				m.unidadesView.Erro = msg.Err.Error()
			}
		} else {
			m.statusMsg = "✔ " + msg.Mensagem
			switch msg.Modulo {
			case "veiculos":
				m.veiculosView.Carregando = true
				cmds = append(cmds, m.loadVeiculosCmd(), m.loadDashboardCmd())
			case "pecas":
				m.pecasView.Carregando = true
				cmds = append(cmds, m.loadPecasCmd(), m.loadDashboardCmd())
			case "vendas":
				m.vendasView.Carregando = true
				cmds = append(cmds, m.loadVendasCmd(), m.loadDashboardCmd())
			case "agenda":
				m.agendaView.Carregando = true
				cmds = append(cmds, m.loadAgendamentosCmd())
			case "reservas":
				m.reservasView.Carregando = true
				cmds = append(cmds, m.loadReservasCmd())
			case "propostas":
				m.propostasView.Carregando = true
				cmds = append(cmds, m.loadPropostasCmd())
			case "contatos":
				m.contatosView.Carregando = true
				cmds = append(cmds, m.loadContatosCmd())
			case "unidades":
				m.unidadesView.Carregando = true
				cmds = append(cmds, m.loadUnidadesCmd())
			}
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
			m.setAutenticado(true)
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
				m.loadPropostasCmd(),
				m.loadContatosCmd(),
			)
		}

	case PerfilLoadedMsg:
		if msg.Err == nil && msg.Usuario != nil {
			m.loginView.Logado = true
			m.loginView.Usuario = msg.Usuario.Nome
			m.setAutenticado(true)
			m.statusMsg = "Usuário: " + msg.Usuario.Nome
			cmds = append(cmds,
				m.loadDashboardCmd(),
				m.loadVendasCmd(),
				m.loadAgendamentosCmd(),
				m.loadReservasCmd(),
				m.loadPropostasCmd(),
				m.loadContatosCmd(),
			)
		}

	case login.LogoutMsg:
		m.cfg.Token = ""
		_ = m.cfg.Save()
		m.client.SetToken("")
		m.statusMsg = "Pronto (Desconectado)"
		m.activeTab = 0
		m.setAutenticado(false)
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
		m.propostasView, cmd = m.propostasView.Update(msg)
		cmds = append(cmds, cmd)
	case 7:
		var cmd tea.Cmd
		m.contatosView, cmd = m.contatosView.Update(msg)
		cmds = append(cmds, cmd)
	case 8:
		var cmd tea.Cmd
		m.unidadesView, cmd = m.unidadesView.Update(msg)
		cmds = append(cmds, cmd)
	case 9:
		if !m.loginView.Logado {
			var cmd tea.Cmd
			m.loginView, cmd = m.loginView.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m *AppModel) setAutenticado(autenticado bool) {
	m.veiculosView.IsAutenticado = autenticado
	m.pecasView.IsAutenticado = autenticado
	m.vendasView.IsAutenticado = autenticado
	m.agendaView.IsAutenticado = autenticado
	m.reservasView.IsAutenticado = autenticado
	m.propostasView.IsAutenticado = autenticado
	m.contatosView.IsAutenticado = autenticado
	m.unidadesView.IsAutenticado = autenticado
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
		{Key: "7", Title: "Propostas"},
		{Key: "8", Title: "SAC"},
		{Key: "9", Title: "Lojas"},
		{Key: "0", Title: loginTabTitle},
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
		body = m.propostasView.View()
	case 7:
		body = m.contatosView.View()
	case 8:
		body = m.unidadesView.View()
	case 9:
		body = m.loginView.View()
	}

	shortcuts := []components.Shortcut{
		{Key: "1-0", Desc: "Mudar Aba"},
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
