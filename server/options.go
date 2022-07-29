package server

import (
	"backend/db"

	log "github.com/sirupsen/logrus"
)

type Option func(server *Server)

func WithDB(db *db.DB) Option {
	return func(s *Server) {
		s.db = db
	}
}

func WithLogger(log *log.Logger) Option {
	return func(s *Server) {
		s.log = log
	}
}
