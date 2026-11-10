package forms

import (
	"fmt"
	"strconv"

	"github.com/UrielAndrade/TUI-StarConcessionaria/internal/model"
	"github.com/charmbracelet/huh"
	tea "github.com/charmbracelet/bubbletea"
)

type NovaVendaResultMsg struct {
	Payload *model.RegistrarVendaRequest
	Cancel  bool
}

type NovaVendaForm struct {
	Form             *huh.Form
	VeiculoIDStr     string
	NomeCliente      string
	EmailCliente     string
	TelefoneCliente  string
	FormaPagamentoStr string
	ValorFinalStr    string
}

func NewNovaVendaForm(veiculoIDPadrao int, valorPadrao float64) *NovaVendaForm {
	f := &NovaVendaForm{
		VeiculoIDStr:      strconv.Itoa(veiculoIDPadrao),
		FormaPagamentoStr: "1",
	}
	if valorPadrao > 0 {
		f.ValorFinalStr = fmt.Sprintf("%.2f", valorPadrao)
	}

	f.Form = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("ID do Veículo Vendido").
				Placeholder("Ex: 1, 2, 3...").
				Value(&f.VeiculoIDStr).
				Validate(func(s string) error {
					v, err := strconv.Atoi(s)
					if err != nil || v <= 0 {
						return fmt.Errorf("informe um ID numérico de veículo válido")
					}
					return nil
				}),

			huh.NewInput().
				Title("Nome Completo do Cliente").
				Placeholder("Ex: Marcelo dos Santos").
				Value(&f.NomeCliente).
				Validate(func(s string) error {
					if len(s) < 3 {
						return fmt.Errorf("o nome do cliente deve ter ao menos 3 caracteres")
					}
					return nil
				}),

			huh.NewInput().
				Title("E-mail do Cliente").
				Placeholder("cliente@email.com").
				Value(&f.EmailCliente).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("o e-mail é obrigatório")
					}
					return nil
				}),

			huh.NewInput().
				Title("Telefone de Contato").
				Placeholder("(11) 98765-4321").
				Value(&f.TelefoneCliente).
				Validate(func(s string) error {
					if len(s) < 8 {
						return fmt.Errorf("informe um telefone de contato válido")
					}
					return nil
				}),

			huh.NewSelect[string]().
				Title("Forma de Pagamento").
				Options(
					huh.NewOption("1 - À vista (TED / Pix integral)", "1"),
					huh.NewOption("2 - Financiamento Bancário", "2"),
					huh.NewOption("3 - Cartão de Crédito (Parcelado)", "3"),
					huh.NewOption("4 - Pix Instantâneo", "4"),
					huh.NewOption("5 - Consórcio Contemplado", "5"),
				).
				Value(&f.FormaPagamentoStr),

			huh.NewInput().
				Title("Valor Final da Venda (R$)").
				Placeholder("Ex: 175000.00").
				Value(&f.ValorFinalStr).
				Validate(func(s string) error {
					v, err := strconv.ParseFloat(s, 64)
					if err != nil || v <= 0 {
						return fmt.Errorf("informe um valor monetário positivo válido")
					}
					return nil
				}),
		),
	).WithTheme(huh.ThemeCharm()).WithKeyMap(FormKeyMap())

	return f
}

func (f *NovaVendaForm) Update(msg tea.Msg) (*NovaVendaForm, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.String() == "esc" {
		return f, func() tea.Msg {
			return NovaVendaResultMsg{Cancel: true}
		}
	}

	form, cmd := f.Form.Update(msg)
	if updated, ok := form.(*huh.Form); ok {
		f.Form = updated
	}

	if f.Form.State == huh.StateCompleted {
		vID, _ := strconv.Atoi(f.VeiculoIDStr)
		fID, _ := strconv.Atoi(f.FormaPagamentoStr)
		valor, _ := strconv.ParseFloat(f.ValorFinalStr, 64)

		payload := &model.RegistrarVendaRequest{
			VeiculoID:        vID,
			NomeCliente:      f.NomeCliente,
			EmailCliente:     f.EmailCliente,
			TelefoneCliente:  f.TelefoneCliente,
			FormaPagamentoID: fID,
			ValorFinal:       valor,
		}

		return f, func() tea.Msg {
			return NovaVendaResultMsg{Payload: payload}
		}
	}

	return f, cmd
}

func (f *NovaVendaForm) View() string {
	return f.Form.View()
}
