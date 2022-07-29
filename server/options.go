package server

import (
	"github.com/TakeAway-Inc/backend/db"
)

type Option func(server *Server)

func WithDB(db *db.DB) Option {
	return func(s *Server) {
		s.db = db
	}
}
