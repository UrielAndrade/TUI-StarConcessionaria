package forms

import (
	"fmt"
	"strconv"

	"github.com/UrielAndrade/TUI-StarConcessionaria/internal/model"
	"github.com/charmbracelet/huh"
	tea "github.com/charmbracelet/bubbletea"
)

// === FORMULÁRIO DE CADASTRO DE NOVA PEÇA ===

type NovaPecaResultMsg struct {
	Payload *model.SalvarPecaRequest
	Cancel  bool
}

type NovaPecaForm struct {
	Form       *huh.Form
	Nome       string
	Descricao  string
	PrecoStr   string
	EstoqueStr string
}

func NewNovaPecaForm() *NovaPecaForm {
	f := &NovaPecaForm{
		EstoqueStr: "10",
	}

	f.Form = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Nome da Peça").
				Placeholder("Ex: Pastilha de Freio Cerâmica").
				Value(&f.Nome).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("o nome da peça é obrigatório")
					}
					return nil
				}),

			huh.NewInput().
				Title("Descrição / Aplicação").
				Placeholder("Ex: Jogo dianteiro para Corolla 2020+").
				Value(&f.Descricao),

			huh.NewInput().
				Title("Preço Unitário (R$)").
				Placeholder("Ex: 280.50").
				Value(&f.PrecoStr).
				Validate(func(s string) error {
					v, err := strconv.ParseFloat(s, 64)
					if err != nil || v <= 0 {
						return fmt.Errorf("informe um preço positivo")
					}
					return nil
				}),

			huh.NewInput().
				Title("Quantidade Inicial em Estoque").
				Placeholder("Ex: 15").
				Value(&f.EstoqueStr).
				Validate(func(s string) error {
					v, err := strconv.Atoi(s)
					if err != nil || v < 0 {
						return fmt.Errorf("informe uma quantidade válida (>= 0)")
					}
					return nil
				}),
		),
	).WithTheme(huh.ThemeCharm()).WithKeyMap(FormKeyMap())

	return f
}

func (f *NovaPecaForm) Update(msg tea.Msg) (*NovaPecaForm, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.String() == "esc" {
		return f, func() tea.Msg {
			return NovaPecaResultMsg{Cancel: true}
		}
	}

	form, cmd := f.Form.Update(msg)
	if updated, ok := form.(*huh.Form); ok {
		f.Form = updated
	}

	if f.Form.State == huh.StateCompleted {
		preco, _ := strconv.ParseFloat(f.PrecoStr, 64)
		qtd, _ := strconv.Atoi(f.EstoqueStr)

		payload := &model.SalvarPecaRequest{
			Nome:              f.Nome,
			Descricao:         f.Descricao,
			Preco:             preco,
			QuantidadeEstoque: qtd,
		}

		return f, func() tea.Msg {
			return NovaPecaResultMsg{Payload: payload}
		}
	}

	return f, cmd
}

func (f *NovaPecaForm) View() string {
	return f.Form.View()
}

// === FORMULÁRIO DE BAIXA / VENDA DIRETA DE PEÇA ===

type VenderPecaResultMsg struct {
	PecaID  int
	Payload *model.RegistrarVendaPecaRequest
	Cancel  bool
}

type VenderPecaForm struct {
	Form        *huh.Form
	PecaID      int
	PecaNome    string
	QtdMax      int
	QtdStr      string
	NomeCliente string
	Email       string
	Telefone    string
}

func NewVenderPecaForm(pecaID int, pecaNome string, estoqueAtual int) *VenderPecaForm {
	f := &VenderPecaForm{
		PecaID:   pecaID,
		PecaNome: pecaNome,
		QtdMax:   estoqueAtual,
		QtdStr:   "1",
	}

	f.Form = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title(fmt.Sprintf("Quantidade a Vender (Estoque Disponível: %d un.)", estoqueAtual)).
				Placeholder("1").
				Value(&f.QtdStr).
				Validate(func(s string) error {
					v, err := strconv.Atoi(s)
					if err != nil || v <= 0 {
						return fmt.Errorf("informe uma quantidade maior que zero")
					}
					if v > estoqueAtual {
						return fmt.Errorf("quantidade solicitada (%d) excede o estoque disponível (%d)", v, estoqueAtual)
					}
					return nil
				}),

			huh.NewInput().
				Title("Nome do Cliente").
				Placeholder("Ex: Carlos Eduardo").
				Value(&f.NomeCliente).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("o nome do cliente é obrigatório")
					}
					return nil
				}),

			huh.NewInput().
				Title("E-mail do Cliente").
				Placeholder("cliente@email.com").
				Value(&f.Email),

			huh.NewInput().
				Title("Telefone de Contato").
				Placeholder("(11) 98765-4321").
				Value(&f.Telefone),
		),
	).WithTheme(huh.ThemeCharm()).WithKeyMap(FormKeyMap())

	return f
}

func (f *VenderPecaForm) Update(msg tea.Msg) (*VenderPecaForm, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.String() == "esc" {
		return f, func() tea.Msg {
			return VenderPecaResultMsg{Cancel: true}
		}
	}

	form, cmd := f.Form.Update(msg)
	if updated, ok := form.(*huh.Form); ok {
		f.Form = updated
	}

	if f.Form.State == huh.StateCompleted {
		qtd, _ := strconv.Atoi(f.QtdStr)
		payload := &model.RegistrarVendaPecaRequest{
			Quantidade:      qtd,
			NomeCliente:     f.NomeCliente,
			EmailCliente:    f.Email,
			TelefoneCliente: f.Telefone,
		}

		return f, func() tea.Msg {
			return VenderPecaResultMsg{
				PecaID:  f.PecaID,
				Payload: payload,
			}
		}
	}

	return f, cmd
}

func (f *VenderPecaForm) View() string {
	return f.Form.View()
}
