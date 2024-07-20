package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

func GameExists(ctx *context.Context, rdb *redis.Client, gameId string) bool {
	val, err := rdb.Get(*ctx, "game:"+gameId).Result()

	if err == redis.Nil || err != nil || val == "" {
		return false
	}

	return true
}

func SaveGame(ctx *context.Context, rdb *redis.Client, gameId string, game *Game) {
	stored, _ := json.Marshal(game)
	rdb.Set(*ctx, "game:"+gameId, []byte(stored), 12*time.Hour)

	// Publish an event to the game:gameId topic to notify other players
	err := rdb.Publish(*ctx, "game:"+gameId, "updated").Err()
	if err != nil {
		log.Println(err)
	}
}

func LoadGame(ctx *context.Context, rdb *redis.Client, gameId string) (*Game, error) {
	var game Game

	stored, err := rdb.Get(*ctx, "game:"+gameId).Result()
	if err != nil || stored == "" {
		return nil, errors.New("Unable to find game")
	}

	err = json.Unmarshal([]byte(stored), &game)
	if err != nil {
		return nil, errors.New("Unable to unmarshal game")
	}

	return &game, nil
}

func DeleteGame(ctx *context.Context, rdb *redis.Client, gameId string) error {
	err := rdb.Del(*ctx, "game:"+gameId).Err()
	if err != nil {
		return err
	}

	return nil
}
