package forms

import (
	"fmt"
	"strconv"

	"github.com/UrielAndrade/TUI-StarConcessionaria/internal/model"
	"github.com/charmbracelet/huh"
	tea "github.com/charmbracelet/bubbletea"
)

type AdicionarFotoResultMsg struct {
	VeiculoID int
	Payload   *model.AdicionarFotoRequest
	Cancel    bool
}

type AdicionarFotoForm struct {
	Form         *huh.Form
	VeiculoID    int
	VeiculoDesc  string
	URL          string
	PrincipalStr string
	OrdemStr     string
}

func NewAdicionarFotoForm(veiculoID int, veiculoDesc string) *AdicionarFotoForm {
	f := &AdicionarFotoForm{
		VeiculoID:    veiculoID,
		VeiculoDesc:  veiculoDesc,
		PrincipalStr: "false",
		OrdemStr:     "1",
	}

	f.Form = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title(fmt.Sprintf("Adicionar Foto para Veículo #%d: %s", veiculoID, veiculoDesc)).
				Placeholder("https://images.unsplash.com/...").
				Value(&f.URL).
				Validate(func(s string) error {
					if len(s) < 10 {
						return fmt.Errorf("informe uma URL válida para a imagem")
					}
					return nil
				}),

			huh.NewSelect[string]().
				Title("Definir como Imagem Principal (Capa)?").
				Options(
					huh.NewOption("Não (Foto secundária na galeria)", "false"),
					huh.NewOption("Sim (Definir como foto de capa)", "true"),
				).
				Value(&f.PrincipalStr),

			huh.NewInput().
				Title("Ordem de Exibição").
				Placeholder("1").
				Value(&f.OrdemStr).
				Validate(func(s string) error {
					v, err := strconv.Atoi(s)
					if err != nil || v < 1 {
						return fmt.Errorf("informe uma ordem numérica positiva (>= 1)")
					}
					return nil
				}),
		),
	).WithTheme(huh.ThemeCharm()).WithKeyMap(FormKeyMap())

	return f
}

func (f *AdicionarFotoForm) Update(msg tea.Msg) (*AdicionarFotoForm, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.String() == "esc" {
		return f, func() tea.Msg {
			return AdicionarFotoResultMsg{Cancel: true}
		}
	}

	form, cmd := f.Form.Update(msg)
	if updated, ok := form.(*huh.Form); ok {
		f.Form = updated
	}

	if f.Form.State == huh.StateCompleted {
		ordem, _ := strconv.Atoi(f.OrdemStr)
		principal := f.PrincipalStr == "true"

		payload := &model.AdicionarFotoRequest{
			URL:       f.URL,
			Principal: principal,
			Ordem:     ordem,
		}

		return f, func() tea.Msg {
			return AdicionarFotoResultMsg{
				VeiculoID: f.VeiculoID,
				Payload:   payload,
			}
		}
	}

	return f, cmd
}

func (f *AdicionarFotoForm) View() string {
	return f.Form.View()
}
