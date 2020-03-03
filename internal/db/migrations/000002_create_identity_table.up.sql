CREATE TABLE IF NOT EXISTS identity
(
	iid SERIAL
		CONSTRAINT identities_pk
			PRIMARY KEY,
	created_at TIMESTAMPTZ DEFAULT now(),
	PROVIDER VARCHAR(60) NOT NULL,
	uid VARCHAR(250) NOT NULL,
	pw_hash BYTEA,
	confirmed_at TIMESTAMPTZ,
    account_iid INTEGER REFERENCES account(iid)
);
