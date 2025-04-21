create extension if not exists "pgcrypto";

create schema if not exists schema_name;

create table if not exists schema_name.urls
(
    id          uuid        not null
        default gen_random_uuid()
        constraint urls_pk
        primary key,
    user_id     uuid        not null,
    group_id    uuid        not null,
    "generated"   bool      default true    not null,
    short_link  varchar(10) not null,
    url         text not null,
    created_at  timestamp default now() not null,
    expire_at   timestamp not null
);