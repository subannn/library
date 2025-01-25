package server

import (
	"EffectiveMobileTestTask/internal/handlers"
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log/slog"
	"net/http"
	"time"
)

type Server struct {
	logger  *slog.Logger
	echo    *echo.Echo
	handler *handlers.Handler
}

func NewServer(logger *slog.Logger, handler *handlers.Handler) *Server {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.HideBanner = true

	e.POST("/saveSong", handler.SaveSong)
	e.PUT("/updateSong/:id", handler.UpdateSong)
	e.DELETE("/deleteSong/:id", handler.DeleteSong)
	e.GET("/getSongText/:id", handler.GetSongText)
	e.GET("/getSongs", handler.GetSongs)
	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	return &Server{
		logger:  logger,
		echo:    e,
		handler: handler,
	}
}

func (s *Server) Start(port int) error {
	err := s.echo.Start(fmt.Sprintf(":%d", port))
	if err != nil {

		return err
	}

	s.logger.Info("Server started")
	return err
}

func (s *Server) Stop(timeout time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := s.echo.Shutdown(ctx); err != nil {
		s.logger.Info("Error during server shutdown: %v", err)
	} else {
		s.logger.Info("Server shut down gracefully")
	}
}
