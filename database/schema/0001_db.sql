USE axolobot;
alter table mention
	add is_done boolean default false not null;
alter table mention
	add time timestamp default now() not null;
