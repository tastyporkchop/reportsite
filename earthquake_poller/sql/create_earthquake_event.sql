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
