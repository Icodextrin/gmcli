# Implementing a New Module

This guide shows the expected pattern for adding a feature module.

## 1. Create Module Package

Create a new folder:

```text
internal/modules/<feature>/
```

Add `module.go` implementing `internal/module.Module`.

## 2. Optional: Create Domain Package

If the feature has reusable logic, create:

```text
internal/domain/<feature>/
```

Keep domain packages free of Bubble Tea/Bubbles/Lip Gloss imports.

## 3. Implement Required Methods

Use this minimum skeleton:

```go
type Model struct {
    width, height int
    focused       bool
}

func New() module.Module { return &Model{} }
func (m *Model) ID() module.ID { return module.ID("your-id") }
func (m *Model) Title() string { return "Your Module" }
func (m *Model) Init() tea.Cmd { return nil }
func (m *Model) Update(msg tea.Msg) (module.Module, tea.Cmd) { return m, nil }
func (m *Model) View() tea.View { return tea.NewView("placeholder") }
func (m *Model) SetSize(w, h int) { m.width, m.height = w, h }
func (m *Model) SetFocus(f bool) { m.focused = f }
func (m *Model) ShortHelp() []key.Binding { return nil }
func (m *Model) FullHelp() [][]key.Binding { return nil }
```

## 4. Register in Entrypoints

Wire the module in both:

- `main.go`
- `cmd/gmcli/main.go`

Example:

```go
m := app.New(
    dicemodule.New(),
    yourmodule.New(),
)
```

## 5. Keybinding Guidance

- Global bindings are handled in `internal/app/model.go`.
- Module-local bindings belong in the module package.
- Avoid global bindings that conflict with text input (for example numeric shortcuts if modules accept numeric input).

## 6. Sizing and Focus Checklist

- Respect `SetSize` and propagate dimensions to subcomponents.
- Respect `SetFocus`; inactive modules should not consume local keys.
- Return only module content from `View`; leave shell chrome to app layer.

## 7. Recommended Test Coverage

- Domain package unit tests first.
- Module-level tests for key handling, focus behavior, and size handling.
- App-level smoke test for module switching and quit behavior.
