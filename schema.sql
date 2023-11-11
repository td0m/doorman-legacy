-- drop table changes;
-- drop table tuples;
-- drop table roles;
-- drop table objects;

-- create table objects(
--   id text primary key,
--   key text unique
-- );

create table roles(
  id text primary key,
  verbs text[] not null default '{}'
);

create table tuples(
  subject text not null,
  role text not null references roles(id),
  object text not null,

  primary key(subject, role, object)
);

-- already indexed for listing connections (from primary key), but need to support the same in reverse
create index "tuples_idx_reverse_lookup" on tuples(object, role);

create table changes(
  id text primary key,
  type text not null,
  payload jsonb not null,
  status text not null default 'pending',
  created_at timestamptz not null default now()
);
