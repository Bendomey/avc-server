module github.com/Bendomey/avc-server

// +heroku goVersion go1.15
go 1.15

require (
	github.com/Bendomey/goutilities v0.0.0-20201104205146-d5b8f238bf1b
	github.com/certifi/gocertifi v0.0.0-20200922220541-2c3bb06c6054 // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/getsentry/raven-go v0.2.0
	github.com/getsentry/sentry-go v0.9.0
	github.com/gin-gonic/gin v1.6.3
	github.com/go-gormigrate/gormigrate/v2 v2.0.0
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/go-redis/redis/v8 v8.4.11
	github.com/gobuffalo/envy v1.9.0 // indirect
	github.com/gofrs/uuid v4.0.0+incompatible
	github.com/graphql-go/graphql v0.7.9
	github.com/graphql-go/graphql-go-handler v0.2.3
	github.com/heroku/x v0.0.26
	github.com/joho/godotenv v1.3.0
	github.com/kehindesalaam/go-paystack v0.0.0-20171103231913-9286b18021f6
	github.com/kehindesalaam/mapstructure v0.0.0-20170924220807-6b3baa82ce04 // indirect
	github.com/mailgun/mailgun-go v2.0.0+incompatible
	github.com/mailgun/mailgun-go/v4 v4.3.3
	github.com/mitchellh/mapstructure v1.4.1
	github.com/sirupsen/logrus v1.7.0
	github.com/vektah/gqlparser v1.3.1 // indirect
	gorm.io/driver/postgres v1.0.8
	gorm.io/gorm v1.20.12
)
