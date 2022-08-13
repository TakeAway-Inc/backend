package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/TakeAway-Inc/backend/api"
)

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

	resp := api.GetPaymentOptionsResponse(paymentOptions)

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
	resp := api.GetRestaurantResponse{
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
