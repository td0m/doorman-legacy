create table entities(
  _id text not null,
  _type text not null,
  attrs jsonb,

  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),

  primary key(_type, _id)
);

create index "entities.idx-id-type" on entities(_id, _type);

-- directed graph
-- todo: consider another table for all "caches"
create table relations(
  _id text primary key,
  name text not null,
  from_id text not null,
  from_type text not null,
  to_id text not null,
  to_type text not null,
  indirect bool not null,

  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),

  constraint "relations.fkey-from" foreign key (from_type, from_id) references entities(_type, _id),
  constraint "relations.fkey-to" foreign key (to_type, to_id) references entities(_type, _id)
);

create index "relations.idx-from-to"
on relations(from_type, from_id, to_type, to_id);

create index "relations.idx-from"
on relations(from_type, from_id);

create index "relations.idx-to"
on relations(to_type, to_id);

create table dependencies(
  relation_id text not null,
  dependency_id text not null,

  primary key(relation_id, dependency_id),

  constraint "dependencies.fkey-relation_id" foreign key (relation_id) references relations(_id) on delete cascade,
  constraint "dependencies.fkey-dependency_id" foreign key (dependency_id) references relations(_id)
);

-- removes relations that depend on the one being deleted
create or replace function remove_dependents()
returns trigger
language plpgsql
as $$
begin
  delete from relations
  where _id in (
    select relation_id from dependencies where dependency_id=old._id
  );
  return old;
end;
$$;

-- delete dependent relations before deleting a given relation
create trigger trg_delete_dependencies
before delete on relations
for each row
when (old.indirect = false) -- prevents infinite loop, as we cannot depend on "indirect" relations.
execute procedure remove_dependents();

