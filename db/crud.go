package db

import (
	"context"

	"github.com/TakeAway-Inc/backend/api"

	"github.com/pkg/errors"
)

func (db *DB) GetRestaurantID(ctx context.Context, restaurantNameOrId string) (string, error) {
	var restaurantID string

	q := `
		SELECT restaurant_id from restaurant where url_name = $1 or restaurant_id::text = $1`

	err := db.pool.QueryRow(ctx, q, restaurantNameOrId).Scan(&restaurantID)
	if err != nil {
		return "", errors.Wrap(err, "can't get restaurant id")
	}

	return restaurantID, err
}

func (db *DB) GetPaymentOptions(ctx context.Context, restaurantId string) ([]api.PaymentOption, error) {
	var pos []api.PaymentOption

	q := `
		select shown_name, description, image_url, forward_url from payment_option po
			left join restaurant_payment_option rpo on po.payment_option_id = rpo.payment_option_id where rpo.restaurant_id = $1
			order by rpo.payment_option_order`

	rows, err := db.pool.Query(ctx, q, restaurantId)
	if err != nil {
		return nil, errors.Wrap(err, "can't get payment options")
	}
	defer rows.Close()

	for rows.Next() {
		var po api.PaymentOption
		err := rows.Scan(&po.ShownName, &po.Description, &po.ImageUrl, &po.PaymentForwardUrl)
		if err != nil {
			return nil, errors.Wrap(err, "can't scan payment option")
		}
		pos = append(pos, po)
	}

	return pos, nil
}

func (db *DB) GetRestaurantStyle(ctx context.Context, id string) (api.RestaurantStyle, error) {
	var rs api.RestaurantStyle
	q := `
	SELECT
	r.restaurant_id, r.restaurant_shown_name, icon_url, background_color
	FROM
	restaurant_style
	rs
	left
	join
	restaurant
	r
	on
	r.restaurant_id = rs.restaurant_id
	WHERE
	r.restaurant_id = $1
	`

	err := db.pool.QueryRow(ctx, q, id).Scan(&rs.Id, &rs.RestaurantShownName, &rs.IconUrl, &rs.BackgroundColor)
	if err != nil {
		return rs, errors.Wrap(err, "can't get restaurant style")
	}

	return rs, err
}

func (db *DB) GetCategories(ctx context.Context, restaurantId string) ([]api.Category, error) {
	var categories []api.Category

	q := `
	SELECT
	c.category_id, c.category_shown_name
	FROM
	category
	c
	WHERE
	c.restaurant_id = $1
	`

	rows, err := db.pool.Query(ctx, q, restaurantId)
	if err != nil {
		return categories, errors.Wrap(err, "can't get categories")
	}
	defer rows.Close()

	for rows.Next() {
		var c api.Category
		err := rows.Scan(&c.Id, &c.ShownName)
		if err != nil {
			return categories, errors.Wrap(err, "can't scan category")
		}
		categories = append(categories, c)
	}

	return categories, nil
}

func (db *DB) GetDishes(ctx context.Context, restaurantId string) ([]api.Dish, error) {
	var dishes []api.Dish

	q := `
	SELECT
	d.dish_id, d.category_id, d.dish_shown_name, d.dish_description, d.dish_price, d.dish_image_url, d.dish_image_url, d.weight, d.calories, d.protein, d.fat, d.carbo
	FROM
	dish
	d
	left
	join
	category
	c
	on
	c.category_id = d.category_id
	left
	join
	restaurant
	r
	on
	r.restaurant_id = c.restaurant_id
	WHERE
	r.restaurant_id = $1
	`

	rows, err := db.pool.Query(ctx, q, restaurantId)
	if err != nil {
		return dishes, errors.Wrap(err, "can't get dishes")
	}
	defer rows.Close()

	for rows.Next() {
		var d api.Dish
		err := rows.Scan(&d.Id, &d.CategoryId, &d.ShownName, &d.Description, &d.Price.Amount, &d.ImageUrl, &d.PreviewImageUrl, &d.Weight, &d.Calories, &d.Proteins, &d.Fats, &d.Carbohydrates)
		if err != nil {
			return dishes, errors.Wrap(err, "can't scan dish")
		}
		dishes = append(dishes, d)
	}

	return dishes, nil
}
