package forms

import (
	"fmt"
	"strconv"

	"github.com/UrielAndrade/TUI-StarConcessionaria/internal/model"
	"github.com/charmbracelet/huh"
	tea "github.com/charmbracelet/bubbletea"
)

type NovaPropostaResultMsg struct {
	Payload *model.SolicitarPropostaRequest
	Cancel  bool
}

type StatusPropostaResultMsg struct {
	PropostaID    int
	NotasResposta string
	Cancel        bool
}

type NovaPropostaForm struct {
	Form           *huh.Form
	VeiculoIDStr   string
	UnidadeIDStr   string
	NomeCliente    string
	EmailCliente   string
	Telefone       string
	FormaPagamento string
	PossuiTrocaStr string
	DadosTroca     string
	ValorEntrada   string
}

func NewNovaPropostaForm(veiculoID int, unidadeID int) *NovaPropostaForm {
	f := &NovaPropostaForm{
		VeiculoIDStr:   strconv.Itoa(veiculoID),
		UnidadeIDStr:   strconv.Itoa(unidadeID),
		FormaPagamento: "Financiamento",
		PossuiTrocaStr: "false",
		ValorEntrada:   "10000.00",
	}

	if f.UnidadeIDStr == "0" {
		f.UnidadeIDStr = "1"
	}
	if f.VeiculoIDStr == "0" {
		f.VeiculoIDStr = "1"
	}

	f.Form = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("ID do Veículo Desejado").
				Placeholder("1").
				Value(&f.VeiculoIDStr).
				Validate(func(s string) error {
					v, err := strconv.Atoi(s)
					if err != nil || v < 1 {
						return fmt.Errorf("informe um ID de veículo válido")
					}
					return nil
				}),

			huh.NewSelect[string]().
				Title("Concessionária").
				Options(
					huh.NewOption("Matriz São Paulo (ID 1)", "1"),
					huh.NewOption("Filial Campinas (ID 2)", "2"),
				).
				Value(&f.UnidadeIDStr),

			huh.NewInput().
				Title("Nome do Comprador / Proponente").
				Placeholder("Ex: Carlos Drumond").
				Value(&f.NomeCliente).
				Validate(func(s string) error {
					if len(s) < 3 {
						return fmt.Errorf("o nome do cliente é obrigatório")
					}
					return nil
				}),

			huh.NewInput().
				Title("E-mail de Contato").
				Placeholder("carlos@email.com").
				Value(&f.EmailCliente).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("o e-mail é obrigatório")
					}
					return nil
				}),

			huh.NewInput().
				Title("Telefone / WhatsApp").
				Placeholder("(11) 97777-8888").
				Value(&f.Telefone),

			huh.NewSelect[string]().
				Title("Condição de Pagamento Pretendida").
				Options(
					huh.NewOption("Financiamento Bancário", "Financiamento"),
					huh.NewOption("À Vista / Transferência", "À Vista"),
					huh.NewOption("Consórcio Contemplado", "Consórcio"),
				).
				Value(&f.FormaPagamento),

			huh.NewInput().
				Title("Valor de Entrada Sugerido (R$)").
				Placeholder("15000.00").
				Value(&f.ValorEntrada).
				Validate(func(s string) error {
					v, err := strconv.ParseFloat(s, 64)
					if err != nil || v < 0 {
						return fmt.Errorf("informe um valor numérico válido")
					}
					return nil
				}),

			huh.NewSelect[string]().
				Title("Oferece Veículo Usado na Troca?").
				Options(
					huh.NewOption("Não (Sem veículo na troca)", "false"),
					huh.NewOption("Sim (Avaliar veículo usado)", "true"),
				).
				Value(&f.PossuiTrocaStr),

			huh.NewInput().
				Title("Detalhes do Veículo da Troca (Marca, Modelo, Ano)").
				Placeholder("Ex: VW Gol 1.0 2018 70.000km").
				Value(&f.DadosTroca),
		),
	).WithTheme(huh.ThemeCharm()).WithKeyMap(FormKeyMap())

	return f
}

func (f *NovaPropostaForm) Update(msg tea.Msg) (*NovaPropostaForm, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.String() == "esc" {
		return f, func() tea.Msg {
			return NovaPropostaResultMsg{Cancel: true}
		}
	}

	form, cmd := f.Form.Update(msg)
	if updated, ok := form.(*huh.Form); ok {
		f.Form = updated
	}

	if f.Form.State == huh.StateCompleted {
		vID, _ := strconv.Atoi(f.VeiculoIDStr)
		uID, _ := strconv.Atoi(f.UnidadeIDStr)
		entrada, _ := strconv.ParseFloat(f.ValorEntrada, 64)
		possuiTroca := f.PossuiTrocaStr == "true"

		payload := &model.SolicitarPropostaRequest{
			VeiculoID:       vID,
			UnidadeID:       uID,
			NomeCliente:     f.NomeCliente,
			EmailCliente:    f.EmailCliente,
			TelefoneCliente: f.Telefone,
			FormaPagamento:  f.FormaPagamento,
			PossuiTroca:     possuiTroca,
			DadosTroca:      f.DadosTroca,
			ValorEntrada:    entrada,
		}

		return f, func() tea.Msg {
			return NovaPropostaResultMsg{Payload: payload}
		}
	}

	return f, cmd
}

func (f *NovaPropostaForm) View() string {
	return f.Form.View()
}

// === FORMS: PARECER / STATUS DA PROPOSTA ===

type StatusPropostaForm struct {
	Form          *huh.Form
	PropostaID    int
	Protocolo     string
	NotasResposta string
}

func NewStatusPropostaForm(id int, protocolo string) *StatusPropostaForm {
	f := &StatusPropostaForm{
		PropostaID: id,
		Protocolo:  protocolo,
	}

	f.Form = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title(fmt.Sprintf("Parecer Comercial para a Proposta #%d (%s)", id, protocolo)).
				Placeholder("Ex: Aprovado com 30%% de entrada e taxa de 1.19%% a.m.").
				Value(&f.NotasResposta).
				Validate(func(s string) error {
					if len(s) < 3 {
						return fmt.Errorf("o parecer da proposta é obrigatório")
					}
					return nil
				}),
		),
	).WithTheme(huh.ThemeCharm()).WithKeyMap(FormKeyMap())

	return f
}

func (f *StatusPropostaForm) Update(msg tea.Msg) (*StatusPropostaForm, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.String() == "esc" {
		return f, func() tea.Msg {
			return StatusPropostaResultMsg{Cancel: true}
		}
	}

	form, cmd := f.Form.Update(msg)
	if updated, ok := form.(*huh.Form); ok {
		f.Form = updated
	}

	if f.Form.State == huh.StateCompleted {
		return f, func() tea.Msg {
			return StatusPropostaResultMsg{
				PropostaID:    f.PropostaID,
				NotasResposta: f.NotasResposta,
			}
		}
	}

	return f, cmd
}

func (f *StatusPropostaForm) View() string {
	return f.Form.View()
}
