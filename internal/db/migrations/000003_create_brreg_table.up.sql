create table if not exists brreg
(
	iid serial
		constraint brreg_pk
			primary key,
	created_at timestamptz default now(),
	registered_at timestamptz not null,
	source varchar(250) not null,
	vat_num varchar(250) not null,
	name varchar(250) not null,
    number_of_employees integer not null default 0,
    org_type_code varchar(10),
    org_kind_code varchar(10),
    org_kind_description varchar(250),
    is_vat_registered boolean,
    is_bankrupt boolean,
    is_under_liquidation boolean,
    is_under_forced_liquidation boolean
);
