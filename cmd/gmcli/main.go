package main

import (
	"fmt"
	"log"

	"gmcli/internal/app"
	dicemodule "gmcli/internal/modules/dice"
	initiativemodule "gmcli/internal/modules/initiative"
	namesmodule "gmcli/internal/modules/names"

	tea "charm.land/bubbletea/v2"
)

func main() {
	m := app.New(
		dicemodule.New(),
		initiativemodule.New(),
		namesmodule.New(),
	)
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Printf("ERROR: %v", err)
		log.Fatal(err)
	}
}
