create table relations(
  "from" text not null,
  name text,
  "to" text not null,
  via text[] not null default '{}',

  primary key("from", name, "to", via)
);

create table schema(
  value jsonb not null default '{}'
);

