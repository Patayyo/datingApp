create table users (
    id serial primary key,
    username varchar(255) not null,
    email varchar(255) not null,
    password varchar(255) not null,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp
);

create table profiles (
    id serial primary key,
    user_id int references users(id) on delete cascade,
    bio text,
    interests text,
    age varchar(255) not null,
    gender varchar(255) not null,
    location varchar(255) not null,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp
);

create table matches (
    id serial primary key,
    user_id int references users(id) on delete cascade,
    matched_id int references users(id) on delete cascade,
    matched boolean not null,
    created_at timestamp not null default current_timestamp,
    updated_at timestamp not null default current_timestamp
);
