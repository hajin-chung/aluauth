package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/lucsky/cuid"
	"github.com/redis/go-redis/v9"
)

func main() {
	SERVER_PORT := os.Getenv("SERVER_PORT")
	REDIS_ADDR := os.Getenv("REDIS_ADDR")
	PASSWORD := os.Getenv("PASSWORD")

	log.Printf("SERVER_PORT=%s", SERVER_PORT)
	log.Printf("REDIS_ADDR=%s", REDIS_ADDR)
	log.Printf("PASSWORD=%s", PASSWORD)

	if REDIS_ADDR == "" {
		REDIS_ADDR = "127.0.0.1:6379"
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: REDIS_ADDR,
		Password: "",
		DB: 0,
	})
	redisCtx := context.Background()

	app := fiber.New()

	app.Static("/", "./public")

	app.Get("/check", func (c *fiber.Ctx) error {
		sessionId := c.Cookies("session_id")
		if sessionId == "" {
			return c.Status(401).Send(nil)
		}
		
		exists, _ := rdb.SIsMember(redisCtx, "sessions", sessionId).Result()
		if exists {
			return c.Status(200).Send(nil)
		} else {
			return c.Status(401).Send(nil)
		}
	})

	app.Get("/login", func (c *fiber.Ctx) error {
		return c.SendFile("./public/login.html")
	})

	app.Post("/login", func (c *fiber.Ctx) error {
		password := string(c.Body()[:])
		if password == PASSWORD {
			sessionId := cuid.New()

			// put session id to redis
			_, err := rdb.SAdd(redisCtx, "sessions", sessionId).Result()
			if err != nil {
				return c.Status(500).Send(nil)
			}

			// set cooke of .deps.me
			cookie := fiber.Cookie {
				Name: "session_id",
				Value: sessionId,
				Domain: "deps.me",
			}
			c.Cookie(&cookie)
			return c.Status(200).Send(nil)
		}
		// redirect to login
		return c.Status(401).Send(nil)
	})

	app.Listen(fmt.Sprintf(":%s", SERVER_PORT))
}
