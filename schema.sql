create table entities(
  id text primary key,
  attrs jsonb,

  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);

create table relations(
  id text primary key,
  name text,
  "from" text not null,
  "to" text not null,

  constraint "relations.fkey-from" foreign key ("from") references entities(id) on delete cascade,
  constraint "relations.fkey-to" foreign key ("to") references entities(id) on delete cascade
);

-- for computing caches
create index "relations.idx-from" on relations("from");
create index "relations.idx-to" on relations("to");

create unlogged table cache(
  id text primary key,
  "from" text not null,
  from_type text not null generated always as (split_part("from", ':', 1)) stored,
  name text,
  "to" text not null,
  to_type text not null generated always as (split_part("to", ':', 1)) stored,

  constraint "cache.fkey-from" foreign key ("from") references entities(id) on delete cascade,
  constraint "cache.fkey-to" foreign key ("to") references entities(id) on delete cascade
);

create index "cache.idx-from-to" on cache("from", "to");
-- "list accessible"
-- list accessible by type? do we want that?
-- e.g. list all posts I can access.
-- computed col?
create index "cache.idx-from" on cache("from");
create index "cache.idx-to" on cache("to");

-- listing access or listing things with access can be done efficiently by type
-- todo: consider if an index is really needed here
create index "cache.idx-from-to_type" on cache("from", to_type);
create index "cache.idx-to-from_type" on cache("to", from_type);
