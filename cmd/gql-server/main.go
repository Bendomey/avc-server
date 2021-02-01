package main

import (
	log "github.com/Bendomey/avc-server/internal/logger"
	"github.com/Bendomey/avc-server/internal/orm"
	"github.com/Bendomey/avc-server/internal/redis"
	"github.com/Bendomey/avc-server/pkg/server"
)

func main() {
	//connects to redis
	redis.Factory()

	// creates a new ORM instance to send it to our server
	_, err := orm.Factory()
	if err != nil {
		log.Panic("[ORM ERR] :: ", err)
	}

	// server invoked here
	server.Run()
}
