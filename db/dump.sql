create extension if not exists citext;

create table if not exists users
(
    id             serial primary key,
    role           smallint not null,
    first_name     varchar(32) not null,
    second_name    varchar(32) not null,
    last_name      varchar(32) not null,
    phone          varchar(16) not null,
    email          citext unique,
    email_verified boolean default false,
    password       varchar(64) not null
);

create table if not exists parents
(
    id                     serial primary key,
    user_id                int not null,
    constraint fk_users_parents foreign key (user_id) references users (id),
    passport               varchar(8),
    passport_verified      boolean default false,
    passport_uploaded_time timestamptz,
    dir_path               varchar(128)
);

create table if not exists children
(
    id                      serial primary key,
    user_id                 int not null,
    constraint fk_users_children foreign key (user_id) references users (id),
    birth_date              date not null,
    first_stage_done        boolean default false,
    first_stage_start_time  timestamptz,
    second_stage_done       boolean default false,
    second_stage_start_time timestamptz,
    third_stage_done        boolean default false,
    third_stage_start_time  timestamptz,
    passport                varchar(8),
    place_of_residence      varchar(128) not null,
    place_of_registration   varchar(128) not null,
    dir_path                varchar(128),
    remarks                 varchar(1024)
);

create table if not exists parents_children
(
    id           serial primary key,
    parent_id    int not null,
    constraint fk_pc_parents foreign key (parent_id) references parents (id),
    child_id     int not null,
    constraint fk_pc_children foreign key (child_id) references children (id),
    -- мб через флаги?
    relationship smallint not null, 
    unique(parent_id, child_id)
);

drop table parents_children;
drop table children;
drop table parents;

