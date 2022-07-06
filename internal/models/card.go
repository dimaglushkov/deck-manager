package models

import (
	"encoding/json"
	"fmt"
)

var cardValues map[string]string
var cardSuites map[rune]string

func init() {
	cardValues = make(map[string]string, 13)
	for i := 2; i <= 10; i++ {
		str := fmt.Sprintf("%d", i)
		cardValues[str] = str
	}
	cardValues["J"] = "JACK"
	cardValues["Q"] = "QUEEN"
	cardValues["K"] = "KING"
	cardValues["A"] = "ACE"

	cardSuites = map[rune]string{
		'C': "CLUBS",
		'H': "HEARTS",
		'S': "SPADES",
		'D': "DIAMONDS",
	}
}

type Card struct {
	Code string
}

func (c *Card) UnmarshalJSON(data []byte) error {
	var d map[string]interface{}
	if err := json.Unmarshal(data, &d); err != nil {
		return err
	}

	c.Code = d["code"].(string)
	return nil
}

func (c *Card) MarshalJSON() ([]byte, error) {
	var data = make(map[string]string)
	var n = len(c.Code)

	data["value"] = cardValues[c.Code[:n-1]]
	data["suite"] = cardSuites[[]rune(c.Code)[n-1]]
	data["code"] = c.Code

	return json.Marshal(data)
}

func IsCodeValid(code string) bool {
	n := len(code)
	if n != 2 && n != 3 {
		return false
	}

	value := code[:n-1]
	suite := []rune(code)[n-1]

	if _, ok := cardValues[value]; !ok {
		return false
	}

	if _, ok := cardSuites[suite]; !ok {
		return false
	}
	return true
}
