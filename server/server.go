package server

import (
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
	ServerAddr string
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

	mux.Route("/api", func(r chi.Router) {
		r.Mount("/", api.Handler(s))
	})

	s.httpServer = &http.Server{
		Handler: mux,
		Addr:    c.ServerAddr,
	}
}

func (s *Server) GetApiMenuRestaurantId(w http.ResponseWriter, r *http.Request, restaurantId string) {
	// TODO implement me
	panic("implement me")
}
