CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
create table if not exists account
(
    user_id        uuid    default uuid_generate_v1()    not null
        constraint account_pk
            primary key,
    email          varchar(150)                          not null,
    google_id      varchar(2048)                         not null,
    access_token   varchar(2048)                         not null,
    expires_in     timestamp                             not null,
    google_picture varchar,
    name           varchar default ''::character varying not null
);

create unique index if not exists account_email_uindex
    on account (email);

create unique index if not exists account_google_id_uindex
    on account (google_id);

create unique index if not exists account_uuid_uindex
    on account (user_id);

