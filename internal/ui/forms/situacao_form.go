package forms

import (
	"fmt"

	"github.com/charmbracelet/huh"
	tea "github.com/charmbracelet/bubbletea"
)

type AlterarSituacaoResultMsg struct {
	VeiculoID     int
	NovaSituacao  string
	Justificativa string
	Cancel        bool
}

type AlterarSituacaoForm struct {
	Form          *huh.Form
	VeiculoID     int
	VeiculoDesc   string
	NovaSituacao  string
	Justificativa string
	Cancelado     bool
}

func NewAlterarSituacaoForm(veiculoID int, veiculoDesc, situacaoAtual string) *AlterarSituacaoForm {
	f := &AlterarSituacaoForm{
		VeiculoID:    veiculoID,
		VeiculoDesc:  veiculoDesc,
		NovaSituacao: situacaoAtual,
	}

	f.Form = huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title(fmt.Sprintf("Alterar Situação do Veículo #%d: %s", veiculoID, veiculoDesc)).
				Options(
					huh.NewOption("Disponível (Pronto para venda)", "Disponivel"),
					huh.NewOption("Reservado (Em negociação/sinal)", "Reservado"),
					huh.NewOption("Vendido (Baixa do estoque)", "Vendido"),
				).
				Value(&f.NovaSituacao),

			huh.NewInput().
				Title("Justificativa / Motivo da Alteração").
				Placeholder("Ex: Entrada de sinal pelo cliente, venda faturada, etc.").
				Value(&f.Justificativa).
				Validate(func(s string) error {
					if len(s) < 3 {
						return fmt.Errorf("a justificativa deve ter ao menos 3 caracteres")
					}
					return nil
				}),
		),
	).WithTheme(huh.ThemeCharm())

	return f
}

func (f *AlterarSituacaoForm) Update(msg tea.Msg) (*AlterarSituacaoForm, tea.Cmd) {
	if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.String() == "esc" {
		f.Cancelado = true
		return f, func() tea.Msg {
			return AlterarSituacaoResultMsg{Cancel: true}
		}
	}

	form, cmd := f.Form.Update(msg)
	if updated, ok := form.(*huh.Form); ok {
		f.Form = updated
	}

	if f.Form.State == huh.StateCompleted {
		return f, func() tea.Msg {
			return AlterarSituacaoResultMsg{
				VeiculoID:     f.VeiculoID,
				NovaSituacao:  f.NovaSituacao,
				Justificativa: f.Justificativa,
			}
		}
	}

	return f, cmd
}

func (f *AlterarSituacaoForm) View() string {
	return f.Form.View()
}
