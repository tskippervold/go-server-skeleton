create table if not exists company
(
	iid serial
		constraint company_pk
			primary key,
	created_at timestamptz default now(),
	brreg_iid integer references brreg(iid),
    created_by_account_iid integer references account(iid)
);
