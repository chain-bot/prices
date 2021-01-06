create table meta_migration(
    migration_number integer not null,
    migration_name text not null,
    date date not null,
    primary key (migration_number)
);