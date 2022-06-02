alter table if exists transaction
    drop constraint check_split_max;

alter table if exists transaction
    drop column split_type;

