package dap

import (
	"context"
	"fmt"
	"net/http"

	"golang.org/x/sync/errgroup"
)

type Server struct {
	*Options
}

func NewServer(options *Options) *Server {
	return &Server{
		Options: options,
	}
}

func (server *Server) Run() error {
	eg, ctx := errgroup.WithContext(context.Background())
	ctx, cancel := context.WithCancel(ctx)

	for port, backend := range server.Backends {
		mux := http.NewServeMux()
		mux.HandleFunc(server.HealthCheck, HandleHealth)
		mux.Handle("/", NewAuthHandler(ctx, backend, &server.AuthOptions))

		eg.Go(func() error {
			return http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
		})
	}

	err := eg.Wait()
	cancel()

	return err
}
