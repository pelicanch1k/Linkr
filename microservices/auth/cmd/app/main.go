package main

import (
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/pelicanch1k/Linkr/auth/internal/config"
	"github.com/pelicanch1k/Linkr/auth/internal/config/db"
	"github.com/pelicanch1k/Linkr/auth/internal/mvc/handler"
	"github.com/pelicanch1k/Linkr/auth/internal/mvc/repository/postgres"
	"github.com/pelicanch1k/Linkr/auth/internal/mvc/service"
	"github.com/pelicanch1k/Linkr/auth/internal/router"
	"github.com/pelicanch1k/Linkr/auth/pkg/database"
	"github.com/pelicanch1k/ProductGatewayAPI/pkg/logging"
	"github.com/spf13/viper"
)

func main() {
    logger := logging.GetLogger()

	envPath := filepath.Join("../..", ".env")

    if err := godotenv.Load(envPath); err != nil {
		logger.Fatalf("error loading env variables: %s", err.Error())
	}

	if err := initConfig(); err != nil {
		logger.Fatalf("error init configs: %s", err.Error())
	}

	configDB := db.NewPostgresConfig()
	configAuth := config.NewAuthConfig() 


	driver, err := database.NewPostgresDriver(configDB)
	if err != nil {
		logger.Fatalf("failed to initialize db: %s", err.Error())
	}

	repos := postgres.NewAuthRepository(driver)
	service := service.NewService(repos, configAuth)
	handler := handler.NewHandler(service, logger)

	router := router.NewRouter(handler)
	
	router.Listen(":" + viper.GetString("port"))
}

func initConfig() error {
	viper.AddConfigPath("../../config")
	viper.SetConfigName("local")
	return viper.ReadInConfig()
}