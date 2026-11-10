package forms

import (
	"fmt"
	"strconv"

	"github.com/UrielAndrade/TUI-StarConcessionaria/internal/model"
	"github.com/charmbracelet/huh"
	tea "github.com/charmbracelet/bubbletea"
)

// === NOVO CHAMADO SAC ===

type NovoContatoResultMsg struct {
	Payload *model.NovoContatoRequest
	Cancel  bool
}

type NovoContatoForm struct {
	Form         *huh.Form
	Assunto      string
	Nome         string
	Email        string
	Telefone     string
	Mensagem     string
	UnidadeIDStr string
}

func NewNovoContatoForm() *NovoContatoForm {
	f := &NovoContatoForm{
		UnidadeIDStr: "1",
	}

	f.Form = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Assunto da Mensagem").
				Placeholder("Ex: Dúvida sobre garantia, proposta de compra...").
				Value(&f.Assunto).
				Validate(func(s string) error {
					if len(s) < 3 {
						return fmt.Errorf("o assunto é obrigatório")
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
				Title("Nome do Solicitante").
				Placeholder("Ex: Fernando Torres").
				Value(&f.Nome).
				Validate(func(s string) error {
					if len(s) < 3 {
						return fmt.Errorf("o nome é obrigatório")
					}
					return nil
				}),

			huh.NewInput().
				Title("E-mail para Retorno").
				Placeholder("fernando@email.com").
				Value(&f.Email).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("o e-mail é obrigatório")
					}
					return nil
				}),

			huh.NewInput().
				Title("Telefone").
				Placeholder("(11) 98765-4321").
				Value(&f.Telefone),

			huh.NewInput().
				Title("Mensagem / Descrição").
				Placeholder("Escreva aqui o conteúdo do atendimento...").
				Value(&f.Mensagem).
				Validate(func(s string) error {
					if len(s) < 5 {
						return fmt.Errorf("a mensagem deve conter ao menos 5 caracteres")
					}
					return nil
				}),
		),
	).WithTheme(huh.ThemeCharm())

	return f
}

func (f *NovoContatoForm) Update(msg tea.Msg) (*NovoContatoForm, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.String() == "esc" {
		return f, func() tea.Msg {
			return NovoContatoResultMsg{Cancel: true}
		}
	}

	form, cmd := f.Form.Update(msg)
	if updated, ok := form.(*huh.Form); ok {
		f.Form = updated
	}

	if f.Form.State == huh.StateCompleted {
		uID, _ := strconv.Atoi(f.UnidadeIDStr)
		payload := &model.NovoContatoRequest{
			Assunto:   f.Assunto,
			Nome:      f.Nome,
			Email:     f.Email,
			Telefone:  f.Telefone,
			Mensagem:  f.Mensagem,
			UnidadeID: &uID,
		}

		return f, func() tea.Msg {
			return NovoContatoResultMsg{Payload: payload}
		}
	}

	return f, cmd
}

func (f *NovoContatoForm) View() string {
	return f.Form.View()
}

// === RESPONDER CHAMADO SAC ===

type ResponderContatoResultMsg struct {
	ContatoID int
	Notas     string
	Cancel    bool
}

type ResponderContatoForm struct {
	Form      *huh.Form
	ContatoID int
	Protocolo string
	Notas     string
}

func NewResponderContatoForm(id int, protocolo string) *ResponderContatoForm {
	f := &ResponderContatoForm{
		ContatoID: id,
		Protocolo: protocolo,
	}

	f.Form = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title(fmt.Sprintf("Registrar Resposta para Chamado #%d (%s)", id, protocolo)).
				Placeholder("Ex: Entrado em contato via WhatsApp e proposta encaminhada.").
				Value(&f.Notas).
				Validate(func(s string) error {
					if len(s) < 3 {
						return fmt.Errorf("a nota de resposta é obrigatória")
					}
					return nil
				}),
		),
	).WithTheme(huh.ThemeCharm())

	return f
}

func (f *ResponderContatoForm) Update(msg tea.Msg) (*ResponderContatoForm, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.String() == "esc" {
		return f, func() tea.Msg {
			return ResponderContatoResultMsg{Cancel: true}
		}
	}

	form, cmd := f.Form.Update(msg)
	if updated, ok := form.(*huh.Form); ok {
		f.Form = updated
	}

	if f.Form.State == huh.StateCompleted {
		return f, func() tea.Msg {
			return ResponderContatoResultMsg{
				ContatoID: f.ContatoID,
				Notas:     f.Notas,
			}
		}
	}

	return f, cmd
}

func (f *ResponderContatoForm) View() string {
	return f.Form.View()
}
