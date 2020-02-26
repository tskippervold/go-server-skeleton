package env

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/tskippervold/golang-base-server/internal/utils"
)

type Env struct {
	DB  *gorm.DB
	Log *utils.Log
}

func ConnectDatabase(host string, port string, name string, user string, pass string) (*gorm.DB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		host, port, name, user, pass)

	db, err := gorm.Open("postgres", connStr)

	if db != nil {
		db.SingularTable(true)
	}

	return db, err
}

func (e *Env) MigrateDatabase() error {
	e.Log.Info("Migrating database")
	return nil
	// return e.DB.AutoMigrate(&models.Account{}, &models.Organization{}, &models.OTP{}).Error
}
