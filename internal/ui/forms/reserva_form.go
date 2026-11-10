package forms

import (
	"fmt"
	"strconv"

	"github.com/UrielAndrade/TUI-StarConcessionaria/internal/model"
	"github.com/charmbracelet/huh"
	tea "github.com/charmbracelet/bubbletea"
)

// === SOLICITAR RESERVA ===

type NovaReservaResultMsg struct {
	Payload *model.SolicitarReservaRequest
	Cancel  bool
}

type NovaReservaForm struct {
	Form            *huh.Form
	VeiculoIDStr    string
	NomeCliente     string
	EmailCliente    string
	TelefoneCliente string
	Observacoes     string
}

func NewNovaReservaForm(veiculoIDPadrao int) *NovaReservaForm {
	f := &NovaReservaForm{
		VeiculoIDStr: strconv.Itoa(veiculoIDPadrao),
	}

	f.Form = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("ID do Veículo para Reserva").
				Placeholder("Ex: 1, 2, 3...").
				Value(&f.VeiculoIDStr).
				Validate(func(s string) error {
					v, err := strconv.Atoi(s)
					if err != nil || v <= 0 {
						return fmt.Errorf("informe o ID numérico do veículo")
					}
					return nil
				}),

			huh.NewInput().
				Title("Nome do Interessado").
				Placeholder("Ex: Patrícia Mendes").
				Value(&f.NomeCliente).
				Validate(func(s string) error {
					if len(s) < 3 {
						return fmt.Errorf("o nome é obrigatório")
					}
					return nil
				}),

			huh.NewInput().
				Title("E-mail").
				Placeholder("patricia@email.com").
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
						return fmt.Errorf("o telefone é obrigatório")
					}
					return nil
				}),

			huh.NewInput().
				Title("Observações / Condições").
				Placeholder("Ex: Sinal de R$ 5.000 pago via Pix, reserva por 48h.").
				Value(&f.Observacoes),
		),
	).WithTheme(huh.ThemeCharm()).WithKeyMap(FormKeyMap())

	return f
}

func (f *NovaReservaForm) Update(msg tea.Msg) (*NovaReservaForm, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.String() == "esc" {
		return f, func() tea.Msg {
			return NovaReservaResultMsg{Cancel: true}
		}
	}

	form, cmd := f.Form.Update(msg)
	if updated, ok := form.(*huh.Form); ok {
		f.Form = updated
	}

	if f.Form.State == huh.StateCompleted {
		vID, _ := strconv.Atoi(f.VeiculoIDStr)
		payload := &model.SolicitarReservaRequest{
			VeiculoID:       vID,
			NomeCliente:     f.NomeCliente,
			EmailCliente:    f.EmailCliente,
			TelefoneCliente: f.TelefoneCliente,
			Observacoes:     f.Observacoes,
		}

		return f, func() tea.Msg {
			return NovaReservaResultMsg{Payload: payload}
		}
	}

	return f, cmd
}

func (f *NovaReservaForm) View() string {
	return f.Form.View()
}

// === ATUALIZAR STATUS DE RESERVA ===

type StatusReservaResultMsg struct {
	ReservaID  int
	NovoStatus string
	Motivo     string
	Cancel     bool
}

type StatusReservaForm struct {
	Form       *huh.Form
	ReservaID  int
	NovoStatus string
	Motivo     string
}

func NewStatusReservaForm(id int, statusAtual string) *StatusReservaForm {
	f := &StatusReservaForm{
		ReservaID:  id,
		NovoStatus: statusAtual,
	}

	f.Form = huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title(fmt.Sprintf("Alterar Status da Reserva #%d", id)).
				Options(
					huh.NewOption("Ativa (Bloqueio mantido)", "Ativa"),
					huh.NewOption("Concluída (Venda efetivada)", "Concluida"),
					huh.NewOption("Cancelada (Prazo expirado / desistência)", "Cancelada"),
				).
				Value(&f.NovoStatus),

			huh.NewInput().
				Title("Motivo / Justificativa").
				Placeholder("Ex: Pagamento integral concluído ou desistência do lead.").
				Value(&f.Motivo),
		),
	).WithTheme(huh.ThemeCharm()).WithKeyMap(FormKeyMap())

	return f
}

func (f *StatusReservaForm) Update(msg tea.Msg) (*StatusReservaForm, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.String() == "esc" {
		return f, func() tea.Msg {
			return StatusReservaResultMsg{Cancel: true}
		}
	}

	form, cmd := f.Form.Update(msg)
	if updated, ok := form.(*huh.Form); ok {
		f.Form = updated
	}

	if f.Form.State == huh.StateCompleted {
		return f, func() tea.Msg {
			return StatusReservaResultMsg{
				ReservaID:  f.ReservaID,
				NovoStatus: f.NovoStatus,
				Motivo:     f.Motivo,
			}
		}
	}

	return f, cmd
}

func (f *StatusReservaForm) View() string {
	return f.Form.View()
}
