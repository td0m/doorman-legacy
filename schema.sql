-- drop table objects;
-- drop table roles;
-- drop table tuples;
-- drop index "tuples_reverse_lookup";

create table objects(
  id text primary key
);

create table roles(
  id text primary key,
  verbs text[] not null default '{}'
);

create table tuples(
  subject text not null references objects(id),
  role text not null references roles(id),
  object text not null references objects(id),

  primary key(subject, role, object)
);

-- already indexed for listing connections (from primary key), but need to support the same in reverse
create index "tuples_reverse_lookup" on tuples(object, role);

