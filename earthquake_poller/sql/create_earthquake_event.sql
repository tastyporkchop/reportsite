create table earthquake_event (
	id serial,
	event_id varchar,
	title varchar,
	updated timestamp,
	link varchar,
	summary varchar,
	summary_type varchar,
	point varchar,
	elevation varchar,
	primary key(id)
)
;

create index earthquake_event_event_id_idx on earthquake_event(event_id);
