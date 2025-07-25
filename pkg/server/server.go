package server

import (
	"context"
	"log"
	"net/http"
	"time"
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run(port string, hander http.Handler) error {
	s.httpServer = &http.Server{
		Addr:           ":" + port,
		Handler:        hander,
		MaxHeaderBytes: 1 << 20, // 1MB
		ReadTimeout:    time.Second * 10,
		WriteTimeout:   time.Second * 10,
	}
	log.Printf("Server started on http://localhost:%s", port)
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
