package models

import (
	"fmt"
	"github.com/google/uuid"
	"math/rand"
	"time"
)

const numOfCards = 52

type Deck struct {
	UUID      uuid.UUID `json:"deck_id"`
	Shuffled  bool      `json:"shuffled"`
	Remaining int       `json:"remaining"`
	Cards     []Card    `json:"cards,omitempty"`
}

func NewDeck(shuffled bool, cards ...Card) *Deck {
	deck := new(Deck)
	deck.UUID = uuid.New()
	if len(cards) == 0 {
		deck.Remaining = numOfCards
		deck.Cards = fullListOfCards()
	} else {
		deck.Remaining = len(cards)
		deck.Cards = cards
	}

	if shuffled {
		deck.Shuffled = true
		deck.shuffle()
	}
	return deck
}

func (d *Deck) shuffle() {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(d.Cards), func(i, j int) { d.Cards[i], d.Cards[j] = d.Cards[j], d.Cards[i] })
}

func (d *Deck) DrawCards(cardsNum int) (cards []Card) {
	if cardsNum > d.Remaining {
		cardsNum = d.Remaining
	}
	cards = d.Cards[:cardsNum]
	d.Cards = d.Cards[cardsNum:]
	d.Remaining -= cardsNum
	return cards
}

func fullListOfCards() []Card {
	var cards = make([]Card, numOfCards)
	var index int

	for _, suite := range []rune("CDHS") {
		for i := 2; i <= 10; i++ {
			cards[index] = Card{Code: fmt.Sprintf("%d%c", i, suite)}
			index++
		}
		for _, i := range []rune("JQKA") {
			cards[index] = Card{Code: fmt.Sprintf("%c%c", i, suite)}
			index++
		}
	}
	return cards
}
