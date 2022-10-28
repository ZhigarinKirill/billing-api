package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/ZhigarinKirill/billing-api/config"
	"github.com/ZhigarinKirill/billing-api/internal/app/apiserver"
	"github.com/ZhigarinKirill/billing-api/internal/app/apiserver/handler"
	"github.com/ZhigarinKirill/billing-api/internal/app/repository"
	"github.com/ZhigarinKirill/billing-api/internal/app/service"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func main() {

	if err := config.InitConfig(); err != nil {
		log.Fatal().Err(err).Msg("error occurred while reading configuration file")
	}

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal().Err(err).Msg("error occurred while reading .env file")
	}

	db, err := repository.NewPostgresDB(&repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		User:     viper.GetString("db.user"),
		DBName:   viper.GetString("db.db_name"),
		SSLMode:  viper.GetString("db.ssl_mode"),
		Password: "0000",
	})
	if err != nil {
		log.Fatal().Err(err).Msg("error occurred while connecting to the postgres database")
	}

	repo := repository.NewRepository(db)
	services := service.NewService(repo)
	h := handler.NewHandler(services)

	server := apiserver.NewServer(viper.GetString("port"), h.InitRoutes())

	go func() {
		if err := server.Start(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("error occurred while running http server")
		}
	}()

	log.Print("Server started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	if err := server.Shutdown(context.Background()); err != nil {
		log.Fatal().Err(err).Msg("error occurred on server shutting down")
	}

	log.Print("Server stopped")

}
