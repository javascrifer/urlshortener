package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"github.com/nikas-lebedenko/urlshortener/api"
	"github.com/nikas-lebedenko/urlshortener/repository/mongodb"
	"github.com/nikas-lebedenko/urlshortener/repository/redis"
	"github.com/nikas-lebedenko/urlshortener/shortener"
)

func main() {
	repository := newRepo()

	if repository == nil {
		panic("Repository not specified")
	}

	service := shortener.NewRedirectService(repository)
	handler := api.NewHandler(service)

	r := newRouter()
	r.Get("/{code}", handler.Get)
	r.Post("/", handler.Post)

	errs := make(chan error, 2)
	go func() {
		fmt.Println("Listening on port :8000")
		errs <- http.ListenAndServe(":8000", r)

	}()
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()

	fmt.Printf("Terminated %s", <-errs)
}

func newRepo() shortener.RedirectRepository {
	switch os.Getenv("URL_DB") {
	case "redis":
		url := os.Getenv("REDIS_URL")
		repo, err := redis.NewRepository(url)
		if err != nil {
			log.Fatal(err)
		}
		return repo
	case "mongo":
		url := os.Getenv("MONGO_URL")
		db := os.Getenv("MONGO_DB")
		timeout, _ := strconv.Atoi(os.Getenv("MONGO_TIMEOUT"))
		repo, err := mongodb.NewRepository(url, db, timeout)
		if err != nil {
			log.Fatal(err)
		}
		return repo
	}
	return nil
}

func newRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	return r
}
