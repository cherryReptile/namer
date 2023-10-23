create schema if not exists persons;

create type persons.gender as enum ('male', 'female');

create table if not exists persons.persons_table
(
    id         bigserial primary key,
    name       varchar(250)                 not null,
    surname    varchar(250)                 not null,
    patronymic varchar(250),
    age        smallint,
    gender     persons.gender default 'male',
    nation     varchar(100),
    created_at timestamp      default now() not null,
    updated_at timestamp
);

create or replace function update_updated_at() returns trigger
    language plpgsql
as
$$
begin
    new.updated_at = now();
    return new;
end;
$$;

create trigger persons_updated_at_trigger
    before update
    on persons.persons_table
    for each row
execute function update_updated_at();

create index on persons.persons_table ((lower(name)));
create index on persons.persons_table ((lower(surname)));
