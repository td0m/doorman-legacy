create table relations(
  id text primary key,
  name text,
  "from" text not null,
  "to" text not null
);

create table dependencies(
  relation_id text not null references relations(id),
  depends_on text not null references relations(id),
  primary key(relation_id, depends_on)
);

create table schema(
  value jsonb not null default '{}'
);

