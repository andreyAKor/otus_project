package app

import (
	"context"
	"io"

	"github.com/andreyAKor/otus_project/internal/http/server"

	"github.com/rs/zerolog/log"
)

var _ io.Closer = (*App)(nil)

type App struct {
	srv *server.Server
}

func New(srv *server.Server) (*App, error) {
	return &App{srv}, nil
}

// Run application.
func (a *App) Run(ctx context.Context) error {
	go func() {
		if err := a.srv.Run(ctx); err != nil {
			log.Fatal().Err(err).Msg("http-server listen fail")
		}
	}()

	return nil
}

// Close application.
func (a *App) Close() error {
	return nil
}
