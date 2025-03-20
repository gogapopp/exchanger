package main

import (
	"context"
	"exchanger/internal/repository"
	"exchanger/internal/server"
	"exchanger/internal/server/handlers"
	"exchanger/internal/service"
	"log"
	"net/http"
	"os"
	"os/signal"
)

const addr = ":8080"

func main() {
	ctx := context.Background()

	repository, err := repository.New(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer repository.Close()

	currencyService := service.NewCurrencyService(repository)
	exchangeService := service.NewExchangeRateService(repository)
	convertService := service.NewConvertService(repository, repository)

	handlers := handlers.New(currencyService, exchangeService, convertService)

	routes := server.Routes(handlers)

	httpServer := &http.Server{
		Addr:    addr,
		Handler: routes,
	}

	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()
	log.Println("server running at addres: ", addr)

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint)
	<-sigint
}
