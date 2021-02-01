package main

import (
	"github.com/Bendomey/avc-server/internal/redis"
	"github.com/Bendomey/avc-server/pkg/server"
)

func main() {
	//connects to redis
	redis.Factory()

	// server invoked here
	server.Run()
}
