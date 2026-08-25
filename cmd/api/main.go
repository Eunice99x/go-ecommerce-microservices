package main

import (
	"log"

	"github.com/eunice99x/goMicro/db"
	"github.com/eunice99x/goMicro/internal/handler"
	"github.com/eunice99x/goMicro/internal/repository"
	"github.com/eunice99x/goMicro/internal/service"
)

func main() {
	db, err := db.NewDatabase()
	if err != nil {
		log.Fatalf("error opening db: %v", err)
	}
	defer db.Close()

	log.Printf("successfully connected to database")

	store := repository.NewPostgresStorer(db.GetDB())
	service := service.NewService(store)
	hld := handler.NewHandler(service)

	r := handler.RegisterRoutes(hld)
	
	if err := handler.Start(":3000", r); err != nil {
    	log.Fatal(err)
	}
}
