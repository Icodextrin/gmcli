package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const gap = "\n\n"

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("ERROR: %v", err)
		log.Fatal(err)
	}
}

type (
	errMsg error
)

type keyMap struct {
	Up    key.Binding
	Down  key.Binding
	Enter key.Binding
	Help  key.Binding
	Quit  key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Enter, k.Help, k.Quit},
	}
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up"),
		key.WithHelp("↑", "prev roll"),
	),
	Down: key.NewBinding(
		key.WithKeys("down"),
		key.WithHelp("↓", "next roll"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "roll"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

type model struct {
	viewport         viewport.Model
	messages         []string
	rolls            []string
	rollsIndex       int
	textarea         textarea.Model
	senderStyle      lipgloss.Style
	critSuccessStyle lipgloss.Style
	critFailStyle    lipgloss.Style
	critBothStyle    lipgloss.Style
	helpStyle        lipgloss.Style
	keys             keyMap
	help             help.Model
	quitting         bool
	err              error
}

func initialModel() model {
	ta := textarea.New()
	ta.Placeholder = "<total num of rolls>#<num dice>d<num sides>[+,-]<modifier>"
	ta.Focus()

	ta.Prompt = "> "
	ta.CharLimit = 280

	ta.SetWidth(30)
	ta.SetHeight(1)

	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()

	ta.ShowLineNumbers = false

	vp := viewport.New(30, 5)
	vp.SetContent(`Rollem if you gottem...`)

	ta.KeyMap.InsertNewline.SetEnabled(false)

	return model{
		textarea:         ta,
		messages:         []string{},
		viewport:         vp,
		keys:             keys,
		help:             help.New(),
		helpStyle:        lipgloss.NewStyle().Foreground(lipgloss.Color("#FF75B7")),
		senderStyle:      lipgloss.NewStyle().Foreground(lipgloss.Color("5")),
		critSuccessStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("2")),
		critFailStyle:    lipgloss.NewStyle().Foreground(lipgloss.Color("1")),
		critBothStyle:    lipgloss.NewStyle().Foreground(lipgloss.Color("3")),
		err:              nil,
	}
}

func (m model) Init() tea.Cmd {
	return textarea.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		tiCmd tea.Cmd
		vpCmd tea.Cmd
	)

	m.textarea, tiCmd = m.textarea.Update(msg)
	m.viewport, vpCmd = m.viewport.Update(msg)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.viewport.Width = msg.Width
		m.textarea.SetWidth(msg.Width)
		m.viewport.Height = msg.Height - m.textarea.Height() - lipgloss.Height(gap)
		m.help.Width = msg.Width

		if len(m.messages) > 0 {
			m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(strings.Join(m.messages, "\n")))
		}
		m.viewport.GotoBottom()

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Quit):
			m.quitting = true
			return m, tea.Quit
		case key.Matches(msg, m.keys.Enter):
			userInput := m.textarea.Value()
			if userInput != "" {
				result, err := m.RollDiceString(userInput)
				if err != nil {
					return m, nil
				}

				m.messages = append(m.messages, m.senderStyle.Render(userInput)+": "+result)
				m.rolls = append(m.rolls, userInput)
				m.viewport.SetContent(lipgloss.NewStyle().Width(m.viewport.Width).Render(strings.Join(m.messages, "\n")))
				m.textarea.Reset()
				m.viewport.GotoBottom()
				m.rollsIndex = len(m.rolls)
			}
		case key.Matches(msg, m.keys.Up):
			if len(m.rolls) == 0 {
				return m, nil
			}
			m.rollsIndex--
			if m.rollsIndex < 0 {
				m.rollsIndex = len(m.rolls) - 1
			}
			m.textarea.Reset()
			m.textarea.SetValue(m.rolls[m.rollsIndex])
		case key.Matches(msg, m.keys.Down):
			if len(m.rolls) == 0 || m.rollsIndex == len(m.rolls) {
				return m, nil
			}
			m.rollsIndex++
			if m.rollsIndex == len(m.rolls) {
				m.rollsIndex = 0
			}
			m.textarea.Reset()
			m.textarea.SetValue(m.rolls[m.rollsIndex])
		case key.Matches(msg, m.keys.Help):
			m.help.ShowAll = !m.help.ShowAll
		}

	case errMsg:
		m.err = msg
		return m, nil
	}

	return m, tea.Batch(tiCmd, vpCmd)
}

func (m model) View() string {
	if m.quitting {
		return "Bye!\n"
	}

	helpView := m.help.View(m.keys)

	return fmt.Sprintf(
		"%s%s%s%s%s",
		m.viewport.View(),
		gap,
		m.textarea.View(),
		gap,
		helpView,
	)
}

func (m model) RollDiceString(userInput string) (string, error) {
	dice, err := ParseDiceString(userInput)
	if err != nil {
		return "", err
	}
	results := dice.Roll()

	var styledResults []string

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
