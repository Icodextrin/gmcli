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
	NumDice  int
	Sides    int
	Modifier int
}

// Roll executes the dice roll and returns the result
func (d DiceRoll) Roll() int {
	total := 0
	for range d.NumDice {
		total += rand.IntN(d.Sides) + 1
	}
	return total + d.Modifier
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
	re := regexp.MustCompile(`(\d+)d(\d+)([+-]\d+)?`)
	matches := re.FindStringSubmatch(input)

	if matches == nil {
		return nil, errors.New("invalid dice notation format")
	}

	// Parse number of dice (default to 1 if not specified)
	numDice := 1
	if matches[1] != "" {
		var err error
		numDice, err = strconv.Atoi(matches[1])
		if err != nil {
			return nil, fmt.Errorf("invalid number of dice: %s", matches[1])
		}
	}

	// Parse number of sides
	sides, err := strconv.Atoi(matches[2])
	if err != nil {
		return nil, fmt.Errorf("invalid number of sides: %s", matches[2])
	}

	// Parse modifier (default to 0 if not specified)
	modifier := 0
	if matches[3] != "" {
		modifier, err = strconv.Atoi(matches[3])
		if err != nil {
			return nil, fmt.Errorf("invalid modifier: %s", matches[3])
		}
	}

	// Validate inputs
	if numDice <= 0 {
		return nil, errors.New("number of dice must be positive")
	}
	if sides <= 0 {
		return nil, errors.New("number of sides must be positive")
	}

	return &DiceRoll{
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
