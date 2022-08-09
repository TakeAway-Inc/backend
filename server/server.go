package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/TakeAway-Inc/backend/api"
	"github.com/TakeAway-Inc/backend/db"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/pkg/errors"
)

type Server struct {
	httpServer *http.Server
	staticUrl  string

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
	StaticUrl  string `yaml:"static_url"`
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
	s.staticUrl = c.StaticUrl

	corsOptions := cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}

	mux := chi.NewRouter()
	mux.Use(middleware.NoCache)
	mux.Use(cors.Handler(corsOptions))
	mux.Use(middleware.SetHeader("Content-Type", applicationJSONContentType))

	mux.Route("/", func(r chi.Router) {
		r.Mount("/", api.Handler(s))
		r.Mount("/static", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		}))
	})

	s.httpServer = &http.Server{
		Handler: mux,
		Addr:    c.ServerAddr,
	}
}

func (s *Server) GetApiMenuRestaurantId(w http.ResponseWriter, r *http.Request, restaurantId string) {
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
