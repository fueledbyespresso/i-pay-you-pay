alter table transaction
    add split_type int;

alter table transaction
    add constraint check_split_max
        check (transaction.split_type >= 1);

