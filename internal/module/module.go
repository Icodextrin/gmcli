package module

import (
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
)

type ID string

const (
	DiceID       ID = "dice"
	InitiativeID ID = "initiative"
	NamesID      ID = "names"
)

type Module interface {
	ID() ID
	Title() string

	Init() tea.Cmd
	Update(tea.Msg) (Module, tea.Cmd)
	View() tea.View

	SetSize(width, height int)
	SetFocus(focused bool)

	ShortHelp() []key.Binding
	FullHelp() [][]key.Binding
}
