package main

import (
	log "github.com/Bendomey/avc-server/internal/logger"
	"github.com/Bendomey/avc-server/internal/orm"
	"github.com/Bendomey/avc-server/internal/redis"
	"github.com/Bendomey/avc-server/internal/services"
	"github.com/Bendomey/avc-server/pkg/server"
)

func main() {
	//connects to redis
	rdb := redis.Factory()

	// creates a new ORM instance to send it to our server
	orm, err := orm.Factory()
	if err != nil {
		log.Panic("[ORM ERR] :: ", err)
	}

	//start services here
	services := services.Factory(orm, rdb)

	// server invoked here
	server.Run(services)
}
