package main

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/gommon/log"
	"gopkg.in/olahol/melody.v1"
)

func WatchGame(ctx *context.Context, rdb *redis.Client, session *melody.Session, gameId string, done <-chan bool) {
	pubsub := rdb.Subscribe(*ctx, "game:"+gameId)

	ch := pubsub.Channel()

	for {
		select {
		case _ = <-ch:
			game, err := LoadGame(ctx, rdb, gameId)

			if err != nil {
				log.Error("Unable to load game")
				return
			}

			SendGameUpdate(session, gameId, game, game.State == GameAbandoned)
			break

		case <-done:
			fmt.Println("Unsubscribing from game")
			pubsub.Close()
			return
		}
	}

}
