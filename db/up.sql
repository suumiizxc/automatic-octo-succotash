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

create table content (
	id serial primary key,
	name varchar(255),
	published_year int,
	page_count int,
	content_scope text,
	is_purchase int,
	is_rent int,
	rent_price numeric(10,2),
	author varchar(255),
	content_purpose text,
	brief_description text,
	summary text,
	description text,
	table_list text,
	image_list text,
	cover varchar(127),
	feature_start_date timestamp,
	feature_end_date timestamp,
	is_featured int,
	feature_image varchar(127),
	preview_file varchar(127),
	language varchar(255),
	rent_month int,
	is_free int, 
	is_published int,
	created_at timestamp,
	updated_at timestamp,
	created_user_id int,
	updated_user_id int
);

create table content_author
(
	id serial primary key,
	author_type varchar(255),
	name varchar(255),
	image varchar(127),
	description text,
	created_at timestamp,
	updated_at timestamp,
	created_user_id int,
	updated_user_id int
);

create table content_file 
(
	id serial primary key,
	content_id int,
	attachment varchar(255),
	file_type varchar(100),
	file_url varchar(255),
	source_url varchar(255),
	created_at timestamp,
	updated_at timestamp,
	created_user_id int,
	updated_user_id int
);

create table blog 
(
	id serial primary key,
	language varchar(127),
	title varchar(511),
	poster varchar(127),
	cover varchar(127),
	brief_description text,
	description text,
	category_id int,
	is_published int,
	is_featured int,
	feature_start_date timestamp,
	feature_end_date timestamp,
	created_at timestamp,
	updated_at timestamp,
	created_user_id int,
	updated_user_id int
);

create table blog_category 
(
	id serial primary key,
	name_mn varchar(255),
	name_en varchar(255),
	is_featured int, 
	created_at timestamp,
	updated_at timestamp,
	created_user_id int,
	updated_user_id int
);