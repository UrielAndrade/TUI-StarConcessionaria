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
		Tipo:      "Novo", // Obrigatório pela API: "Novo" ou "Seminovo"
		LojaIDStr: "1",
		AnoStr:    "2024",
	}

	// Todos os campos em um único grupo com FormKeyMap permite
	// navegar livremente para frente (Enter, Tab, ↓) e para trás (Shift+Tab, Backtab, ↑)
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
				Title("Tipo / Condição do Veículo (Novo ou Seminovo)").
				Options(
					huh.NewOption("Novo (Zero KM)", "Novo"),
					huh.NewOption("Seminovo (Usado)", "Seminovo"),
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
				Placeholder("Ex: 15000 (ou 0 para Novo)").
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
	).WithTheme(huh.ThemeCharm()).WithKeyMap(FormKeyMap())

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

// === FORMULÁRIO DE EDIÇÃO DE VEÍCULO ===

type EditarVeiculoResultMsg struct {
	VeiculoID int
	Payload   *model.SalvarVeiculoRequest
	Cancel    bool
}

type EditarVeiculoForm struct {
	Form        *huh.Form
	VeiculoID   int
	LojaIDStr   string
	Tipo        string
	Marca       string
	Modelo      string
	AnoStr      string
	PrecoStr    string
	KmStr       string
	Cor         string
	Descricao   string
}

func NewEditarVeiculoForm(v *model.VeiculoCatalogoResponse) *EditarVeiculoForm {
	tipoInicial := v.Tipo
	if tipoInicial != "Novo" && tipoInicial != "Seminovo" {
		tipoInicial = "Seminovo"
	}

	f := &EditarVeiculoForm{
		VeiculoID: v.ID,
		Tipo:      tipoInicial,
		LojaIDStr: strconv.Itoa(v.LojaID),
		Marca:     v.Marca,
		Modelo:    v.Modelo,
		AnoStr:    strconv.Itoa(v.Ano),
		PrecoStr:  fmt.Sprintf("%.2f", v.Preco),
		KmStr:     strconv.Itoa(v.Quilometragem),
		Cor:       v.Cor,
	}

	if f.LojaIDStr == "0" {
		f.LojaIDStr = "1"
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
				Title("Tipo / Condição do Veículo (Novo ou Seminovo)").
				Options(
					huh.NewOption("Novo (Zero KM)", "Novo"),
					huh.NewOption("Seminovo (Usado)", "Seminovo"),
				).
				Value(&f.Tipo),

			huh.NewInput().
				Title("Marca").
				Value(&f.Marca).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("a marca é obrigatória")
					}
					return nil
				}),

			huh.NewInput().
				Title("Modelo").
				Value(&f.Modelo).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("o modelo é obrigatório")
					}
					return nil
				}),

			huh.NewInput().
				Title("Ano de Fabricação").
				Value(&f.AnoStr).
				Validate(func(s string) error {
					v, err := strconv.Atoi(s)
					if err != nil || v < 1950 || v > 2030 {
						return fmt.Errorf("informe um ano válido")
					}
					return nil
				}),

			huh.NewInput().
				Title("Preço de Venda (R$)").
				Value(&f.PrecoStr).
				Validate(func(s string) error {
					v, err := strconv.ParseFloat(s, 64)
					if err != nil || v <= 0 {
						return fmt.Errorf("informe um preço positivo")
					}
					return nil
				}),

			huh.NewInput().
				Title("Quilometragem (km)").
				Value(&f.KmStr).
				Validate(func(s string) error {
					v, err := strconv.Atoi(s)
					if err != nil || v < 0 {
						return fmt.Errorf("informe uma quilometragem válida")
					}
					return nil
				}),

			huh.NewInput().
				Title("Cor").
				Value(&f.Cor),

			huh.NewInput().
				Title("Descrição / Opcionais").
				Value(&f.Descricao),
		),
	).WithTheme(huh.ThemeCharm()).WithKeyMap(FormKeyMap())

	return f
}

func (f *EditarVeiculoForm) Update(msg tea.Msg) (*EditarVeiculoForm, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.String() == "esc" {
		return f, func() tea.Msg {
			return EditarVeiculoResultMsg{Cancel: true}
		}
	}

	form, cmd := f.Form.Update(msg)
	if updated, ok := form.(*huh.Form); ok {
		f.Form = updated
	}

	if f.Form.State == huh.StateCompleted {
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
			return EditarVeiculoResultMsg{
				VeiculoID: f.VeiculoID,
				Payload:   payload,
			}
		}
	}

	return f, cmd
}

func (f *EditarVeiculoForm) View() string {
	return f.Form.View()
}
