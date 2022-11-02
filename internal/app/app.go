package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/Meystergod/placements-api-service/pkg/shutdown"
	"golang.org/x/sync/errgroup"
	"net"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/Meystergod/placements-api-service/internal/config"
	"github.com/Meystergod/placements-api-service/internal/handlers/placements"
	httpclient "github.com/Meystergod/placements-api-service/pkg/client"
	"github.com/Meystergod/placements-api-service/pkg/logging"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
)

type App struct {
	logger     *logging.Logger
	cfg        *config.Config
	router     *httprouter.Router
	httpServer *http.Server
}

func NewApp(cfg *config.Config, logger *logging.Logger) (App, error) {
	logger.Info("router initializing")
	router := httprouter.New()

	logger.Info("http client initializing")
	client := httpclient.NewClient()

	logger.Info("handler initializing")
	handler := placements.NewHandler(logger, cfg, client)
	handler.Register(router)

	return App{
		logger: logger,
		cfg:    cfg,
		router: router,
	}, nil
}

func (s *App) Run(ctx context.Context) error {
	grp, ctx := errgroup.WithContext(ctx)
	grp.Go(func() error {
		return s.startHTTP(ctx)
	})

	return grp.Wait()
}

func (s *App) startHTTP(ctx context.Context) error {
	logger := s.logger.WithFields(map[string]interface{}{
		"IP":   s.cfg.HTTP.IP,
		"Port": s.cfg.HTTP.Port,
	})
	logger.Info("HTTP server initializing")

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", s.cfg.HTTP.IP, s.cfg.HTTP.Port))
	if err != nil {
		s.logger.WithError(err).Fatal("failed to create listener")
	}

	logger = s.logger.WithFields(map[string]interface{}{
		"AllowedMethods":     s.cfg.HTTP.CORS.AllowedMethods,
		"AllowedOrigins":     s.cfg.HTTP.CORS.AllowedOrigins,
		"AllowCredentials":   s.cfg.HTTP.CORS.AllowCredentials,
		"AllowedHeaders":     s.cfg.HTTP.CORS.AllowedHeaders,
		"OptionsPassthrough": s.cfg.HTTP.CORS.OptionsPassthrough,
		"ExposedHeaders":     s.cfg.HTTP.CORS.ExposedHeaders,
		"Debug":              s.cfg.HTTP.CORS.Debug,
	})

	logger.Info("cors initializing")
	c := cors.New(cors.Options{
		AllowedMethods:     s.cfg.HTTP.CORS.AllowedMethods,
		AllowedOrigins:     s.cfg.HTTP.CORS.AllowedOrigins,
		AllowCredentials:   s.cfg.HTTP.CORS.AllowCredentials,
		AllowedHeaders:     s.cfg.HTTP.CORS.AllowedHeaders,
		OptionsPassthrough: s.cfg.HTTP.CORS.OptionsPassthrough,
		ExposedHeaders:     s.cfg.HTTP.CORS.ExposedHeaders,
		Debug:              s.cfg.HTTP.CORS.Debug,
	})

	handler := c.Handler(s.router)

	s.httpServer = &http.Server{
		Handler:      handler,
		WriteTimeout: 250 * time.Millisecond,
		ReadTimeout:  250 * time.Millisecond,
	}

	go shutdown.Graceful(s.logger, []os.Signal{syscall.SIGABRT, syscall.SIGQUIT, syscall.SIGHUP, os.Interrupt, syscall.SIGTERM}, s.httpServer)

	logger.Info("application initialized and started")

	if err = s.httpServer.Serve(listener); err != nil {
		switch {
		case errors.Is(err, http.ErrServerClosed):
			logger.Warning("server shutdown")
		default:
			logger.Fatal(err)
		}
	}

	return err
}
