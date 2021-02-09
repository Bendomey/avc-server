package orm

import (
	"fmt"

	log "github.com/Bendomey/avc-server/internal/logger"
	"github.com/Bendomey/avc-server/internal/orm/migration"
	"github.com/getsentry/raven-go"
	"gorm.io/gorm"

	"github.com/Bendomey/avc-server/pkg/utils"
	// Imports the database dialect of choice
	"gorm.io/driver/postgres"
)

var host, user, password, dbname, port, sslmode string
var autoMigrate, seedDB bool

// ORM struct to holds the gorm pointer to db
type ORM struct {
	DB *gorm.DB
}

func init() {
	host = utils.MustGet("GORM_HOST")
	user = utils.MustGet("GORM_USER")
	password = utils.MustGet("GORM_PASSWORD")
	dbname = utils.MustGet("GORM_DBNAME")
	port = utils.MustGet("GORM_PORT")
	autoMigrate = utils.MustGetBool("GORM_AUTOMIGRATE")
	seedDB = utils.MustGetBool("GORM_SEED_DB")
}

// Factory creates a db connection with the selected dialect and connection string
func Factory() (*ORM, error) {
	db, err := gorm.Open(postgres.Open(fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", host, user, password, dbname, port, sslmode)), &gorm.Config{})
	if err != nil {
		raven.CaptureError(err, nil)
		log.Panic("[ORM] err: ", err)
	}
	log.Info("[ORM] :: Database connection initialized.")
	orm := &ORM{
		DB: db,
	}
	if autoMigrate {
		err = migration.ServiceAutoMigration(orm.DB, seedDB)
		raven.CaptureError(err, nil)
	}

	return orm, err
}
