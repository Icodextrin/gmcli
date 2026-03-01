package dice

import "math/rand/v2"

type RollResult struct {
	Total       int
	HasCritFail bool
	HasCritSucc bool
}

func (r RollResult) CritStatus() string {
	if r.HasCritFail && r.HasCritSucc {
		return "both"
	}
	if r.HasCritSucc {
		return "success"
	}
	if r.HasCritFail {
		return "fail"
	}
	return "normal"
}

func (d DiceRoll) Roll() []RollResult {
	rolls := make([]RollResult, 0, d.NumRolls)
	for range d.NumRolls {
		total := 0
		hasCritFail := false
		hasCritSucc := false

		for range d.NumDice {
			dieResult := rand.IntN(d.Sides) + 1
			total += dieResult
			if dieResult == 1 {
				hasCritFail = true
			}
			if dieResult == d.Sides {
				hasCritSucc = true
			}
		}

		rolls = append(rolls, RollResult{
			Total:       total + d.Modifier,
			HasCritFail: hasCritFail,
			HasCritSucc: hasCritSucc,
		})
	}
	return rolls
}
