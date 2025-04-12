package metrics

import (
	"context"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Server struct {
	cfg        *Config
	httpServer *http.Server
}

func NewServer(cfg *Config) (*Server, error) {
	err := cfg.Validate()
	if err != nil {
		return nil, err
	}

	return &Server{cfg: cfg}, nil
}

func (s *Server) Run(_ context.Context) error {
	router := httprouter.New()
	router.Handler(http.MethodGet, path, promhttp.Handler())

	s.httpServer = &http.Server{
		Addr:              s.cfg.address,
		Handler:           router,
		ReadTimeout:       s.cfg.readTimeout,
		WriteTimeout:      s.cfg.writeTimeout,
		ReadHeaderTimeout: s.cfg.readHeaderTimeout,
	}

	return s.httpServer.ListenAndServe()
}

func (s *Server) Close() error {
	return s.httpServer.Close()
}
