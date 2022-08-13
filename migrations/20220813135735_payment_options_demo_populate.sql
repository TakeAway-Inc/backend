-- +goose Up
INSERT INTO payment_option (payment_option_id, shown_name, description, image_url, forward_url)
VALUES ('2f94860c-7455-4e4a-a02f-34a2a300df3f', 'Наличные', 'аааааа', '%static%/istockphoto-1153024228-612x612.jpg', null);
INSERT INTO payment_option (payment_option_id, shown_name, description, image_url, forward_url)
VALUES ('9ba50d86-0ad4-4a6b-b027-6495da9049ba', 'Карта', 'фффф', '', null);

insert into restaurant_payment_option(restaurant_id, payment_option_id, payment_option_order)
VALUES ((select restaurant_id from restaurant where url_name = 'demo'), '2f94860c-7455-4e4a-a02f-34a2a300df3f', 1);
insert into restaurant_payment_option(restaurant_id, payment_option_id, payment_option_order)
VALUES ((select restaurant_id from restaurant where url_name = 'demo'), '9ba50d86-0ad4-4a6b-b027-6495da9049ba', 2);

-- +goose Down

delete
from restaurant_payment_option
where restaurant_id = (select restaurant_id from restaurant where url_name = 'demo');

delete
from payment_option
where payment_option_id in ('2f94860c-7455-4e4a-a02f-34a2a300df3f', '9ba50d86-0ad4-4a6b-b027-6495da9049ba');
