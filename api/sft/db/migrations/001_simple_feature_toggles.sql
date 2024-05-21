create extension if not exists "uuid-ossp";

create schema sft;

create table if not exists sft.feature_toggles (
    id uuid primary key default uuid_generate_v4(),
    feature_name text not null,
    toggle_meta jsonb not null,
    enabled boolean not null default false
);

create table if not exists sft.users (
    id uuid primary key default uuid_generate_v4(),
    first_name text not null,
    last_name text not null
);

---- create above / drop below ----

drop table if exists sft.feature_toggles;
drop table if exists sft.users;
