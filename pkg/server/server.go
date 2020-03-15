package server

import (
	"context"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gofrp/fp-multiuser/pkg/server/controller"

	"github.com/gin-gonic/gin"
)

type Config struct {
	BindAddress string
	Tokens      map[string]string
}

type Server struct {
	cfg Config

	s    *http.Server
	done chan struct{}
}

func New(cfg Config) (*Server, error) {
	s := &Server{
		cfg:  cfg,
		done: make(chan struct{}),
	}
	if err := s.init(); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Server) Run() error {
	l, err := net.Listen("tcp", s.cfg.BindAddress)
	if err != nil {
		return err
	}
	log.Printf("HTTP server listen on %s", l.Addr().String())
	go func() {
		if err = s.s.Serve(l); err != http.ErrServerClosed {
			log.Printf("error shutdown HTTP server: %v", err)
		}
	}()
	<-s.done
	return nil
}

func (s *Server) Stop() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.s.Shutdown(ctx); err != nil {
		log.Fatalf("shutdown HTTP server error: %v", err)
	}
	log.Printf("HTTP server exited")
	close(s.done)
	return nil
}

func (s *Server) init() error {
	if err := s.initHTTPServer(); err != nil {
		log.Printf("init HTTP server error: %v", err)
		return err
	}
	return nil
}

func (s *Server) initHTTPServer() error {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()

	s.s = &http.Server{
		Handler: engine,
	}

	controller.NewOpController(s.cfg.Tokens).Register(engine)
	return nil
}
