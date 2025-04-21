create extension if not exists "pgcrypto";

create schema if not exists schema_name;

create table if not exists schema_name.users
(
    id          uuid        not null
        default gen_random_uuid()
        constraint urls_pk
            primary key,
    telegram_id     bigint null,
    is_staff bool default false not null,
    is_admin bool default false not null,
    created_at  timestamp default now() not null
);