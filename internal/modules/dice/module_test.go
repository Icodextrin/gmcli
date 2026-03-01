package dice

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
)

func newTestModel(t *testing.T) *Model {
	t.Helper()
	mod := New()
	m, ok := mod.(*Model)
	if !ok {
		t.Fatalf("New() returned %T, want *Model", mod)
	}
	m.SetSize(80, 24)
	return m
}

func enterKey() tea.KeyPressMsg { return tea.KeyPressMsg(tea.Key{Code: tea.KeyEnter}) }
func upKey() tea.KeyPressMsg    { return tea.KeyPressMsg(tea.Key{Code: tea.KeyUp}) }
func downKey() tea.KeyPressMsg  { return tea.KeyPressMsg(tea.Key{Code: tea.KeyDown}) }

func rollInput(t *testing.T, m *Model, input string) {
	t.Helper()
	m.input.SetValue(input)
	_, _ = m.Update(enterKey())
}

func TestEnterRollAddsMessageAndHistory(t *testing.T) {
	m := newTestModel(t)

	rollInput(t, m, "2d8")

	if len(m.messages) != 1 {
		t.Fatalf("messages = %d, want 1", len(m.messages))
	}
	if len(m.rolls) != 1 {
		t.Fatalf("roll history = %d, want 1", len(m.rolls))
	}
	if m.rolls[0] != "2d8" {
		t.Fatalf("history entry = %q, want %q", m.rolls[0], "2d8")
	}
	if m.historyIx != 1 {
		t.Fatalf("history index = %d, want 1", m.historyIx)
	}
	if got := m.input.Value(); got != "" {
		t.Fatalf("input should be reset after roll, got %q", got)
	}
	if !strings.Contains(m.messages[0], "2d8") {
		t.Fatalf("message %q does not contain roll expression", m.messages[0])
	}
}

func TestInvalidInputDoesNotAddHistory(t *testing.T) {
	m := newTestModel(t)
	m.input.SetValue("not-a-roll")

	_, _ = m.Update(enterKey())

	if len(m.messages) != 0 {
		t.Fatalf("messages = %d, want 0", len(m.messages))
	}
	if len(m.rolls) != 0 {
		t.Fatalf("roll history = %d, want 0", len(m.rolls))
	}
	if got := m.input.Value(); got != "not-a-roll" {
		t.Fatalf("input value = %q, want %q", got, "not-a-roll")
	}
}

func TestHistoryNavigationUpAndDown(t *testing.T) {
	m := newTestModel(t)

	rollInput(t, m, "1d6")
	rollInput(t, m, "2d4")

	_, _ = m.Update(upKey())
	if got := m.input.Value(); got != "2d4" {
		t.Fatalf("after first up, input = %q, want %q", got, "2d4")
	}

	_, _ = m.Update(upKey())
	if got := m.input.Value(); got != "1d6" {
		t.Fatalf("after second up, input = %q, want %q", got, "1d6")
	}

	_, _ = m.Update(downKey())
	if got := m.input.Value(); got != "2d4" {
		t.Fatalf("after first down, input = %q, want %q", got, "2d4")
	}

	_, _ = m.Update(downKey())
	if got := m.input.Value(); got != "1d6" {
		t.Fatalf("after second down, input = %q, want %q", got, "1d6")
	}
}

func TestUnfocusedModuleDoesNotRoll(t *testing.T) {
	m := newTestModel(t)
	m.SetFocus(false)
	m.input.SetValue("1d20")

	_, _ = m.Update(enterKey())

	if len(m.messages) != 0 {
		t.Fatalf("messages = %d, want 0", len(m.messages))
	}
	if len(m.rolls) != 0 {
		t.Fatalf("roll history = %d, want 0", len(m.rolls))
	}
	if m.input.Focused() {
		t.Fatal("input should be blurred when module is unfocused")
	}
}

func TestSetSizeEnforcesMinimumViewportHeight(t *testing.T) {
	m := newTestModel(t)

	m.SetSize(20, 1)

	if got := m.viewport.Height(); got != 3 {
		t.Fatalf("viewport height = %d, want 3", got)
	}
	if got := m.viewport.Width(); got != 20 {
		t.Fatalf("viewport width = %d, want 20", got)
	}
}
