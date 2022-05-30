create table ledger_group
(
    description varchar(256) not null,
    id          serial
        constraint ledger_group_pk
            primary key
);

create unique index ledger_group_id_uindex
    on ledger_group (id);

create table transaction
(
    total               integer generated always as identity,
    description         varchar(256) default ''::character varying not null,
    time_of_transaction date                                       not null,
    time_of_record      timestamp                                  not null,
    recorder            uuid                                       not null
        constraint transaction_account_user_id_fk
            references account
            on delete cascade,
    "group"             integer,
    id                  serial
        constraint transaction_pk
            primary key
);

comment on column transaction.time_of_transaction is 'When the transaction occured';

comment on column transaction.time_of_record is 'When the transaction was recorded';

comment on column transaction."group" is 'Belongs to group of ID';

create unique index transaction_id_uindex
    on transaction (id);

create table user_transaction_bridge
(
    transaction_id integer not null
        constraint account_transaction_bridge_transaction_id_fk
            references transaction,
    user_id        uuid    not null
        constraint user_transaction_bridge_account_user_id_fk
            references account,
    constraint user_transaction_bridge_pk
        primary key (transaction_id, user_id)
);


