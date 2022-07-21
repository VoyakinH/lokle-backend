create extension if not exists citext;

-- auto-generated definition
create table users
(
    id             bigserial
        constraint users_pk
            primary key,
    role           smallint              not null,
    first_name     varchar(32)           not null,
    second_name    varchar(32)           not null,
    last_name      varchar(32),
    phone          varchar(16)           not null,
    email          citext                not null,
    email_verified boolean default false not null,
    password       varchar(64)           not null
);

alter table users
    owner to lokle_admin;

create unique index users_email_uindex
    on users (email);

create unique index users_id_uindex
    on users (id);



-- auto-generated definition
create table parents
(
    id                bigserial
        constraint parents_pk
            primary key,
    user_id           bigint                not null
        constraint parents_users_id_fk
            references users
            on update cascade on delete cascade,
    passport          varchar(64) default ''::character varying,
    passport_verified boolean default false not null,
    dir_path          varchar(128) default ''::character varying
);

alter table parents
    owner to lokle_admin;

create unique index parents_id_uindex
    on parents (id);

create unique index parents_user_id_uindex
    on parents (user_id);






-- auto-generated definition
create table children
(
    id                    bigserial
        constraint children_pk
            primary key,
    user_id               bigint             not null
        constraint children_users_id_fk
            references users
            on update cascade on delete cascade,
    birth_date            bigint             not null,
    done_stage            smallint default 0 not null,
    passport              varchar(64)  default ''::character varying,
    place_of_residence    varchar(128) default ''::character varying,
    place_of_registration varchar(128) default ''::character varying,
    dir_path              varchar(128) default ''::character varying
);

alter table children
    owner to lokle_admin;

create unique index children_id_uindex
    on children (id);

create unique index children_user_id_uindex
    on children (user_id);




-- auto-generated definition
create table parents_children
(
    id           bigserial
        constraint parents_children_pk
            primary key,
    parent_id    bigint      not null
        constraint parents_children_parents_id_fk
            references parents
            on update cascade on delete cascade,
    child_id     bigint      not null
        constraint parents_children_children_id_fk
            references children
            on update cascade on delete cascade,
    relationship varchar(16)
);

alter table parents_children
    owner to lokle_admin;

create unique index parents_children_id_uindex
    on parents_children (id);





-- auto-generated definition
create table registration_requests
(
    id          bigserial
        constraint registration_requests_pk
            primary key,
    user_id     bigint                                           not null
        constraint registration_requests_users_id_fk
            references users
            on update cascade on delete cascade,
    manager_id  bigint
        constraint registration_requests_users_id_fk_2
            references users
            on update cascade on delete cascade,
    type        smallint                                         not null,
    status      varchar(16) default 'pending'::character varying not null,
    create_time bigint                                           not null,
    message     varchar(1024) default ''::character varying
);

alter table registration_requests
    owner to lokle_admin;

create unique index registration_requests_id_uindex
    on registration_requests (id);



drop table if exists registration_requests cascade;

drop table if exists parents_children cascade;

drop table if exists parents cascade;

drop table if exists children cascade;

drop table if exists users cascade;

