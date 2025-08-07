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
)

func main() {
	config, err := settings.NewConfig()
	if err != nil {
		panic(err)
	}
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.DBName)
	db, err := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	if err != nil {
		panic(errors.Wrap(err, "failed to connect database"))
	}
	err = db.AutoMigrate(&repositories.Subscription{})
	if err != nil {
		panic(errors.Wrap(err, "failed to merge database"))
	}
	repository := repositories.NewRepository(db)

	service := services.NewService(repository)

	presentation := presentations.NewPresentation(service)
	presentation.BuildApp()
}
