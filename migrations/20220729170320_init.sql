-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

create table restaurant
(
    restaurant_id         uuid primary key default uuid_generate_v4(),
    restaurant_shown_name text not null
);

create table restaurant_style
(
    restaurant_id    uuid primary key references restaurant (restaurant_id),
    icon_url         text not null,
    background_color text not null
);

create table category
(
    restaurant_id       uuid references restaurant (restaurant_id),
    category_id         uuid primary key default uuid_generate_v4(),
    category_shown_name text not null
);

create table dish
(
    dish_id            uuid primary key default uuid_generate_v4(),
    category_id        uuid references category (category_id),
    dish_shown_name    text    not null,
    dish_description   text    not null default '',
    dish_price         integer not null,
    dish_image_url     text    not null,
    is_unavailable     boolean not null default false,
    unavailable_reason text    not null default ''
);

-- +goose Down

drop table dish;
drop table category;
drop table restaurant_style;
drop table restaurant;
