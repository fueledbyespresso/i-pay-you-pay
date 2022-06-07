create table IF NOT EXISTS user_group_bridge
(
    "group" int
        constraint user_group_bridge_ledger_group_id_fk
            references ledger_group,
    "user"  uuid
        constraint user_group_bridge_account_user_id_fk
            references account
);

alter table transaction
    add constraint transaction_ledger_group_id_fk
        foreign key ("group") references ledger_group;

