package model

import (
	"time"
)

type BaseModel struct {
	IID       int       `db:"iid"`
	CreatedAt time.Time `db:"created_at"`
}
