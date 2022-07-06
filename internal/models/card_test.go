package models

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCard_MarshalJSON(t *testing.T) {
	card := Card{Code: "10D"}
	data, err := card.MarshalJSON()
	assert.NoError(t, err)
	assert.Contains(t, string(data), "DIAMONDS")
	assert.Contains(t, string(data), "10")
}

func TestCard_UnmarshalJSON(t *testing.T) {
	cardJSONFull := []byte("{\"value\": \"ACE\", \"suit\": \"SPADES\", \"code\": \"AS\"}")
	cardJSONCode := []byte("{\"code\": \"AS\"}")

	var cardFull, cardCode Card
	err := json.Unmarshal(cardJSONCode, &cardCode)
	assert.NoError(t, err)
	err = json.Unmarshal(cardJSONFull, &cardFull)
	assert.NoError(t, err)
	assert.Equal(t, cardFull, cardCode)
}

func TestCard_IsCodeValid(t *testing.T) {
	assert.True(t, IsCodeValid("2C"))
	assert.True(t, IsCodeValid("10H"))
	assert.True(t, IsCodeValid("AD"))
	assert.False(t, IsCodeValid("11D"))
	assert.False(t, IsCodeValid("10"))
	assert.False(t, IsCodeValid("10X"))
	assert.False(t, IsCodeValid(""))
	assert.False(t, IsCodeValid("22D"))
}
