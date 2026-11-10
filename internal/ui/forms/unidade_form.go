package forms

import (
	"fmt"
	"strconv"

	"github.com/UrielAndrade/TUI-StarConcessionaria/internal/model"
	"github.com/charmbracelet/huh"
	tea "github.com/charmbracelet/bubbletea"
)

type NovaUnidadeResultMsg struct {
	Payload *model.SalvarUnidadeRequest
	Cancel  bool
}

type NovaUnidadeForm struct {
	Form        *huh.Form
	Nome        string
	Telefone    string
	CEP         string
	Endereco    string
	Numero      string
	Complemento string
	Bairro      string
	Cidade      string
	Estado      string
	LatStr      string
	LongStr     string
}

func NewNovaUnidadeForm() *NovaUnidadeForm {
	f := &NovaUnidadeForm{
		Estado:  "SP",
		LatStr:  "-23.5505",
		LongStr: "-46.6333",
	}

	f.Form = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Nome da Unidade / Filial").
				Placeholder("Ex: Star Concessionária - Filial Santos").
				Value(&f.Nome).
				Validate(func(s string) error {
					if len(s) < 3 {
						return fmt.Errorf("o nome da unidade é obrigatório")
					}
					return nil
				}),

			huh.NewInput().
				Title("Telefone").
				Placeholder("(13) 3200-1000").
				Value(&f.Telefone).
				Validate(func(s string) error {
					if len(s) < 8 {
						return fmt.Errorf("o telefone é obrigatório")
					}
					return nil
				}),

			huh.NewInput().
				Title("CEP (apenas números)").
				Placeholder("11010000").
				Value(&f.CEP).
				Validate(func(s string) error {
					if len(s) < 8 {
						return fmt.Errorf("informe o CEP com 8 dígitos")
					}
					return nil
				}),

			huh.NewInput().
				Title("Logradouro / Endereço").
				Placeholder("Av. Ana Costa").
				Value(&f.Endereco).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("o endereço é obrigatório")
					}
					return nil
				}),

			huh.NewInput().
				Title("Número").
				Placeholder("500").
				Value(&f.Numero),

			huh.NewInput().
				Title("Bairro").
				Placeholder("Gonzaga").
				Value(&f.Bairro),

			huh.NewInput().
				Title("Cidade").
				Placeholder("Santos").
				Value(&f.Cidade).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("a cidade é obrigatória")
					}
					return nil
				}),

			huh.NewInput().
				Title("UF / Estado (Ex: SP, RJ)").
				Placeholder("SP").
				Value(&f.Estado).
				Validate(func(s string) error {
					if len(s) != 2 {
						return fmt.Errorf("informe a sigla do estado com 2 letras")
					}
					return nil
				}),

			huh.NewInput().
				Title("Latitude").
				Placeholder("-23.9608").
				Value(&f.LatStr),

			huh.NewInput().
				Title("Longitude").
				Placeholder("-46.3336").
				Value(&f.LongStr),
		),
	).WithTheme(huh.ThemeCharm()).WithKeyMap(FormKeyMap())

	return f
}

func (f *NovaUnidadeForm) Update(msg tea.Msg) (*NovaUnidadeForm, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.String() == "esc" {
		return f, func() tea.Msg {
			return NovaUnidadeResultMsg{Cancel: true}
		}
	}

	form, cmd := f.Form.Update(msg)
	if updated, ok := form.(*huh.Form); ok {
		f.Form = updated
	}

	if f.Form.State == huh.StateCompleted {
		lat, _ := strconv.ParseFloat(f.LatStr, 64)
		long, _ := strconv.ParseFloat(f.LongStr, 64)

		payload := &model.SalvarUnidadeRequest{
			Nome:        f.Nome,
			Telefone:    f.Telefone,
			CEP:         f.CEP,
			Endereco:    f.Endereco,
			Numero:      f.Numero,
			Complemento: f.Complemento,
			Bairro:      f.Bairro,
			Cidade:      f.Cidade,
			Estado:      f.Estado,
			Latitude:    lat,
			Longitude:   long,
		}

		return f, func() tea.Msg {
			return NovaUnidadeResultMsg{Payload: payload}
		}
	}

	return f, cmd
}

func (f *NovaUnidadeForm) View() string {
	return f.Form.View()
}
