-- +goose Up
create table restaurant_order_status
(
    value text not null primary key
);

insert into restaurant_order_status
values ('created'),
       ('processing'),
       ('canceled'),
       ('done');

create table restaurant_order
(
    order_id      uuid primary key default uuid_generate_v4(),
    restaurant_id uuid not null references restaurant (restaurant_id),
    status        text not null    default 'created' references restaurant_order_status (value),
    order_comment text not null    default ''
);


create table order_position
(
    order_id uuid    not null references restaurant_order (order_id),
    dish_id  uuid    not null references dish (dish_id),
    amount   integer not null
);

-- +goose Down
drop table order_position;

drop table restaurant_order;
