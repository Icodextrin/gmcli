# Testing Guide

This project uses standard Go tests (`testing` package) for shell and module behavior.

## Run Tests

Run all tests:

```bash
go test ./...
```

Run one package:

```bash
go test ./internal/app
go test ./internal/modules/dice
```

Run one test by name:

```bash
go test ./internal/modules/dice -run TestEnterRollAddsMessageAndHistory
```

## Current Test Files

- `internal/app/model_test.go`
- `internal/modules/dice/module_test.go`

## Testing Strategy

1. Shell tests (`internal/app`)
- Use a fake module that implements `internal/module.Module`.
- Assert global behavior:
  - window resize handling
  - tab navigation and focus changes
  - help toggle and quit behavior
  - message forwarding rules

2. Module tests (`internal/modules/dice`)
- Cast `New()` result to concrete `*Model` for direct state assertions.
- Send key events as `tea.KeyPressMsg`.
- Assert module-local behavior:
  - valid roll updates history and messages
  - invalid roll does not mutate history
  - up/down recall behavior
  - unfocused module ignores local actions
  - sizing logic applies minimum viewport constraints

## Key Testing Patterns

Create key messages:

```go
tea.KeyPressMsg(tea.Key{Code: tea.KeyEnter})
tea.KeyPressMsg(tea.Key{Code: tea.KeyTab})
tea.KeyPressMsg(tea.Key{Code: tea.KeyTab, Mod: tea.ModShift})
tea.KeyPressMsg(tea.Key{Text: "q", Code: 'q'})
```

Test quit command:

```go
_, cmd := m.Update(tea.KeyPressMsg(tea.Key{Text: "q", Code: 'q'}))
msg := cmd()
_, ok := msg.(tea.QuitMsg)
```

## Guidelines for New Tests

- Keep tests deterministic where possible.
- Prefer behavior assertions over implementation details.
- For shell tests, verify routing/focus/resizing first.
- For module tests, verify local state transitions and key handling.
- Add domain tests in `internal/domain/<feature>` when business logic grows.
