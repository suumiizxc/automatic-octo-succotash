create table account (
	ID serial primary key,
	password varchar(127),
	user_name varchar(255),
	first_name varchar(63),
	last_name varchar(63),
	email varchar(255),
	is_active int,
	avatar_image varchar(255),
	role int,
	created_at timestamp ,
	updated_at timestamp
);

