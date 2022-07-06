package redis

import (
	"encoding/json"
	"fmt"

	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"

	. "github.com/dimaglushkov/toggl-test-assignment/internal"
	"github.com/dimaglushkov/toggl-test-assignment/internal/models"
)

type Cache struct {
	connPool *redis.Pool
}

func New(redisHost, redisPort string) (*Cache, error) {
	redisUrl := fmt.Sprintf("redis://%s:%s", redisHost, redisPort)
	connPool := redis.Pool{
		MaxIdle:   60,
		MaxActive: 1000,
		Dial: func() (redis.Conn, error) {
			return redis.DialURL(redisUrl)
		},
	}

	conn := connPool.Get()
	defer conn.Close()
	_, err := conn.Do("PING")
	if err != nil {
		return nil, fmt.Errorf("error while creating redis cache: %v", err)
	}

	return &Cache{connPool: &connPool}, nil
}

func (c *Cache) Set(deck models.Deck) error {
	conn := c.connPool.Get()
	defer conn.Close()

	uuidString := deck.UUID.String()
	deckDataJSON, err := json.Marshal(deck)
	if err != nil {
		return err
	}

	_, err = conn.Do("SET", uuidString, deckDataJSON)
	return err
}

func (c *Cache) Get(uuid uuid.UUID) (*models.Deck, error) {
	conn := c.connPool.Get()
	defer conn.Close()

	uuidString := uuid.String()
	deckDataJSON, err := conn.Do("GET", uuidString)
	if err != nil {
		return nil, err
	}

	if deckDataJSON == nil || len(deckDataJSON.([]byte)) == 0 {
		return nil, NewUnknownUUIDError(uuid)
	}

	var deck models.Deck
	err = json.Unmarshal(deckDataJSON.([]byte), &deck)
	if err != nil {
		return nil, err
	}
	return &deck, nil
}
