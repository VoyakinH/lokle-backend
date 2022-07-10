create extension if not exists citext;

create table if not exists parents
(
    id             serial primary key,
    first_name     varchar(32) not null,
    second_name    varchar(32) not null,
    last_name      varchar(32) not null,
    phone          varchar(16) not null,
    email          citext unique,
    email_verified boolean default false,
    password       varchar(64) not null,
    access_enabled boolean default false,
    passport       varchar(8),
    dir_path       varchar(128)
);

create table if not exists children
(
    id                    serial primary key,
    first_name            varchar(32) not null,
    second_name           varchar(32) not null,
    last_name             varchar(32) not null,
    phone                 varchar(16) not null,
    email                 citext unique,
    password              varchar(64) not null,
    birth_date            timestamp not null,
    access_enabled        boolean default false,
    passport              varchar(8),
    place_of_residence    varchar(128) not null,
    place_of_registration varchar(128) not null,
    dir_path              varchar(128)
);

create table if not exists parents_children
(
    id       serial primary key,
    parent_id int not null,
    constraint fk_pc_parents foreign key (parent_id) references parents (id),
    child_id  int not null,
    constraint fk_pc_children foreign key (child_id) references children (id),
    -- мб через флаги?
    relationship varchar(32) not null, 
    unique(parent_id, child_id)
);
