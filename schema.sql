create table entities(
  _id text primary key,
	attrs jsonb not null
	-- TODO: unique on json "sub" if _id starts with users/
);

-- directed graph
create table relations(
  _id text primary key,
  "from" text not null references entities(_id),
	"to" text not null references entities(_id),
	attrs jsonb not null
);

