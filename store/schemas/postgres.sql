create table tuples(
	u text not null,
	v text not null,
	label text not null,

	primary key(u, label, v)
);

