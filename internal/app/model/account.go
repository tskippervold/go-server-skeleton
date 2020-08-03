package model

import (
	"database/sql"

	"github.com/lib/pq"

	"github.com/go-playground/validator"
	"github.com/jmoiron/sqlx"
)

type Account struct {
	BaseModel
	Email           string         `db:"email" validate:"required,email" json:"email"`
	Type            pq.StringArray `db:"type" validate:"dive,eq=regular|eq=consultant"`
	Summary         sql.NullString `db:"summary"`
	AreaOfExpertise pq.StringArray `db:"area_of_expertise"`
	Certifications  pq.StringArray `db:"certifications"`
	CompanyIID      sql.NullInt64  `db:"company_iid"`
}

type AccountType string

const (
	AccountTypeRegular    AccountType = "regular"
	AccountTypeConsultant AccountType = "consultant"
)

func NewAccount(email string, t AccountType) Account {
	return Account{
		Email: email,
		Type:  []string{string(t)},
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

func (a *Account) Validate() error {
	v := validator.New()
	return v.Struct(a)
}

func (a *Account) Insert(tx *sqlx.Tx) (int, error) {
	q := `INSERT INTO account(email, type)
		  VALUES(:email, :type)
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

func (a *Account) Update(db *sqlx.DB) error {
	q := `UPDATE account
		  SET
		  	summary=:summary,
			area_of_expertise=:areaOfExpertise,
			certifications=:certifications,
			type=:type
		  WHERE iid=:iid`

	_, err := db.NamedExec(q, map[string]interface{}{
		"iid":             a.IID,
		"summary":         a.Summary,
		"areaOfExpertise": a.AreaOfExpertise,
		"certifications":  a.Certifications,
		"type":            a.Type,
	})

	return err
}
