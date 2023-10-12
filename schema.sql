create table entity_types(
  _id text primary key,
  attrs jsonb,

  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);

insert into entity_types(_id, attrs) values('collection', '{}');
insert into entity_types(_id, attrs) values('user', '{}');
insert into entity_types(_id, attrs) values('role', '{}');
insert into entity_types(_id, attrs) values('permission', '{}');

create table entities(
  _id text not null,
  _type text not null,
  attrs jsonb,

  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),

  primary key(_type, _id),
  constraint "entities.fkey-_type" foreign key (_type) references entity_types(_id)
);

create index "entities.idx-id-type" on entities(_id, _type);

-- directed graph
create table relations(
  _id text primary key,
  name text,
  from_id text not null,
  from_type text not null,
  to_id text not null,
  to_type text not null,

  constraint "relations.fkey-from" foreign key (from_type, from_id) references entities(_type, _id) on delete cascade,
  constraint "relations.fkey-to" foreign key (to_type, to_id) references entities(_type, _id) on delete cascade
);

create index "relations.idx-from" on relations(from_type, from_id);
create index "relations.idx-to" on relations(to_type, to_id);

-- TODO: ensure they get deleted after relation deleted
create unlogged table cache(
  _id text primary key,
  name text,
  from_id text not null,
  from_type text not null,
  to_id text not null,
  to_type text not null,

  constraint "cache.fkey-from" foreign key (from_type, from_id) references entities(_type, _id) on delete cascade,
  constraint "cache.fkey-to" foreign key (to_type, to_id) references entities(_type, _id) on delete cascade
);

create index "cache.idx-from-to" on cache(from_type, from_id, to_type, to_id);
create index "cache.idx-from" on cache(from_type, from_id);
create index "cache.idx-to" on cache(to_type, to_id);

create unlogged table dependencies(
  relation_id text not null,
  cache_id text not null,

  primary key(relation_id, cache_id),

  -- relation cannot be removed if any cached dependencies lingering
  constraint "dependencies.fkey-relation_id" foreign key (relation_id) references relations(_id) on delete cascade,
  -- cache is dropped = remove dependencies linked to it
  constraint "dependencies.fkey-cache_id" foreign key (cache_id) references cache(_id) on delete cascade
);


-- removes cached relations that depend on the one being deleted
create or replace function remove_dependent_cache()
returns trigger
language plpgsql
as $$
begin
  delete from cache
  where _id in (
    select cache_id from dependencies where relation_id=old._id
  );
  return old;
end;
$$;

-- delete dependent cache relations before deleting a given relation
create trigger trg_delete_dependent_cache
before delete on relations
for each row
execute procedure remove_dependent_cache();

-- TODO: ensure if dependency dropped then cache also dropped
