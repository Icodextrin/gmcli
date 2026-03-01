package app

import (
	"testing"

	"gmcli/internal/module"

	"charm.land/bubbles/v2/key"
	tea "charm.land/bubbletea/v2"
)

type markerMsg struct{}

type fakeModule struct {
	id    module.ID
	title string

	focused      bool
	focusCalls   []bool
	sizeCalls    int
	lastWidth    int
	lastHeight   int
	updateCalls  int
	lastUpdate   tea.Msg
	initCmd      tea.Cmd
	shortHelpOut []key.Binding
	fullHelpOut  [][]key.Binding
}

func newFakeModule(id module.ID, title string) *fakeModule {
	return &fakeModule{id: id, title: title}
}

func (m *fakeModule) ID() module.ID { return m.id }
func (m *fakeModule) Title() string { return m.title }
func (m *fakeModule) Init() tea.Cmd { return m.initCmd }
func (m *fakeModule) Update(msg tea.Msg) (module.Module, tea.Cmd) {
	m.updateCalls++
	m.lastUpdate = msg
	return m, nil
}
func (m *fakeModule) View() tea.View { return tea.NewView(m.title) }
func (m *fakeModule) SetSize(width, height int) {
	m.sizeCalls++
	m.lastWidth = width
	m.lastHeight = height
}
func (m *fakeModule) SetFocus(focused bool) {
	m.focused = focused
	m.focusCalls = append(m.focusCalls, focused)
}
func (m *fakeModule) ShortHelp() []key.Binding  { return m.shortHelpOut }
func (m *fakeModule) FullHelp() [][]key.Binding { return m.fullHelpOut }
func textKey(s string) tea.KeyPressMsg          { return tea.KeyPressMsg(tea.Key{Text: s, Code: []rune(s)[0]}) }
func codeKey(code rune, mods ...tea.KeyMod) tea.KeyPressMsg {
	var mod tea.KeyMod
	for _, m := range mods {
		mod |= m
	}
	return tea.KeyPressMsg(tea.Key{Code: code, Mod: mod})
}
func setWindowSize(m *Model, width, height int) {
	_, _ = m.Update(tea.WindowSizeMsg{Width: width, Height: height})
}
func assertActive(t *testing.T, m *Model, id module.ID) {
	t.Helper()
	if m.active != id {
		t.Fatalf("active module = %q, want %q", m.active, id)
	}
}
func assertFocused(t *testing.T, m *fakeModule, focused bool) {
	t.Helper()
	if m.focused != focused {
		t.Fatalf("focused = %t, want %t", m.focused, focused)
	}
}
func assertSize(t *testing.T, m *fakeModule, w, h int) {
	t.Helper()
	if m.lastWidth != w || m.lastHeight != h {
		t.Fatalf("size = (%d,%d), want (%d,%d)", m.lastWidth, m.lastHeight, w, h)
	}
}
func assertUpdateCalls(t *testing.T, m *fakeModule, n int) {
	t.Helper()
	if m.updateCalls != n {
		t.Fatalf("update calls = %d, want %d", m.updateCalls, n)
	}
}
func assertFocusCallCount(t *testing.T, m *fakeModule, n int) {
	t.Helper()
	if len(m.focusCalls) != n {
		t.Fatalf("focus calls = %d, want %d", len(m.focusCalls), n)
	}
}
func assertContainsMsgType[T any](t *testing.T, msg tea.Msg) {
	t.Helper()
	if _, ok := msg.(T); !ok {
		t.Fatalf("message type = %T, want %T", msg, *new(T))
	}
}
func assertCmdMsgType[T any](t *testing.T, cmd tea.Cmd) {
	t.Helper()
	if cmd == nil {
		t.Fatal("command was nil")
	}
	assertContainsMsgType[T](t, cmd())
}

func TestNewPanicsWithoutModules(t *testing.T) {
	t.Helper()
	defer func() {
		if recover() == nil {
			t.Fatal("expected panic when constructing app without modules")
		}
	}()
	_ = New()
}

func TestInitDelegatesToActiveModule(t *testing.T) {
	a := newFakeModule(module.DiceID, "Dice")
	a.initCmd = func() tea.Msg { return markerMsg{} }
	b := newFakeModule(module.NamesID, "Names")

	m := New(a, b)
	cmd := m.Init()
	if cmd == nil {
		t.Fatal("expected init command from active module")
	}
	assertContainsMsgType[markerMsg](t, cmd())
}

func TestWindowSizeResizesActiveModule(t *testing.T) {
	a := newFakeModule(module.DiceID, "Dice")
	b := newFakeModule(module.NamesID, "Names")
	m := New(a, b)

	setWindowSize(m, 120, 42)

	if !m.ready {
		t.Fatal("expected model to be ready after window size message")
	}
	assertSize(t, a, 120, 38)
	if b.sizeCalls != 0 {
		t.Fatalf("inactive module resize calls = %d, want 0", b.sizeCalls)
	}
}

func TestTabSwitchesModuleAndFocus(t *testing.T) {
	a := newFakeModule(module.DiceID, "Dice")
	b := newFakeModule(module.InitiativeID, "Initiative")
	m := New(a, b)
	setWindowSize(m, 100, 30)

	_, _ = m.Update(codeKey(tea.KeyTab))

	assertActive(t, m, module.InitiativeID)
	assertFocused(t, a, false)
	assertFocused(t, b, true)
	assertSize(t, b, 100, 26)
	assertUpdateCalls(t, a, 0)
	assertUpdateCalls(t, b, 0)
	assertFocusCallCount(t, a, 3) // false(initial), true(active), false(switched away)
	assertFocusCallCount(t, b, 2) // false(initial), true(switched to)
}

func TestShiftTabCyclesBackward(t *testing.T) {
	a := newFakeModule(module.DiceID, "Dice")
	b := newFakeModule(module.InitiativeID, "Initiative")
	c := newFakeModule(module.NamesID, "Names")
	m := New(a, b, c)
	setWindowSize(m, 100, 30)

	_, _ = m.Update(codeKey(tea.KeyTab, tea.ModShift))

	assertActive(t, m, module.NamesID)
	assertFocused(t, a, false)
	assertFocused(t, c, true)
}

func TestGlobalKeysAreNotForwardedToModules(t *testing.T) {
	a := newFakeModule(module.DiceID, "Dice")
	m := New(a)
	setWindowSize(m, 100, 30)

	_, _ = m.Update(textKey("?"))
	_, _ = m.Update(codeKey(tea.KeyTab))

	assertUpdateCalls(t, a, 0)
	if m.showHelp {
		t.Fatal("expected help to be hidden after pressing '?'")
	}
}

func TestNonGlobalKeyIsForwardedToActiveModule(t *testing.T) {
	a := newFakeModule(module.DiceID, "Dice")
	b := newFakeModule(module.InitiativeID, "Initiative")
	m := New(a, b)

	_, _ = m.Update(textKey("2"))

	assertActive(t, m, module.DiceID)
	assertUpdateCalls(t, a, 1)
	assertUpdateCalls(t, b, 0)
}

func TestQuitReturnsQuitCommand(t *testing.T) {
	a := newFakeModule(module.DiceID, "Dice")
	m := New(a)

	_, cmd := m.Update(textKey("q"))
	assertCmdMsgType[tea.QuitMsg](t, cmd)
}
