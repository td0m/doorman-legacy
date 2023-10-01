create table entities(
  _id text not null,
  _type text not null,
  attrs jsonb not null,

  primary key(_type, _id)
  -- TODO: unique on json "sub" if _id starts with users/
);

create index "entities.idx-id-type" on entities(_id, _type);

-- directed graph
create table relations(
  _id text primary key,
  attrs jsonb not null,
  from_id text not null,
  from_type text not null,
  to_id text not null,
  to_type text not null,

  constraint "relations.fkey-from" foreign key (from_type, from_id) references entities(_type, _id),
  constraint "relations.fkey-to" foreign key (to_type, to_id) references entities(_type, _id)
);

create index "relations.idx-from-to"
on relations(from_type, from_id, to_type, to_id);
