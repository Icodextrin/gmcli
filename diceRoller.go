package main

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"regexp"
	"strconv"
)

// DiceRoll represents a parsed dice expression
type DiceRoll struct {
	NumRolls int
	NumDice  int
	Sides    int
	Modifier int
}

type RollResult struct {
	Total       int
	HasCritFail bool
	HasCritSucc bool
}

func (r RollResult) CritStatus() string {
	if r.HasCritFail && r.HasCritSucc {
		return "both"
	} else if r.HasCritSucc {
		return "success"
	} else if r.HasCritFail {
		return "fail"
	} else {
		return "normal"
	}
}

// Roll executes the dice roll and returns the result
// Tracks crit failure and crit success for later formatting
func (d DiceRoll) Roll() []RollResult {
	rollList := []RollResult{}
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
		rollList = append(rollList, RollResult{
			Total:       total + d.Modifier,
			HasCritFail: hasCritFail,
			HasCritSucc: hasCritSucc,
		})
	}
	return rollList
}

// String returns a string representation of the dice roll
func (d DiceRoll) String() string {
	if d.Modifier == 0 {
		return fmt.Sprintf("%dd%d", d.NumDice, d.Sides)
	} else if d.Modifier > 0 {
		return fmt.Sprintf("%dd%d+%d", d.NumDice, d.Sides, d.Modifier)
	} else {
		return fmt.Sprintf("%dd%d%d", d.NumDice, d.Sides, d.Modifier)
	}
}

// ParseDiceString parses a dice string like "2d20+5" into a DiceRoll
func ParseDiceString(input string) (*DiceRoll, error) {
	// Match dice notation: optional number, 'd', number, optional modifier
	re := regexp.MustCompile(`(?:(\d+)#)?(\d*)d(\d+)([+-]\d+)?`)
	matches := re.FindStringSubmatch(input)

	if matches == nil {
		return nil, errors.New("invalid dice notation format")
	}
	rollsStr := matches[1] // Number of times to roll
	diceStr := matches[2]  // Number of dice
	sidesStr := matches[3] // Number of sides (required)
	modStr := matches[4]   // Modifier (+/- number)

	// defaults
	if rollsStr == "" {
		rollsStr = "1"
	}
	if diceStr == "" {
		diceStr = "1"
	}
	if modStr == "" {
		modStr = "+0"
	}
	// number of sides is required
	if sidesStr == "" {
		return nil, errors.New("invalid dice notation format")
	}

	numRolls, err := strconv.Atoi(rollsStr)
	if err != nil {
		return nil, fmt.Errorf("invalid number of rolls: %s", rollsStr)
	}
	numDice, err := strconv.Atoi(diceStr)
	if err != nil {
		return nil, fmt.Errorf("invalid number of dice: %s", diceStr)
	}
	sides, err := strconv.Atoi(sidesStr)
	if err != nil {
		return nil, fmt.Errorf("invalid number of sides: %s", sidesStr)
	}

	mod := modStr[0]

	modifier, err := strconv.Atoi(modStr[1:])
	if err != nil {
		return nil, fmt.Errorf("invalid modifier: %s", matches[4])
	}
	if mod == '-' {
		modifier = -modifier
	}

	// Validate inputs
	if numRolls <= 0 {
		return nil, errors.New("number of rolls must be positive")
	}
	if numDice <= 0 {
		return nil, errors.New("number of dice must be positive")
	}
	if sides <= 0 {
		return nil, errors.New("number of sides must be positive")
	}

	return &DiceRoll{
		NumRolls: numRolls,
		NumDice:  numDice,
		Sides:    sides,
		Modifier: modifier,
	}, nil
}

//func main() {
//	testInputs := []string{
//		"1d20+3",
//		"2d6",
//		"d8+2",
//		"3d10-1",
//		"4d4+5",
//		"1d20 + 3", // invalid - contains spaces
//		"invalid",
//	}
//
//	for _, input := range testInputs {
//		fmt.Printf("Input: %s\n", input)
//
//		dice, err := ParseDiceString(input)
//		if err != nil {
//			fmt.Printf("  Error: %v\n", err)
//			continue
//		}
//
//		fmt.Printf("  Parsed: %s\n", dice)
//		fmt.Printf("  Roll result: %d\n", dice.Roll())
//		fmt.Println()
//	}
//}
