package main

import (
	"log"

	"github.com/nekr0z/muhadi/internal/app"
)

func main() {
	if err := app.New().Run(); err != nil {
		log.Fatal(err)
	}
}
