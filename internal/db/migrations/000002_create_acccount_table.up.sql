CREATE TABLE IF NOT EXISTS account
(
	iid SERIAL
		CONSTRAINT account_pk
			PRIMARY KEY,
	created_at 			TIMESTAMPTZ NOT NULL DEFAULT now(),
	email 				VARCHAR(250),
	type 				VARCHAR(10)[] NOT NULL,
	summary 			TEXT,
	area_of_expertise 	TEXT[],
	certifications 		TEXT[]
);
