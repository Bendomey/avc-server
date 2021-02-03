package services

import (
	"github.com/Bendomey/avc-server/internal/orm"
	"github.com/go-redis/redis/v8"
)

// Services responsible for exposing all services to resolvers
type Services struct {
	AdminServices AdminService
}

//Factory activates all services
func Factory(orm *orm.ORM, rdb *redis.Client) Services {

	//activate admin service
	adminService := NewAdminSvc(orm, rdb)

	return Services{
		AdminServices: adminService,
	}
}
