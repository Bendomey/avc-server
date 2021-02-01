package redis

import (
	"fmt"

	log "github.com/Bendomey/avc-server/internal/logger"
	"github.com/Bendomey/avc-server/pkg/utils"
	"github.com/go-redis/redis/v8"
)

var host, port, password string

func init() {
	host = utils.MustGet("REDIS_HOST")
	port = utils.MustGet("REDIS_PORT")
	password = utils.MustGet("REDIS_PASSWORD")
}

// Factory connects app to reid
func Factory() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: password, // password set
		DB:       0,        // use default DB
	})
	log.Info("Redis :: Connected successfully")
	return rdb
}
