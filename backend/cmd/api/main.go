package main

import (
	"log"

	"github.com/ryanprayoga/messhub/backend/internal/app"
)

func main() {
	application, err := app.New()
	if err != nil {
		log.Fatalf("bootstrap app: %v", err)
	}

	defer application.Close()

	if err := application.Listen(); err != nil {
		log.Fatalf("listen: %v", err)
	}
}
