package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/TakeAway-Inc/backend/api"
)

func doResponse(w http.ResponseWriter, r *http.Request, staticUrl string, resp interface{}) {
	bb, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	bb = bytes.Replace(bb, []byte("%static%"), []byte(staticUrl), -1)

	if _, err = w.Write(bb); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) CreateOrder(w http.ResponseWriter, r *http.Request, restaurantId string) {
	ctx := r.Context()

	var newOrder api.NewOrder

	if err := json.NewDecoder(r.Body).Decode(&newOrder); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	orderId, err := s.db.CreateOrder(ctx, newOrder, restaurantId)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	order, err := s.db.GetOrder(ctx, orderId)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := api.CreateOrderResp(order)
	doResponse(w, r, s.staticUrl, resp)
}

func (s *Server) GetOrdersOfRestaurantByID(w http.ResponseWriter, r *http.Request, restaurantId string) {
	ctx := r.Context()

	orders, err := s.db.GetOrdersOfRestaurant(ctx, restaurantId)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := api.GetOrdersResp(orders)
	doResponse(w, r, s.staticUrl, resp)
}

func (s *Server) GetOrderByID(w http.ResponseWriter, r *http.Request, orderId string) {
	ctx := r.Context()

	order, err := s.db.GetOrder(ctx, orderId)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := api.GetOrderResp(order)
	doResponse(w, r, s.staticUrl, resp)
}

func (s *Server) UpdateOrderByBot(w http.ResponseWriter, r *http.Request, orderId string) {
	// TODO implement me
	panic("implement me")
}

func (s *Server) GetRestaurantPaymentOptions(w http.ResponseWriter, r *http.Request, restaurantId string) {
	ctx := r.Context()

	restaurantId, err := s.db.GetRestaurantID(ctx, restaurantId)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	paymentOptions, err := s.db.GetPaymentOptions(ctx, restaurantId)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := api.GetPaymentOptionsResp(paymentOptions)

	bb, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	bb = bytes.Replace(bb, []byte("%static%"), []byte(s.staticUrl), -1)

	if _, err = w.Write(bb); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) GetRestaurantMenu(w http.ResponseWriter, r *http.Request, restaurantId string) {
	ctx := r.Context()

	restaurantId, err := s.db.GetRestaurantID(ctx, restaurantId)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	categories, err := s.db.GetCategories(ctx, restaurantId)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	dishes, err := s.db.GetDishes(ctx, restaurantId)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	style, err := s.db.GetRestaurantStyle(ctx, restaurantId)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resp := api.GetRestaurantResp{
		Categories: categories,
		Dishes:     dishes,
		Style:      style,
	}

	bb, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	bb = bytes.Replace(bb, []byte("%static%"), []byte(s.staticUrl), -1)

	if _, err = w.Write(bb); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
