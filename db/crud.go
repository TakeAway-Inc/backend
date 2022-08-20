package db

import (
	"context"

	"github.com/TakeAway-Inc/backend/api"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"

	"github.com/pkg/errors"
)

func (db *DB) UpdateOrderByBot(ctx context.Context, orderId string, updatedOrder api.UpdatedOrderByBot) error {
	tx, err := db.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	if updatedOrder.Status != nil {
		q := `
			update restaurant_order
			set status = $1
			where order_id = $2`

		if _, err = tx.Exec(ctx, q, updatedOrder.Status, updatedOrder.Status, orderId); err != nil {
			return errors.Wrap(err, "can't update order")
		}
	}

	return nil
}

func (db *DB) GetOrdersOfRestaurant(ctx context.Context, restaurantId string) ([]api.Order, error) {
	tx, err := db.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	q := `
		select order_id,
			   status,
			   order_comment
		from restaurant_order
		where restaurant_id = $1
		`

	rows, err := tx.Query(ctx, q, restaurantId)
	if err != nil {
		return nil, errors.Wrap(err, "can't get orders of restaurant")
	}

	var orders []api.Order
	for rows.Next() {
		var order api.Order
		if err = rows.Scan(
			&order.OrderId,
			&order.Status,
			&order.Comment,
		); err != nil {
			return nil, errors.Wrap(err, "can't scan order")
		}

		positions, err := db.GetOrderPositions(ctx, order.OrderId)
		if err != nil {
			return nil, errors.Wrap(err, "can't get order positions")
		}

		order.Positions = positions
		orders = append(orders, order)
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, errors.Wrap(err, "can't commit transaction")
	}

	return orders, nil
}

func (db *DB) GetOrder(ctx context.Context, orderId string) (api.Order, error) {
	q := `
		select order_id,
			   restaurant_id,
			   status,
			   order_comment
		from restaurant_order
		where order_id = $1`

	var order api.Order

	if err := db.pool.QueryRow(ctx, q, orderId).Scan(
		&order.OrderId,
		&order.RestaurantId,
		&order.Status,
		&order.Comment,
	); err != nil {
		return order, errors.Wrap(err, "can't get order")
	}

	positions, err := db.GetOrderPositions(ctx, orderId)
	if err != nil {
		return order, errors.Wrap(err, "can't get order positions")
	}
	order.Positions = positions

	return order, nil
}

func (db *DB) GetOrderPositions(ctx context.Context, orderId string) ([]api.OrderPosition, error) {
	q := `
		select op.amount,
			   d.dish_id,
			   category_id,
			   dish_shown_name,
			   dish_description,
			   dish_price,
			   dish_image_url,
			   weight,
			   calories,
			   fat,
			   protein,
			   carbo
		
		from order_position op
				 left join dish d on d.dish_id = op.dish_id
		where order_id = $1
		`

	rows, err := db.pool.Query(ctx, q, orderId)
	if err != nil {
		return nil, errors.Wrap(err, "can't get order positions")
	}

	var positions []api.OrderPosition
	for rows.Next() {
		var pos api.OrderPosition
		pos.Dish.Price.Amount = new(int)
		if err = rows.Scan(
			&pos.Quantity,
			&pos.Dish.DishId,
			&pos.Dish.CategoryId,
			&pos.Dish.ShownName,
			&pos.Dish.Description,
			pos.Dish.Price.Amount,
			&pos.Dish.ImageUrl,
			&pos.Dish.Weight,
			&pos.Dish.Calories,
			&pos.Dish.Fats,
			&pos.Dish.Proteins,
			&pos.Dish.Carbohydrates,
		); err != nil {
			return nil, errors.Wrap(err, "can't scan order position")
		}
		positions = append(positions, pos)
	}

	return positions, nil
}
func (db *DB) CreateOrder(ctx context.Context, order api.NewOrder, restaurantId string) (orderId string, err error) {
	tx, err := db.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return "", err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	orderId = uuid.New().String()

	q := `insert into restaurant_order (order_id, restaurant_id) values ($1, $2)`
	if _, err = tx.Exec(ctx, q, orderId, restaurantId); err != nil {
		return "", errors.Wrap(err, "can't create order")
	}

	q = `insert into order_position (order_id, dish_id, amount) VALUES ($1, $2, $3)`
	for _, position := range order.Positions {
		dish := position.Dish
		quantity := position.Quantity
		if _, err = tx.Exec(ctx, q, orderId, dish.DishId, quantity); err != nil {
			return "", errors.Wrap(err, "can't add position to order")
		}
	}
	if err = tx.Commit(ctx); err != nil {
		return "", errors.Wrap(err, "can't commit tx")
	}

	return orderId, err
}

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
		err := rows.Scan(&d.DishId, &d.CategoryId, &d.ShownName, &d.Description, &d.Price.Amount, &d.ImageUrl, &d.PreviewImageUrl, &d.Weight, &d.Calories, &d.Proteins, &d.Fats, &d.Carbohydrates)
		if err != nil {
			return dishes, errors.Wrap(err, "can't scan dish")
		}
		dishes = append(dishes, d)
	}

	return dishes, nil
}
