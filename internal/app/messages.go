package app

import "gmcli/internal/module"

type SwitchModuleMsg struct {
	Target module.ID
}

type StatusMsg struct {
	Text string
}

type DiceRolledMsg struct {
	Expr   string
	Totals []int
}
