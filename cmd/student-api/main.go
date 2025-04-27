package main

import (
	"context"
	"fmt"
	"github/suryakantdsa/student-api/internal/config"
	"github/suryakantdsa/student-api/internal/http/handlers/student"
	"github/suryakantdsa/student-api/internal/storage/postgress"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	fmt.Println("welcome to student api")

	// load config
	cfg := config.MustLoad()
	// databse setup
	// storage, err := sqlite.New(cfg)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	storage, err := postgress.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	slog.Info("storage intilized", slog.String("env", cfg.Env), slog.String("version", "1.1.1"))
	// setup rounter .
	router := http.NewServeMux()
	router.HandleFunc("POST /api/students", student.New(storage))
	router.HandleFunc("GET /api/students/{id}", student.GetById(storage))

	// setup server /

	server := http.Server{
		Addr:    cfg.Addr,
		Handler: router,
	}

	//gracefull shutdown

	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal("failed to start the server")
		}
	}()

	<-done
	slog.Info("shutdowning the server")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Info("faild to shutdown server ", slog.String("erro", err.Error()))
	}
	// err := server.Shutdown(ctx)
	// if err != nil {
	// 	slog.Info("faild to shutdown server ", slog.String("erro", err.Error()))
	// }

	slog.Info("server shutdown sucessfully")

}
