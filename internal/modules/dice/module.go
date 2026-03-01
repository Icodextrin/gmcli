package dice

import (
	"fmt"
	"gmcli/internal/domain/dice"
	"gmcli/internal/module"
	"strings"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/textarea"
	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

const (
	gap             = "\n\n"
	placeholderText = "Rollem if you gottem..."
)

type keyMap struct {
	Up    key.Binding
	Down  key.Binding
	Enter key.Binding
}

func newKeyMap() keyMap {
	return keyMap{
		Up:    key.NewBinding(key.WithKeys("up"), key.WithHelp("up", "prev roll")),
		Down:  key.NewBinding(key.WithKeys("down"), key.WithHelp("down", "next roll")),
		Enter: key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "roll")),
	}
}

type Model struct {
	viewport viewport.Model
	input    textarea.Model

	messages  []string
	rolls     []string
	historyIx int
	focused   bool

	keys keyMap

	rollStyle        lipgloss.Style
	critSuccessStyle lipgloss.Style
	critFailStyle    lipgloss.Style
	critBothStyle    lipgloss.Style
	placeholderStyle lipgloss.Style
}

var _ module.Module = (*Model)(nil)

func New() module.Module {
	in := textarea.New()
	in.Placeholder = "<total num of rolls>#<num dice>d<num sides>[+,-]<modifier>"
	in.Prompt = "> "
	in.CharLimit = 280
	in.SetHeight(1)
	in.SetWidth(30)
	in.ShowLineNumbers = false
	in.KeyMap.InsertNewline.SetEnabled(false)
	in.Focus()

	styles := in.Styles()
	styles.Focused.CursorLine = lipgloss.NewStyle()
	in.SetStyles(styles)

	vp := viewport.New(
		viewport.WithWidth(30),
		viewport.WithHeight(5),
	)
	vp.SetContent(lipgloss.NewStyle().Foreground(lipgloss.Color("8")).Render(placeholderText))

	return &Model{
		viewport:         vp,
		input:            in,
		keys:             newKeyMap(),
		focused:          true,
		rollStyle:        lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
		critSuccessStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("2")),
		critFailStyle:    lipgloss.NewStyle().Foreground(lipgloss.Color("1")),
		critBothStyle:    lipgloss.NewStyle().Foreground(lipgloss.Color("3")),
		placeholderStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("8")),
	}
}

func (m *Model) ID() module.ID {
	return module.DiceID
}

func (m *Model) Title() string {
	return "Dice Roller"
}

func (m *Model) Init() tea.Cmd {
	return textarea.Blink
}

func (m *Model) SetSize(width, height int) {
	m.input.SetWidth(width)
	vh := height - m.input.Height() - lipgloss.Height(gap)
	if vh < 3 {
		vh = 3
	}
	m.viewport.SetWidth(width)
	m.viewport.SetHeight(vh)
	m.refreshViewport()
	m.viewport.GotoBottom()
}

func (m *Model) SetFocus(focused bool) {
	m.focused = focused
	if focused {
		m.input.Focus()
	} else {
		m.input.Blur()
	}
}

func (m *Model) Update(msg tea.Msg) (module.Module, tea.Cmd) {
	var inCmd tea.Cmd
	m.input, inCmd = m.input.Update(msg)
	var vpCmd tea.Cmd
	m.viewport, vpCmd = m.viewport.Update(msg)

	if !m.focused {
		return m, tea.Batch(inCmd, vpCmd)
	}

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch {
		case key.Matches(msg, m.keys.Enter):
			userInput := strings.TrimSpace(m.input.Value())
			if userInput == "" {
				return m, tea.Batch(inCmd, vpCmd)
			}

			if strings.EqualFold(userInput, "clear") {
				m.clearHistory()
				m.input.Reset()
				return m, tea.Batch(inCmd, vpCmd)
			}

			result, err := m.rollDiceString(userInput)
			if err != nil {
				return m, tea.Batch(inCmd, vpCmd)
			}
			m.messages = append(m.messages, m.rollStyle.Render(userInput)+": "+result)
			m.rolls = append(m.rolls, userInput)
			m.historyIx = len(m.rolls)
			m.refreshViewport()
			m.input.Reset()
			m.viewport.GotoBottom()

		case key.Matches(msg, m.keys.Up):
			if len(m.rolls) == 0 {
				return m, tea.Batch(inCmd, vpCmd)
			}
			m.historyIx--
			if m.historyIx < 0 {
				m.historyIx = len(m.rolls) - 1
			}
			m.input.Reset()
			m.input.SetValue(m.rolls[m.historyIx])

		case key.Matches(msg, m.keys.Down):
			if len(m.rolls) == 0 || m.historyIx == len(m.rolls) {
				return m, tea.Batch(inCmd, vpCmd)
			}
			m.historyIx++
			if m.historyIx >= len(m.rolls) {
				m.historyIx = 0
			}
			m.input.Reset()
			m.input.SetValue(m.rolls[m.historyIx])
		}
	}

	return m, tea.Batch(inCmd, vpCmd)
}

func (m *Model) View() tea.View {
	return tea.NewView(m.viewport.View() + gap + m.input.View())
}

func (m *Model) ShortHelp() []key.Binding {
	return []key.Binding{m.keys.Enter, m.keys.Up, m.keys.Down}
}

func (m *Model) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{m.keys.Enter, m.keys.Up, m.keys.Down},
	}
}

func (m *Model) refreshViewport() {
	base := lipgloss.NewStyle().Width(m.viewport.Width())
	if len(m.messages) == 0 {
		m.viewport.SetContent(base.Render(m.placeholderStyle.Render(placeholderText)))
		return
	}

	content := strings.Join(m.messages, "\n")
	m.viewport.SetContent(base.Render(content))
}

func (m *Model) rollDiceString(userInput string) (string, error) {
	d, err := dice.ParseDiceString(userInput)
	if err != nil {
		return "", err
	}
	results := d.Roll()

	styledResults := make([]string, 0, len(results))
	for _, result := range results {
		resultStr := fmt.Sprintf("%d", result.Total)
		switch result.CritStatus() {
		case "success":
			resultStr = m.critSuccessStyle.Render(resultStr)
		case "fail":
			resultStr = m.critFailStyle.Render(resultStr)
		case "both":
			resultStr = m.critBothStyle.Render(resultStr)
		}
		styledResults = append(styledResults, resultStr)
	}
	return strings.Join(styledResults, " "), nil
}

func (m *Model) clearHistory() {
	m.messages = nil
	m.rolls = nil
	m.historyIx = 0
	m.refreshViewport()
}
