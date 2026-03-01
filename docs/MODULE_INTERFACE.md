# Module Interface

The core contract between the shell and feature modules lives in `internal/module/module.go`.

## Interface Definition

```go
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
```

## Method Responsibilities

`ID() ID`
- Returns stable module identity (`dice`, `initiative`, `names`).
- Used by shell registry and navigation.

`Title() string`
- Human-readable tab/header label.

`Init() tea.Cmd`
- Startup command for the module.
- Return `nil` if no startup command is needed.

`Update(tea.Msg) (Module, tea.Cmd)`
- Handles module-level events and state transitions.
- Returns updated module and optional command.

`View() tea.View`
- Returns the module's content view only.
- Shell composes global chrome around it.

`SetSize(width, height int)`
- Called by shell on window/layout changes.
- Module must resize internal components (viewport/table/inputs).

`SetFocus(focused bool)`
- Called by shell when switching active modules.
- Module should stop local key handling when unfocused.

`ShortHelp() []key.Binding`
- Compact help row for active module.

`FullHelp() [][]key.Binding`
- Expanded help sections for active module.

## Behavioral Expectations

- Avoid direct dependencies between modules.
- Keep domain logic in `internal/domain`.
- Process local keybindings only when focused.
- Prefer deterministic, side-effect-light `Update` logic.
