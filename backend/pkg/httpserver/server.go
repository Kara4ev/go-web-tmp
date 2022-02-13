package httpserver

import (
	"context"
	"net/http"
	"time"
)

const (
	_defaultReadTimeout     = 10 * time.Second
	_defaultWriteTimeout    = 10 * time.Second
	_defaultAddr            = ":8080"
	_defaultShutdownTimeout = 3 * time.Second
)

type Server struct {
	server          *http.Server
	notify          chan error
	shutdownTimeout time.Duration
}

type SConfig struct {
	Hendler         http.Handler
	Addr            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
}

func New(sc SConfig, opts ...Option) *Server {

	addr := _defaultAddr
	readTimeout := _defaultReadTimeout
	writeTimeout := _defaultWriteTimeout
	shutdownTimeout := _defaultShutdownTimeout

	if sc.Addr != "" {
		addr = sc.Addr
	}

	if sc.ReadTimeout != 0 {
		readTimeout = sc.ReadTimeout
	}

	if sc.WriteTimeout != 0 {
		writeTimeout = sc.WriteTimeout
	}

	if sc.ShutdownTimeout != 0 {
		shutdownTimeout = sc.ShutdownTimeout
	}

	httpServer := &http.Server{
		Handler:      sc.Hendler,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		Addr:         addr,
	}

	s := &Server{
		server:          httpServer,
		notify:          make(chan error, 1),
		shutdownTimeout: shutdownTimeout,
	}

	for _, opt := range opts {
		opt(s)
	}

	s.start()

	return s
}

func (s *Server) start() {
	go func() {
		s.notify <- s.server.ListenAndServe()
		close(s.notify)
	}()
}

func (s *Server) Notify() <-chan error {
	return s.notify
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	return s.server.Shutdown(ctx)
}
