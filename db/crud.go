package db

import (
	"context"

	"backend/api"
)

func (db *DB) GetRestaurantStyle(ctx context.Context, id string) (api.RestaurantStyle, error) {
	var rs api.RestaurantStyle
	err := db.pool.QueryRow(ctx, `SELECT  FROM restaurant_style WHERE id = $1`, id).Scan(&rs.ID, &rs.Name)
	return restaraunt, err
}
