-- +goose Up
create table payment_option
(
    payment_option_id uuid primary key default uuid_generate_v4(),
    shown_name        text not null,
    description       text not null,
    image_url         text not null,
    forward_url       text
);

create table restaurant_payment_option
(
    restaurant_id        uuid    not null references restaurant (restaurant_id),
    payment_option_id    uuid    not null references payment_option (payment_option_id),
    payment_option_order integer not null default 0,
    primary key (restaurant_id, payment_option_id)
);

-- +goose Down

drop table restaurant_payment_option;
drop table payment_option;
