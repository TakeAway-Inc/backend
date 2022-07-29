package server

import (
	"backend/db"
)

type Option func(server *Server)

func WithDB(db *db.DB) Option {
	return func(s *Server) {
		s.db = db
	}
}
