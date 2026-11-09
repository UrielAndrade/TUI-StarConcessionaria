package forms

import (
	"fmt"
	"strconv"

	"github.com/UrielAndrade/TUI-StarConcessionaria/internal/model"
	"github.com/charmbracelet/huh"
	tea "github.com/charmbracelet/bubbletea"
)

type NovoVeiculoResultMsg struct {
	Payload *model.SalvarVeiculoRequest
	Cancel  bool
}

type NovoVeiculoForm struct {
	Form        *huh.Form
	LojaIDStr   string
	Tipo        string
	Marca       string
	Modelo      string
	AnoStr      string
	PrecoStr    string
	KmStr       string
	Cor         string
	Descricao   string
	Finalizado  bool
	Cancelado   bool
}

func NewNovoVeiculoForm() *NovoVeiculoForm {
	f := &NovoVeiculoForm{
		Tipo:      "SUV",
		LojaIDStr: "1",
		AnoStr:    "2024",
	}

	f.Form = huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Concessionária / Unidade").
				Options(
					huh.NewOption("Matriz São Paulo (ID 1)", "1"),
					huh.NewOption("Filial Campinas (ID 2)", "2"),
				).
				Value(&f.LojaIDStr),

			huh.NewSelect[string]().
				Title("Categoria / Tipo").
				Options(
					huh.NewOption("SUV", "SUV"),
					huh.NewOption("Sedan", "Sedan"),
					huh.NewOption("Hatch", "Hatch"),
					huh.NewOption("Picape", "Picape"),
					huh.NewOption("Esportivo", "Esportivo"),
				).
				Value(&f.Tipo),

			huh.NewInput().
				Title("Marca").
				Placeholder("Ex: Toyota, Honda, Ford...").
				Value(&f.Marca).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("a marca é obrigatória")
					}
					return nil
				}),

			huh.NewInput().
				Title("Modelo").
				Placeholder("Ex: Corolla Cross, Civic Touring...").
				Value(&f.Modelo).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("o modelo é obrigatório")
					}
					return nil
				}),
		),
		huh.NewGroup(
			huh.NewInput().
				Title("Ano de Fabricação").
				Placeholder("2024").
				Value(&f.AnoStr).
				Validate(func(s string) error {
					v, err := strconv.Atoi(s)
					if err != nil || v < 1950 || v > 2030 {
						return fmt.Errorf("informe um ano válido entre 1950 e 2030")
					}
					return nil
				}),

			huh.NewInput().
				Title("Preço de Venda (R$)").
				Placeholder("Ex: 185000.00").
				Value(&f.PrecoStr).
				Validate(func(s string) error {
					v, err := strconv.ParseFloat(s, 64)
					if err != nil || v <= 0 {
						return fmt.Errorf("informe um preço positivo válido")
					}
					return nil
				}),

			huh.NewInput().
				Title("Quilometragem (km)").
				Placeholder("Ex: 15000").
				Value(&f.KmStr).
				Validate(func(s string) error {
					v, err := strconv.Atoi(s)
					if err != nil || v < 0 {
						return fmt.Errorf("informe uma quilometragem válida (>= 0)")
					}
					return nil
				}),

			huh.NewInput().
				Title("Cor").
				Placeholder("Ex: Branco Perolizado, Preto...").
				Value(&f.Cor),

			huh.NewInput().
				Title("Descrição / Opcionais").
				Placeholder("Ex: Único dono, teto solar, revisado...").
				Value(&f.Descricao),
		),
	).WithTheme(huh.ThemeCharm())

	return f
}

func (f *NovoVeiculoForm) Update(msg tea.Msg) (*NovoVeiculoForm, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.String() == "esc" {
		f.Cancelado = true
		return f, func() tea.Msg {
			return NovoVeiculoResultMsg{Cancel: true}
		}
	}

	form, cmd := f.Form.Update(msg)
	if updated, ok := form.(*huh.Form); ok {
		f.Form = updated
	}

	if f.Form.State == huh.StateCompleted {
		f.Finalizado = true
		lojaID, _ := strconv.Atoi(f.LojaIDStr)
		ano, _ := strconv.Atoi(f.AnoStr)
		preco, _ := strconv.ParseFloat(f.PrecoStr, 64)
		km, _ := strconv.Atoi(f.KmStr)

		payload := &model.SalvarVeiculoRequest{
			LojaID:        lojaID,
			Tipo:          f.Tipo,
			Marca:         f.Marca,
			Modelo:        f.Modelo,
			Ano:           ano,
			Preco:         preco,
			Quilometragem: km,
			Cor:           f.Cor,
			Descricao:     f.Descricao,
		}

		return f, func() tea.Msg {
			return NovoVeiculoResultMsg{Payload: payload}
		}
	}

	return f, cmd
}

func (f *NovoVeiculoForm) View() string {
	return f.Form.View()
}
