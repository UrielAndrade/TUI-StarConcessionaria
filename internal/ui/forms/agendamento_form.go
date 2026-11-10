package forms

import (
	"fmt"
	"strconv"
	"time"

	"github.com/UrielAndrade/TUI-StarConcessionaria/internal/model"
	"github.com/charmbracelet/huh"
	tea "github.com/charmbracelet/bubbletea"
)

// === AGENDAR VISITA ===

type AgendarVisitaResultMsg struct {
	Payload *model.AgendarVisitaRequest
	Cancel  bool
}

type AgendarVisitaForm struct {
	Form            *huh.Form
	UnidadeIDStr    string
	NomeCliente     string
	EmailCliente    string
	TelefoneCliente string
	DataHoraStr     string
}

func NewAgendarVisitaForm() *AgendarVisitaForm {
	f := &AgendarVisitaForm{
		UnidadeIDStr: "1",
		DataHoraStr:  time.Now().Add(24 * time.Hour).Format("2006-01-02 10:00"),
	}

	f.Form = huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Unidade da Concessionária").
				Options(
					huh.NewOption("Matriz São Paulo (ID 1)", "1"),
					huh.NewOption("Filial Campinas (ID 2)", "2"),
				).
				Value(&f.UnidadeIDStr),

			huh.NewInput().
				Title("Nome do Cliente").
				Placeholder("Ex: Carlos Silva").
				Value(&f.NomeCliente).
				Validate(func(s string) error {
					if len(s) < 3 {
						return fmt.Errorf("o nome é obrigatório")
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
				Placeholder("(11) 98888-7777").
				Value(&f.TelefoneCliente).
				Validate(func(s string) error {
					if len(s) < 8 {
						return fmt.Errorf("o telefone é obrigatório")
					}
					return nil
				}),

			huh.NewInput().
				Title("Data e Horário (AAAA-MM-DD HH:MM)").
				Placeholder("2026-11-20 14:30").
				Value(&f.DataHoraStr).
				Validate(func(s string) error {
					_, err := time.Parse("2006-01-02 15:04", s)
					if err != nil {
						return fmt.Errorf("use o formato AAAA-MM-DD HH:MM (ex: 2026-11-20 14:30)")
					}
					return nil
				}),
		),
	).WithTheme(huh.ThemeCharm())

	return f
}

func (f *AgendarVisitaForm) Update(msg tea.Msg) (*AgendarVisitaForm, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.String() == "esc" {
		return f, func() tea.Msg {
			return AgendarVisitaResultMsg{Cancel: true}
		}
	}

	form, cmd := f.Form.Update(msg)
	if updated, ok := form.(*huh.Form); ok {
		f.Form = updated
	}

	if f.Form.State == huh.StateCompleted {
		uID, _ := strconv.Atoi(f.UnidadeIDStr)
		dt, _ := time.Parse("2006-01-02 15:04", f.DataHoraStr)

		payload := &model.AgendarVisitaRequest{
			UnidadeID:       uID,
			NomeCliente:     f.NomeCliente,
			EmailCliente:    f.EmailCliente,
			TelefoneCliente: f.TelefoneCliente,
			DataHoraInicio:  dt,
		}

		return f, func() tea.Msg {
			return AgendarVisitaResultMsg{Payload: payload}
		}
	}

	return f, cmd
}

func (f *AgendarVisitaForm) View() string {
	return f.Form.View()
}

// === AGENDAR TEST-DRIVE ===

type AgendarTestDriveResultMsg struct {
	Payload *model.AgendarTestDriveRequest
	Cancel  bool
}

type AgendarTestDriveForm struct {
	Form            *huh.Form
	VeiculoIDStr    string
	UnidadeIDStr    string
	NomeCliente     string
	EmailCliente    string
	TelefoneCliente string
	CNHCliente      string
	DataHoraStr     string
}

func NewAgendarTestDriveForm(veiculoIDPadrao int) *AgendarTestDriveForm {
	f := &AgendarTestDriveForm{
		VeiculoIDStr: strconv.Itoa(veiculoIDPadrao),
		UnidadeIDStr: "1",
		DataHoraStr:  time.Now().Add(24 * time.Hour).Format("2006-01-02 10:00"),
	}

	f.Form = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("ID do Veículo para Test-Drive").
				Placeholder("Ex: 1, 2, 3...").
				Value(&f.VeiculoIDStr).
				Validate(func(s string) error {
					v, err := strconv.Atoi(s)
					if err != nil || v <= 0 {
						return fmt.Errorf("informe o ID numérico do veículo")
					}
					return nil
				}),

			huh.NewSelect[string]().
				Title("Unidade da Loja").
				Options(
					huh.NewOption("Matriz São Paulo (ID 1)", "1"),
					huh.NewOption("Filial Campinas (ID 2)", "2"),
				).
				Value(&f.UnidadeIDStr),

			huh.NewInput().
				Title("Nome do Condutor").
				Placeholder("Ex: Roberto Antunes").
				Value(&f.NomeCliente).
				Validate(func(s string) error {
					if len(s) < 3 {
						return fmt.Errorf("o nome é obrigatório")
					}
					return nil
				}),

			huh.NewInput().
				Title("E-mail").
				Placeholder("condutor@email.com").
				Value(&f.EmailCliente),

			huh.NewInput().
				Title("Telefone").
				Placeholder("(11) 97777-6666").
				Value(&f.TelefoneCliente),

			huh.NewInput().
				Title("Número da CNH (Habilitação)").
				Placeholder("11 dígitos da CNH").
				Value(&f.CNHCliente).
				Validate(func(s string) error {
					if len(s) < 5 {
						return fmt.Errorf("informe a CNH válida do condutor")
					}
					return nil
				}),

			huh.NewInput().
				Title("Data e Horário (AAAA-MM-DD HH:MM)").
				Placeholder("2026-11-20 15:00").
				Value(&f.DataHoraStr).
				Validate(func(s string) error {
					_, err := time.Parse("2006-01-02 15:04", s)
					if err != nil {
						return fmt.Errorf("use o formato AAAA-MM-DD HH:MM")
					}
					return nil
				}),
		),
	).WithTheme(huh.ThemeCharm())

	return f
}

func (f *AgendarTestDriveForm) Update(msg tea.Msg) (*AgendarTestDriveForm, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.String() == "esc" {
		return f, func() tea.Msg {
			return AgendarTestDriveResultMsg{Cancel: true}
		}
	}

	form, cmd := f.Form.Update(msg)
	if updated, ok := form.(*huh.Form); ok {
		f.Form = updated
	}

	if f.Form.State == huh.StateCompleted {
		vID, _ := strconv.Atoi(f.VeiculoIDStr)
		uID, _ := strconv.Atoi(f.UnidadeIDStr)
		dt, _ := time.Parse("2006-01-02 15:04", f.DataHoraStr)

		payload := &model.AgendarTestDriveRequest{
			VeiculoID:       vID,
			UnidadeID:       uID,
			NomeCliente:     f.NomeCliente,
			EmailCliente:    f.EmailCliente,
			TelefoneCliente: f.TelefoneCliente,
			CNHCliente:      f.CNHCliente,
			DataHoraInicio:  dt,
		}

		return f, func() tea.Msg {
			return AgendarTestDriveResultMsg{Payload: payload}
		}
	}

	return f, cmd
}

func (f *AgendarTestDriveForm) View() string {
	return f.Form.View()
}

// === ATUALIZAR STATUS DE AGENDAMENTO ===

type StatusAgendamentoResultMsg struct {
	AgendamentoID int
	Status        string
	Observacoes   string
	Cancel        bool
}

type StatusAgendamentoForm struct {
	Form          *huh.Form
	AgendamentoID int
	Status        string
	Observacoes   string
}

func NewStatusAgendamentoForm(id int, statusAtual string) *StatusAgendamentoForm {
	f := &StatusAgendamentoForm{
		AgendamentoID: id,
		Status:        statusAtual,
	}

	f.Form = huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title(fmt.Sprintf("Alterar Status do Agendamento #%d", id)).
				Options(
					huh.NewOption("Confirmado (Agendamento aceito)", "Confirmado"),
					huh.NewOption("Concluído (Cliente compareceu)", "Concluido"),
					huh.NewOption("Cancelado (Desistência ou ausência)", "Cancelado"),
				).
				Value(&f.Status),

			huh.NewInput().
				Title("Observações / Parecer").
				Placeholder("Ex: Compareceu na concessionária, cliente aprovou o veículo.").
				Value(&f.Observacoes),
		),
	).WithTheme(huh.ThemeCharm())

	return f
}

func (f *StatusAgendamentoForm) Update(msg tea.Msg) (*StatusAgendamentoForm, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.String() == "esc" {
		return f, func() tea.Msg {
			return StatusAgendamentoResultMsg{Cancel: true}
		}
	}

	form, cmd := f.Form.Update(msg)
	if updated, ok := form.(*huh.Form); ok {
		f.Form = updated
	}

	if f.Form.State == huh.StateCompleted {
		return f, func() tea.Msg {
			return StatusAgendamentoResultMsg{
				AgendamentoID: f.AgendamentoID,
				Status:        f.Status,
				Observacoes:   f.Observacoes,
			}
		}
	}

	return f, cmd
}

func (f *StatusAgendamentoForm) View() string {
	return f.Form.View()
}
