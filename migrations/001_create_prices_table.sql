-- +goose Up
create table if not exists prices (
    id bigint primary key,
    name text not null,
    category text not null,
    price integer not null,
    create_date date not null
);

-- +goose Down
drop table if exists prices;