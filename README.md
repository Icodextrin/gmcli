# gmcli

`gmcli` is a terminal app for running tabletop RPG sessions, built with Go and the Bubble Tea TUI framework.

The app uses a modular shell architecture:

- an app-level shell model handles global layout, focus, and module navigation
- each feature lives in its own module package behind a shared interface
- pure game logic is split into domain packages separate from TUI code

## Current Modules

- `Dice`: implemented, with history and crit highlighting
- `Initiative`: placeholder module
- `Names`: placeholder module

## What Dice Does

- Parses dice notation from the input prompt.
- Rolls one or more dice expressions.
- Prints results in a scrollable log.
- Tracks roll history so previous expressions can be recalled.
- Highlights critical outcomes:
  - crit success (at least one die rolled max value)
  - crit fail (at least one die rolled `1`)
  - both (both conditions happened in the same roll)

## Dice Notation

Format:

`<total rolls>#<num dice>d<num sides>[+|-modifier]`

Examples:

- `d20`
- `2d6+3`
- `4d8-1`
- `3#2d20+5` (roll two d20+5 three times)

Defaults:

- total rolls defaults to `1`
- number of dice defaults to `1`
- modifier defaults to `0`

## Keyboard Controls

- `tab`: next module
- `shift+tab`: previous module
- `?`: toggle help
- `q`, `esc`, `ctrl+c`: quit

Dice module local keys:

- `enter`: roll current expression
- `up`: previous expression from history
- `down`: next expression from history

## Project Structure

```text
.
├── main.go                      # convenience entrypoint
├── cmd/gmcli/main.go            # command entrypoint
├── internal/
│   ├── app/
│   │   ├── model.go             # shell model/router/layout
│   │   └── messages.go          # shared app-level messages
│   ├── module/
│   │   └── module.go            # shared Module interface
│   ├── domain/
│   │   └── dice/                # pure dice parsing/rolling logic
│   └── modules/
│       ├── dice/                # dice UI module
│       ├── initiative/          # placeholder module
│       └── names/               # placeholder module
└── docs/
    ├── ARCHITECTURE.md
    ├── MODULE_INTERFACE.md
    └── IMPLEMENTING_MODULES.md
```

## Module Contract

All modules implement `internal/module.Module` so the shell can route messages and render any module consistently.

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

See detailed docs in [docs/MODULE_INTERFACE.md](docs/MODULE_INTERFACE.md).

## Tech Stack

- Go
- Bubble Tea v2 (`charm.land/bubbletea/v2`)
- Bubbles v2 (`charm.land/bubbles/v2`)
- Lip Gloss v2 (`charm.land/lipgloss/v2`)

## Run

```bash
go run .
```

Or:

```bash
go run ./cmd/gmcli
```

## Build

```bash
go build -o gmcli .
```

## Documentation

- [Architecture](docs/ARCHITECTURE.md)
- [Module Interface](docs/MODULE_INTERFACE.md)
- [Implementing Modules](docs/IMPLEMENTING_MODULES.md)
