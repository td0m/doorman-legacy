create table tuples(
  subject text not null references objects(id),
  role text not null references roles(id),
  object text not null references objects(id),

  primary key(subject, role, object)
);

-- already indexed for listing connections (from primary key), but need to support the same in reverse
create index "tuples_idx_reverse_lookup" on tuples(object, role);

create table changes(
  id text primary key,
  type text not null,
  payload jsonb not null,
  created_at timestamptz not null default now()
);

create table roles(
  id text primary key,
  verbs text[] not null default '{}'
);
