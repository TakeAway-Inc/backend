-- +goose Up
alter table dish add column calories int not null default 0;
alter table dish add column fat int not null default 0;
alter table dish add column protein int not null default 0;
alter table dish add column carbo int not null default 0;

-- +goose Down
alter table dish drop column calories;
alter table dish drop column fat;
alter table dish drop column protein;
alter table dish drop column carbo;
