package model

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type Identity struct {
	BaseModel

	Provider    IdentityType `db:"provider"`
	UID         string       `db:"uid"`
	PWHash      []byte       `db:"pw_hash"`
	ConfirmedAt sql.NullTime `db:"confirmed_at"`
	AccountIID  int          `db:"account_iid"`
}

type IdentityType string

const (
	IdentityTypeEmail     IdentityType = "email"
	IdentityTypeGoogle    IdentityType = "google"
	IdentityTypeMicrosoft IdentityType = "microsoft"
)

func NewIdentity(t IdentityType, email string, accountIID int) Identity {
	return Identity{
		Provider:   t,
		UID:        email,
		AccountIID: accountIID,
	}
}

func GetIdentity(db *sqlx.DB, iType IdentityType, uid string) (Identity, error) {
	q := `SELECT *
		  FROM identity
		  WHERE provider=$1 AND uid=$2`

	var i Identity
	err := db.Get(&i, q, iType, uid)
	return i, err
}

func (i *Identity) Insert(tx *sqlx.Tx) error {
	q := `INSERT INTO identity(provider, uid, pw_hash, account_iid, confirmed_at)
		  VALUES(:provider, :uid, :pw_hash, :account_iid, :confirmed_at)`
	_, err := tx.NamedExec(q, i)
	return err
}
