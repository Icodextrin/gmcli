package app

import (
	"fmt"
	"gmcli/internal/module"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type keyMap struct {
	Next key.Binding
	Prev key.Binding
	Help key.Binding
	Quit key.Binding
}

func newKeyMap() keyMap {
	return keyMap{
		Next: key.NewBinding(key.WithKeys("tab"), key.WithHelp("tab", "next module")),
		Prev: key.NewBinding(key.WithKeys("shift+tab"), key.WithHelp("shift+tab", "prev module")),
		Help: key.NewBinding(key.WithKeys("?"), key.WithHelp("?", "toggle help")),
		Quit: key.NewBinding(key.WithKeys("q", "esc", "ctrl+c"), key.WithHelp("q", "quit")),
	}
}

type helpMap struct {
	global keyMap
	local  module.Module
}

func (h helpMap) ShortHelp() []key.Binding {
	out := []key.Binding{h.global.Next, h.global.Help, h.global.Quit}
	if h.local != nil {
		out = append(out, h.local.ShortHelp()...)
	}
	return out
}

func (h helpMap) FullHelp() [][]key.Binding {
	rows := [][]key.Binding{
		{h.global.Next, h.global.Prev},
		{h.global.Help, h.global.Quit},
	}
	if h.local != nil {
		rows = append(rows, h.local.FullHelp()...)
	}
	return rows
}

type Model struct {
	width, height int
	ready         bool

	modules map[module.ID]module.Module
	order   []module.ID
	active  module.ID

	keys     keyMap
	help     help.Model
	showHelp bool
}

func New(mods ...module.Module) *Model {
	if len(mods) == 0 {
		panic("app.New requires at least one module")
	}

	modsByID := make(map[module.ID]module.Module, len(mods))
	order := make([]module.ID, 0, len(mods))
	for _, mod := range mods {
		id := mod.ID()
		modsByID[id] = mod
		order = append(order, id)
		mod.SetFocus(false)
	}

	active := order[0]
	modsByID[active].SetFocus(true)

	return &Model{
		modules:  modsByID,
		order:    order,
		active:   active,
		keys:     newKeyMap(),
		help:     help.New(),
		showHelp: true,
	}
}

func (m *Model) Init() tea.Cmd {
	return m.modules[m.active].Init()
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.help.SetWidth(msg.Width)
		m.ready = true
		m.resizeActive()
		return m, nil

	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Help):
			m.showHelp = !m.showHelp
			m.resizeActive()
			return m, nil
		case key.Matches(msg, m.keys.Next):
			m.activateByOffset(+1)
			return m, nil
		case key.Matches(msg, m.keys.Prev):
			m.activateByOffset(-1)
			return m, nil
		}
	}

	activeMod := m.modules[m.active]
	nextMod, cmd := activeMod.Update(msg)
	m.modules[m.active] = nextMod
	return m, cmd
}

func (m *Model) View() tea.View {
	if !m.ready {
		v := tea.NewView("Initializing...")
		v.AltScreen = true
		return v
	}

	header := m.renderTabs()
	content := m.modules[m.active].View().Content

	helpView := ""
	if m.showHelp {
		helpView = "\n" + m.help.View(helpMap{
			global: m.keys,
			local:  m.modules[m.active],
		})
	}

	v := tea.NewView(fmt.Sprintf("%s\n%s%s", header, content, helpView))
	v.AltScreen = true
	return v
}

func (m *Model) resizeActive() {
	const chromeRows = 2
	helpRows := 0
	if m.showHelp {
		helpRows = 2
	}
	h := m.height - chromeRows - helpRows
	if h < 3 {
		h = 3
	}
	m.modules[m.active].SetSize(m.width, h)
}

func (m *Model) activate(id module.ID) {
	if _, ok := m.modules[id]; !ok || id == m.active {
		return
	}
	m.modules[m.active].SetFocus(false)
	m.active = id
	m.modules[m.active].SetFocus(true)
	m.resizeActive()
}

func (m *Model) activateByOffset(delta int) {
	if len(m.order) == 0 {
		return
	}

	i := 0
	for idx, id := range m.order {
		if id == m.active {
			i = idx
			break
		}
	}
	n := (i + delta + len(m.order)) % len(m.order)
	m.activate(m.order[n])
}

func (m *Model) renderTabs() string {
	base := lipgloss.NewStyle().Padding(0, 1)
	active := lipgloss.NewStyle().Padding(0, 1).Bold(true).Underline(true)

	out := ""
	for i, id := range m.order {
		label := m.modules[id].Title()
		if id == m.active {
			out += active.Render(label)
		} else {
			out += base.Render(label)
		}
		if i < len(m.order)-1 {
			out += " "
		}
	}
	return out
}
