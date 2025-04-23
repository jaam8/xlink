drop extension if exists "pgcrypto";

alter table user_service.users drop column token;

drop table if exists user_service.users;

drop schema if exists user_service;

