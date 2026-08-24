package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/eunice99x/goMicro/db"
	"github.com/eunice99x/goMicro/internal/repository"
	"github.com/eunice99x/goMicro/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	db, err := db.NewDatabase()
	if err != nil {
		log.Fatalf("error opening db: %v", err)
	}
	defer db.Close()

	log.Printf("successfully connected to database")

	store := repository.NewPostgresStorer(db.GetDB())
	service := service.NewService(store)

	fmt.Print(store, service)

	http.ListenAndServe(":3000", r)
}
