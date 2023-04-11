package app

import (
	"context"
	"dev/profileSaver/internal/config"
	controller "dev/profileSaver/internal/controller/v1"
	"dev/profileSaver/internal/model"
	"dev/profileSaver/internal/repository"
	"dev/profileSaver/internal/server"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run(cfg config.Config) error {
	var err error
	repo := repository.New()
	repo.CreateUser(model.User{
		Email:    "admin",
		Username: "admin",
		Password: "admin",
		Admin:    true,
	})

	handler := controller.New(repo)

	srv := new(server.Server)
	defer func() {
		log.Info().Msg("App Shutting Down")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err = srv.Shutdown(ctx)
		if err != nil {
			log.Error().Err(err)
		}
		log.Info().Msg("Server Stopped")
	}()

	errChan := make(chan error, 1)

	go func() {
		if err = srv.Run(cfg.Server.Port, handler.InitRouter()); err != nil {
			errChan <- err
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	select {
	case <-quit:
	case err = <-errChan:
		return err
	}

	return nil
}
