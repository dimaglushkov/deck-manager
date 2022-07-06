package internal

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/google/uuid"

	"github.com/dimaglushkov/toggl-test-assignment/internal/models"
)

type handler struct {
	repo Repository
}

func NewHandler(repo Repository) *handler {
	return &handler{repo: repo}
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	args, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch r.URL.Path {
	case "/create":
		h.handleCreateDeck(w, r, args)
	case "/open":
		h.handleOpenDeck(w, r, args)
	case "/draw":
		h.handleDrawCard(w, r, args)
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func (h *handler) handleCreateDeck(w http.ResponseWriter, r *http.Request, args url.Values) {
	var deck *models.Deck
	var shuffled = false
	var err error

	if shuffledArg, ok := args["shuffled"]; ok && len(shuffledArg) == 1 {
		shuffledVal := strings.ToLower(shuffledArg[0])
		if shuffledVal == "t" || shuffledVal == "true" || shuffledVal == "1" {
			shuffled = true
		}
	}

	if cardCodesArg, ok := args["cards"]; ok && len(cardCodesArg) > 0 {
		var cardCodes = strings.Split(cardCodesArg[0], ",")
		var cards = make([]models.Card, len(cardCodes))
		for i := range cardCodes {
			if !models.IsCodeValid(cardCodes[i]) {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			cards[i] = models.Card{Code: cardCodes[i]}
		}
		deck = models.NewDeck(shuffled, cards...)
	} else {
		deck = models.NewDeck(shuffled)
	}

	err = h.repo.Set(*deck)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	deck.Cards = nil
	jsonResponse, err := json.Marshal(deck)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = w.Write(jsonResponse)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *handler) handleOpenDeck(w http.ResponseWriter, r *http.Request, args url.Values) {
	if deckIdArg, ok := args["deck_id"]; !ok || len(deckIdArg) != 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	deckId, err := uuid.Parse(args["deck_id"][0])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	deck, err := h.repo.Get(deckId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	jsonResponse, err := json.Marshal(deck)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = w.Write(jsonResponse)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *handler) handleDrawCard(w http.ResponseWriter, r *http.Request, args url.Values) {
	if deckIdArg, ok := args["deck_id"]; !ok || len(deckIdArg) != 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if drawNumArg, ok := args["num"]; !ok || len(drawNumArg) != 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	deckId, err := uuid.Parse(args["deck_id"][0])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	drawCardsNum, err := strconv.ParseInt(args["num"][0], 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	deck, err := h.repo.Get(deckId)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	cards := deck.DrawCards(int(drawCardsNum))
	err = h.repo.Set(*deck)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonResponse, err := json.Marshal(cards)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = w.Write(jsonResponse)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
