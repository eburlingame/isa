package main

import (
	"context"
	"log"
	"os"
	"time"

	"math/rand"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"gopkg.in/olahol/melody.v1"
)

var ctx = context.Background()

func main() {
	rand.Seed(time.Now().Unix())

	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST"),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	r := gin.Default()
	m := melody.New()

	r.GET("/", func(c *gin.Context) {
		http.ServeFile(c.Writer, c.Request, "index.html")
	})

	r.GET("/ws", func(c *gin.Context) {
		m.HandleRequest(c.Writer, c.Request)
	})

	m.HandleMessage(func(s *melody.Session, msg []byte) {
		// fmt.Printf("Got msg: %s\n", msg)
		DispatchMessage(&ctx, rdb, s, msg)
	})

	m.HandleDisconnect(func(s *melody.Session) {
		log.Printf("Session disconnected\n")

		persistentSession, err := GetPersistentSession(s)
		if err != nil {
			return
		}

		// Close the goroutine that is listening for state changes
		persistentSession.UnsubChan <- true
	})

	r.Run(":" + os.Getenv("API_SERVER_PORT"))
}
