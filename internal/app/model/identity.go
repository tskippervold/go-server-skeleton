package model

import (
	"time"

	"github.com/jmoiron/sqlx"
)

type Identity struct {
	BaseModel

	Provider    string     `db:"provider"`
	UID         string     `db:"uid"`
	PWHash      []byte     `db:"pw_hash"`
	ConfirmedAt *time.Time `db:"confirmed_at"`
	AccountIID  int        `db:"account_iid"`
}

const (
	IdentityTypeEmail = "email"
)

func NewIdentityEmail(email string, accountIID int, pwHash []byte) Identity {
	return Identity{
		Provider:   IdentityTypeEmail,
		UID:        email,
		AccountIID: accountIID,
		PWHash:     pwHash,
	}
}

func GetIdentityEmail(db *sqlx.DB, email string) (Identity, error) {
	q := `SELECT *
		  FROM identity
		  WHERE provider=$1 AND uid=$2`

	var i Identity
	err := db.Get(&i, q, IdentityTypeEmail, email)
	return i, err
}

func (i *Identity) Insert(tx *sqlx.Tx) error {
	q := `INSERT INTO identity(provider, uid, pw_hash, account_iid)
		  VALUES(:provider, :uid, :pw_hash, :account_iid)`
	_, err := tx.NamedExec(q, i)
	return err
}
