package models

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDeck_New(t *testing.T) {
	deck1 := NewDeck(false)
	deck2 := NewDeck(false)

	assert.Equal(t, deck1.Cards, deck2.Cards)
	assert.NotEqual(t, deck1.UUID, deck2.UUID)
	assert.Equal(t, 52, len(deck1.Cards))
}

func TestDeck_NewShuffle(t *testing.T) {
	deck1 := NewDeck(true)
	deck2 := NewDeck(true)

	assert.Equal(t, len(deck1.Cards), len(deck2.Cards))
	assert.NotEqual(t, deck1.Cards, deck2.Cards)
}

func TestDeck_NewPartial(t *testing.T) {
	cards := []Card{{"3D"}, {"4D"}, {"10D"}, {"AH"}}
	deck := NewDeck(false, cards...)

	assert.Equal(t, len(cards), deck.Remaining)
	assert.Equal(t, cards, deck.Cards)
}

func TestDeck_NewPartialShuffled(t *testing.T) {
	cards := []Card{{"3D"}, {"4D"}, {"10D"}, {"AH"}}
	var before []Card
	copy(before, cards)
	deck := NewDeck(true, cards...)

	assert.Equal(t, len(cards), deck.Remaining)
	assert.NotEqual(t, before, deck.Cards)
}

func TestDeck_DrawCards(t *testing.T) {
	tests := []struct {
		shuffle                                         bool
		cardsNumToDraw, cardsNumToBeDrawn, cardsNumLeft int
	}{
		{true, 10, 10, 42},
		{true, 100, 52, 0},
		{true, 0, 0, 52},
	}
	for _, test := range tests {
		deck := NewDeck(test.shuffle)
		drawnCards := deck.DrawCards(test.cardsNumToDraw)
		assert.Equal(t, test.cardsNumToBeDrawn, len(drawnCards))
		assert.Equal(t, test.cardsNumLeft, len(deck.Cards))
	}
}

func TestDeck_DrawCardsValues(t *testing.T) {
	tests := []struct {
		cardCodes                       []string
		drawnCardsCodes, leftCardsCodes []string
	}{
		{nil, []string{"2C", "3C", "4C"}, nil},
		{[]string{"10H", "JH"}, []string{"10H", "JH"}, []string{}},
		{[]string{"10H", "JH"}, []string{"10H"}, []string{"JH"}},
		{[]string{"10H", "JH"}, []string{}, []string{"10H", "JH"}},
		{[]string{"10H", "JH"}, []string{}, []string{"10H", "JH"}},
	}

	for _, test := range tests {
		cards := make([]Card, len(test.cardCodes))
		for i := range test.cardCodes {
			cards[i] = Card{Code: test.cardCodes[i]}
		}

		deck := NewDeck(false, cards...)
		drawnCards := deck.DrawCards(len(test.drawnCardsCodes))

		assert.Equal(t, len(test.drawnCardsCodes), len(drawnCards))
		for i := range test.drawnCardsCodes {
			assert.Equal(t, test.drawnCardsCodes[i], drawnCards[i].Code)
		}

		if test.leftCardsCodes != nil {
			assert.Equal(t, len(test.leftCardsCodes), len(deck.Cards))
			assert.Equal(t, len(test.leftCardsCodes), deck.Remaining)
			for i := range test.leftCardsCodes {
				assert.Equal(t, test.leftCardsCodes[i], deck.Cards[i].Code)
			}
			assert.NotContains(t, deck.Cards, drawnCards)
		}
	}
}
