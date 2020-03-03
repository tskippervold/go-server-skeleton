package model

import (
	"github.com/jmoiron/sqlx"
)

type Account struct {
	BaseModel
	Email string `db:"email"`
}

func NewAccount(email string) Account {
	return Account{
		Email: email,
	}
}

func AccountExists(db *sqlx.DB, email string) (bool, error) {
	var c int
	err := db.Get(&c, "SELECT count(*) FROM account WHERE email=$1", email)
	return c > 0, err
}

func GetAccount(db *sqlx.DB, email string) (Account, error) {
	var a Account
	err := db.Get(&a, "SELECT * FROM account WHERE email=$1", email)
	return a, err
}

func (a *Account) Insert(tx *sqlx.Tx) (int, error) {
	q := `INSERT INTO account(email)
		  VALUES(:email)
		  RETURNING iid`

	stmt, err := tx.PrepareNamed(q)
	if err != nil {
		return -1, err
	}

	var iid int
	if err = stmt.Get(&iid, a); err != nil {
		return -1, err
	}

	return iid, err
}
