package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/lucsky/cuid"
	"github.com/redis/go-redis/v9"
)

func main() {
	SERVER_PORT := os.Getenv("SERVER_PORT")
	REDIS_ADDR := os.Getenv("REDIS_ADDR")
	PASSWORD := os.Getenv("PASSWORD")
	SESSION_TTL, err := strconv.Atoi(os.Getenv("SESSION_TTL"))
	if err != nil {
		log.Printf("error while parsing $SESSION_TTL may not be a number defaulting to 30 days: %s", err)
		SESSION_TTL = 30 * 24 * 60 * 60 // 30 days in seconds
	}

	log.Printf("SERVER_PORT=%s", SERVER_PORT)
	log.Printf("REDIS_ADDR=%s", REDIS_ADDR)
	log.Printf("PASSWORD=%s", PASSWORD)
	log.Printf("SESSION_TTL=%d", SESSION_TTL)

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

	app.Use(logger.New(logger.Config{
		TimeZone: "Asia/Seoul",
		TimeFormat: "2006-01-02T15:04:05.999999-07:00",
	}))
	app.Static("/", "./public")

	app.Get("/check", func (c *fiber.Ctx) error {
		sessionId := c.Cookies("session_id")
		if sessionId == "" {
			return c.Status(401).Send(nil)
		}
		
		exists, err := rdb.Exists(redisCtx, fmt.Sprintf("sessions:%s", sessionId)).Result()
		if err != nil {
			log.Printf("error while rdb.Exists: %s\n", err)
			return c.Status(500).Send(nil)
		}
		if exists == 1 {
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
		if password != PASSWORD {
			return c.Status(401).Send(nil)
		}

		sessionId := cuid.New()

		// put session id to redis
		ttl := time.Duration(SESSION_TTL) * time.Second
		_, err := rdb.Set(redisCtx, fmt.Sprintf("sessions:%s", sessionId), 1, ttl).Result()
		if err != nil {
			log.Printf("error while rdb.Exists: %s\n", err)
			return c.Status(500).Send(nil)
		}

		// set cooke of .deps.me
		cookie := fiber.Cookie {
			Name: "session_id",
			Value: sessionId,
			Domain: "deps.me",
			Expires: time.Now().Add(ttl),
		}
		c.Cookie(&cookie)
		return c.Status(200).Send(nil)
	})

	log.Panic(app.Listen(fmt.Sprintf(":%s", SERVER_PORT)))
}
