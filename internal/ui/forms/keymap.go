package forms

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/huh"
)

// FormKeyMap configura atalhos aprimorados permitindo que o usuário retroceda
// aos campos anteriores utilizando Shift+Tab, Backtab ou Seta para Cima (↑),
// além de avançar com Enter, Tab ou Seta para Baixo (↓).
func FormKeyMap() *huh.KeyMap {
	km := huh.NewDefaultKeyMap()
	km.Input.Prev = key.NewBinding(
		key.WithKeys("shift+tab", "backtab", "up"),
		key.WithHelp("↑/shift+tab", "campo anterior"),
	)
	km.Input.Next = key.NewBinding(
		key.WithKeys("enter", "tab", "down"),
		key.WithHelp("↓/enter/tab", "próximo campo"),
	)
	km.Select.Prev = key.NewBinding(
		key.WithKeys("shift+tab", "backtab"),
		key.WithHelp("shift+tab", "campo anterior"),
	)
	return km
}
