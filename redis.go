package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/garyburd/redigo/redis"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

const (
	// Channel name to use with redis
	Channel = "redis"
)

var (
	waitingMessage, availableMessage []byte
	waitSleep                        = time.Second * 10
)

func init() {
	var err error
	waitingMessage, err = json.Marshal(message{
		Handle: "system",
		Text:   "Waiting for redis to be available. Messaging won't work until redis is available",
	})
	if err != nil {
		panic(err)
	}
	availableMessage, err = json.Marshal(message{
		Handle: "system",
		Text:   "Redis is now available & messaging is now possible",
	})
	if err != nil {
		panic(err)
	}
}

// redisWriter publishes messages to the Redis CHANNEL
type redisWriter struct {
	pool     *redis.Pool
	messages chan []byte
}

func newRedisWriter(pool *redis.Pool) redisWriter {
	return redisWriter{
		pool:     pool,
		messages: make(chan []byte, 10000),
	}
}

// run the main redisWriter loop that publishes incoming messages to Redis.
func (rw *redisWriter) run() error {
	conn := rw.pool.Get()
	defer conn.Close()

	for data := range rw.messages {
		if err := writeToRedis(conn, data); err != nil {
			rw.publish(data) // attempt to redeliver later
			return err
		}
	}
	return nil
}

func writeToRedis(conn redis.Conn, data []byte) error {
	if err := conn.Send("PUBLISH", Channel, data); err != nil {
		return errors.Wrap(err, "Unable to publish message to Redis")
	}
	if err := conn.Flush(); err != nil {
		return errors.Wrap(err, "Unable to flush published message to Redis")
	}
	return nil
}

// publish to Redis via channel.
func (rw *redisWriter) publish(data []byte) {
	rw.messages <- data
}
