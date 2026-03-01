package dice

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

type DiceRoll struct {
	NumRolls int
	NumDice  int
	Sides    int
	Modifier int
}

func (d DiceRoll) String() string {
	if d.Modifier == 0 {
		return fmt.Sprintf("%dd%d", d.NumDice, d.Sides)
	}
	if d.Modifier > 0 {
		return fmt.Sprintf("%dd%d+%d", d.NumDice, d.Sides, d.Modifier)
	}
	return fmt.Sprintf("%dd%d%d", d.NumDice, d.Sides, d.Modifier)
}

func ParseDiceString(input string) (*DiceRoll, error) {
	re := regexp.MustCompile(`(?:(\d+)#)?(\d*)d(\d+)([+-]\d+)?`)
	matches := re.FindStringSubmatch(input)
	if matches == nil {
		return nil, errors.New("invalid dice notation format")
	}

	rollsStr := matches[1]
	diceStr := matches[2]
	sidesStr := matches[3]
	modStr := matches[4]

	if rollsStr == "" {
		rollsStr = "1"
	}
	if diceStr == "" {
		diceStr = "1"
	}
	if modStr == "" {
		modStr = "+0"
	}
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

	modSign := modStr[0]
	modifier, err := strconv.Atoi(modStr[1:])
	if err != nil {
		return nil, fmt.Errorf("invalid modifier: %s", matches[4])
	}
	if modSign == '-' {
		modifier = -modifier
	}

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
