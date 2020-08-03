package model

import (
	"time"
)

type BaseModel struct {
	IID       int       `db:"iid" json:"-"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
}
