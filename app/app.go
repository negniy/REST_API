package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"restapi/app/service"
	"restapi/config"
	"restapi/database"
	"syscall"
)

func Run() error {

	cnf, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	db, err := database.InitDB()
	if err != nil {
		log.Fatalf("Failed to find database: %v", err)
	}

	carService := service.NewService(*db)

	router := http.NewServeMux()
	router.HandleFunc("POST /cars", carService.Create)
	router.HandleFunc("GET /cars", carService.GetAll)
	router.HandleFunc("GET /cars/{id}", carService.Get)
	router.HandleFunc("PUT /cars/{id}", carService.Update)
	router.HandleFunc("PATCH /cars/{id}", carService.UpdateSomething)
	router.HandleFunc("DELETE /cars/{id}", carService.Delete)

	srv := http.Server{
		Addr:    cnf.Port,
		Handler: router,
	}

	go func() {
		log.Printf("run server: http://localhost%s", cnf.Port)
		err := srv.ListenAndServe()
		if err != nil {
			log.Printf("error when listen and serve: %s", err)
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(ch)
	sig := <-ch
	log.Printf("%s %v - %s", "Reseived shutdown signal", sig, "")
	return srv.Shutdown(context.Background())
}

func main() {

	if err := Run(); err != nil {
		log.Fatalf("Server stopped with error: %v", err)
	}
}
