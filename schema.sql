create table tuples(
  subject text not null references objects(id),
  role text not null references roles(id),
  object text not null references objects(id),

  primary key(subject, role, object)
);

-- already indexed for listing connections (from primary key), but need to support the same in reverse
create index "tuples_idx_reverse_lookup" on tuples(object, role);

create unlogged table relations(
  subject text not null,
  verb text not null,
  object text not null,

  primary key (subject, verb, object),
  constraint "relations_fkey_subject" foreign key (subject) references objects(id) on delete cascade,
  constraint "relations_fkey_object" foreign key (object) references objects(id) on delete cascade
);

create table changes(
  id text primary key,
  type text not null,
  payload jsonb not null,
  created_at timestamptz not null default now()
);

