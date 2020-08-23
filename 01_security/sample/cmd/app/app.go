package app

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/jackc/pgx/v4/pgxpool"
	"net/http"
)

type Server struct {
	mux chi.Router
	pool *pgxpool.Pool
}

func NewServer(mux chi.Router, pool *pgxpool.Pool) *Server {
	return &Server{mux: mux, pool: pool}
}

func (s *Server) Init() error {
	s.mux.With(middleware.Logger).Get("/test", s.test)
	return nil
}

func (s *Server) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	s.mux.ServeHTTP(writer, request)
}

func (s *Server) test(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("test"))
}
