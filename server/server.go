package server

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

type Config struct {
	Host string
	Port string
}

type Server struct {
	http_server *http.Server
	Config      *Config
}

func New(cfg *Config) *Server {
	return &Server{
		http_server: &http.Server{},
		Config:      cfg,
	}
}

func (s *Server) Run(handler http.Handler) error {
	s.http_server = &http.Server{
		Addr:    s.Config.Host + ":" + s.Config.Port,
		Handler: handler,
	}

	logrus.Infof("Server is now listening on %s",
		s.http_server.Addr)

	return s.http_server.ListenAndServe()
}
