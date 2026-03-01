# gmcli Architecture

This document explains the current application structure after the modular refactor.

## Layers

1. `internal/app`
- The shell (`Model`) is the top-level Bubble Tea model.
- Owns global keybindings, help display, active module selection, and layout.
- Forwards non-global messages to the active module.

2. `internal/module`
- Defines the shared `Module` interface.
- Every feature module must implement this contract.

3. `internal/modules/<feature>`
- UI behavior for each feature.
- Manages local inputs, views, and feature-specific key handling.
- Should not contain reusable game-domain logic.

4. `internal/domain/<feature>`
- Pure domain logic with no Bubble Tea/Bubbles/Lip Gloss dependencies.
- Example: dice parsing and rolling.

## Message and Control Flow

1. Bubble Tea sends `tea.Msg` to the shell (`internal/app.Model`).
2. Shell handles global keys (`tab`, `shift+tab`, `?`, quit).
3. If not globally handled, shell forwards the message to the active module.
4. Module updates local state and returns `(Module, tea.Cmd)`.
5. Shell renders tabs/chrome + active module view + help.

## Focus and Sizing Rules

- Shell owns the full terminal dimensions.
- Shell computes each module's usable content area and calls `SetSize(width, height)`.
- Shell controls active/inactive state via `SetFocus(bool)`.
- Modules should ignore local interaction when unfocused.

## Why This Structure

- Keeps modules independent and replaceable.
- Makes adding new features mostly a matter of adding one module package and wiring it in.
- Separates domain logic from TUI concerns, improving testability.
