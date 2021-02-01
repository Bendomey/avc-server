module github.com/Bendomey/avc-server

// +heroku goVersion go1.15
go 1.15

require (
	github.com/gin-gonic/gin v1.6.3
	github.com/go-redis/redis/v8 v8.4.11
	github.com/graphql-go/graphql v0.7.9
	github.com/graphql-go/graphql-go-handler v0.2.3
	github.com/heroku/x v0.0.26
	github.com/joho/godotenv v1.3.0
	github.com/sirupsen/logrus v1.7.0
	gorm.io/driver/postgres v1.0.8
	gorm.io/gorm v1.20.12
)
