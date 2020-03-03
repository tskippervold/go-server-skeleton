package env

import (
	"github.com/jmoiron/sqlx"
	"github.com/tskippervold/golang-base-server/internal/utils/log"
)

type Env struct {
	DB  *sqlx.DB
	Log *log.Log
}
