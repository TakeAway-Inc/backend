package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"backend/api"
	"backend/db"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/pkg/errors"
)

type Server struct {
	httpServer *http.Server

	db *db.DB
}

func (s *Server) Run() {
	if err := s.httpServer.ListenAndServe(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			fmt.Println("Server closed")
		} else {
			fmt.Println("Error:", err)
		}
	}
}

type Config struct {
	ServerAddr string `yaml:"server_addr"`
}

func NewServer(opts ...Option) *Server {
	s := &Server{}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

const applicationJSONContentType = "application/json"

func (s *Server) AddHTTPServer(c Config) {

	mux := chi.NewRouter()
	mux.Use(middleware.NoCache)
	mux.Use(middleware.SetHeader("Content-Type", applicationJSONContentType))

	mux.Route("/", func(r chi.Router) {
		r.Mount("/", api.Handler(s))
	})

	s.httpServer = &http.Server{
		Handler: mux,
		Addr:    c.ServerAddr,
	}
}

func (s *Server) GetApiMenuRestaurantId(w http.ResponseWriter, r *http.Request, restaurantId string) {
	ctx := r.Context()
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

	if encErr := json.NewEncoder(w).Encode(resp); encErr != nil {
		http.Error(w, encErr.Error(), http.StatusInternalServerError)
		return
	}
}
