package main

import (
	"context"
	"log"
	"os"
	"wallet-app/pkg/handler"
	"wallet-app/pkg/repository"
	"wallet-app/pkg/server"
	"wallet-app/pkg/service"

	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load("config.env"); err != nil {
		log.Fatalf("error loading env variables: %s", err.Error())
	}

	// repo := repository.NewRepository()
	postgres, err := repository.NewPG(context.Background(), repository.Config{
		Host:    os.Getenv("DB_HOST"),
		Port:    os.Getenv("DB_PORT"),
		User:    os.Getenv("DB_USER"),
		Pass:    os.Getenv("DB_PASS"),
		DBName:  os.Getenv("DB_NAME"),
		SSLMode: os.Getenv("DB_SSLMODE"),
	})

	if err != nil {
		log.Fatal(err)
	}

	repo := repository.NewRepository(postgres)
	service := service.NewService(repo)
	handler := handler.NewHandler(service)
	router := handler.RegisterRoutes()
	server := new(server.Server)
	if err := server.Run("8000", router); err != nil {
		log.Fatal(err)
	}
}
