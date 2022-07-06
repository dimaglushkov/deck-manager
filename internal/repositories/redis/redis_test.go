package redis

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/dimaglushkov/deck-manager/internal/models"
)

const host, port = "127.0.0.1", "6379"

var cache *Cache

func TestCache_New(t *testing.T) {
	var err error
	_, err = New(host, "")
	assert.ErrorContains(t, err, "connection")

	_, err = New("", "")
	assert.ErrorContains(t, err, "connection")

	cache, err = New(host, port)
	require.NoErrorf(t, err, "To run these test start redis instance at %s:%s", host, port)
}

func TestCache(t *testing.T) {
	var err error
	assert.NotNil(t, cache)
	deck := models.NewDeck(false)

	err = cache.Set(*deck)
	assert.NoError(t, err)

	savedDeck, err := cache.Get(deck.UUID)
	assert.NoError(t, err)
	assert.Equal(t, deck, savedDeck)

	deck.DrawCards(20)
	err = cache.Set(*deck)
	assert.NoError(t, err)

	savedDeck, err = cache.Get(deck.UUID)
	assert.NoError(t, err)
	assert.Equal(t, deck, savedDeck)

	_, err = cache.Get(uuid.New())
	assert.ErrorContains(t, err, "uuid")
}
