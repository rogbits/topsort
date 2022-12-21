package lib

import (
	"net/http"
)

type Server struct {
	Server *http.Server
	Logger *Logger
}

func NewHttpServer(logger *Logger, api http.Handler) *Server {
	s := new(Server)
	s.Logger = logger
	s.Server = new(http.Server)
	s.Server.Handler = api
	return s
}

func (s *Server) Start(addr string) error {
	s.Server.Addr = addr
	s.Logger.Log("starting server on", addr)
	err := s.Server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (s *Server) Stop() error {
	s.Logger.Log("stopping server")
	err := s.Server.Close()
	if err != nil {
		return err
	}

	return nil
}
