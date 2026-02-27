# gmcli

`gmcli` is a terminal app for running tabletop RPG sessions, built with Go and the Bubble Tea TUI framework.

Right now the implemented module is a dice roller with roll history, critical result highlighting, and keyboard-driven interaction.

## What It Does

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
- `3#1d20+5` (roll one d20+5 three times)

Defaults:

- total rolls defaults to `1`
- number of dice defaults to `1`
- modifier defaults to `0`

## Keyboard Controls

- `enter`: roll current expression
- `up`: previous expression from history
- `down`: next expression from history
- `?`: toggle help display
- `q`, `esc`, `ctrl+c`: quit

## Project Structure

- `main.go`: Bubble Tea app model, input handling, viewport rendering, keybindings.
- `diceRoller.go`: dice parsing and roll logic.

## Tech Stack

- Go
- Bubble Tea v2 (`charm.land/bubbletea/v2`)
- Bubbles v2 (`charm.land/bubbles/v2`)
- Lip Gloss v2 (`charm.land/lipgloss/v2`)

## Run

```bash
go run .
```

## Build

```bash
go build -o gmcli .
```

## Current Scope

This repository is currently focused on the dice roller workflow. Future game-master utilities can be added as additional modules within the same TUI.
