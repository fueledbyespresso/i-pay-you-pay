create table contact
(
    user_id    uuid not null
        constraint contacts_account_user_id_fk
            references account,
    contact_id uuid not null
        constraint contacts_account_user_id_fk_2
            references account,
    constraint contacts_pk
        primary key (user_id, contact_id)
);

