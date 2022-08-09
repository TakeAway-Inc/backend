-- +goose Up
alter table restaurant
    add column url_name text not null default '';

alter table dish add column weight integer not null default 1;

-- +goose Down
alter table dish drop column weight;

alter table restaurant
    drop column url_name;
