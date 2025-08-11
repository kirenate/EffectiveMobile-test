package main

import (
	"fmt"
	"github.com/pkg/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"main.go/presentations"
	"main.go/repositories"
	"main.go/services"
	"main.go/settings"
	"main.go/utils"
)

func main() {
	log := utils.MakeLogger()

	log.Info().Msg("logger.started")

	err := settings.NewConfig()
	if err != nil {
		panic(errors.Wrap(err, "failed to create config"))
	}
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		settings.MyConfig.Host, settings.MyConfig.Port, settings.MyConfig.User, settings.MyConfig.Password, settings.MyConfig.DBName)

	db, err := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	if err != nil {
		panic(errors.Wrap(err, "failed to connect database"))
	}

	err = db.AutoMigrate(&repositories.Subscription{})
	if err != nil {
		panic(errors.Wrap(err, "failed to merge database"))
	}

	log.Info().Msgf("db.%s.started.at.%s:%d", settings.MyConfig.DBName, settings.MyConfig.Host, settings.MyConfig.Port)

	repository := repositories.NewRepository(db)

	service := services.NewService(repository)

	presentation := presentations.NewPresentation(service)
	app := presentation.BuildApp()

	err = app.Listen(settings.MyConfig.Addr)
	if err != nil {
		panic(errors.Wrap(err, "failed to start server"))
	}

	log.Info().Msgf("server.started.at.%s", settings.MyConfig.Addr)

}
