package names

import (
	"gmcli/internal/module"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type Model struct {
	width   int
	height  int
	focused bool
}

var _ module.Module = (*Model)(nil)

func New() module.Module {
	return &Model{
		focused: true,
	}
}

func (m *Model) ID() module.ID {
	return module.NamesID
}

func (m *Model) Title() string {
	return "Names"
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (module.Module, tea.Cmd) {
	return m, nil
}

func (m *Model) View() tea.View {
	body := "Name generator placeholder\n\nUse tab / shift+tab to switch modules."
	if m.width > 0 {
		body = lipgloss.NewStyle().Width(m.width).Render(body)
	}
	if m.height > 0 {
		body = lipgloss.NewStyle().Height(m.height).Render(body)
	}
	return tea.NewView(body)
}

func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
}

func (m *Model) SetFocus(focused bool) {
	m.focused = focused
}

func (m *Model) ShortHelp() []key.Binding {
	return []key.Binding{}
}

func (m *Model) FullHelp() [][]key.Binding {
	return [][]key.Binding{}
}
