package main

import (
	"log"
	"myapp/internal/pkg/app"
)

// air bench
func main() {
	a, err := app.New()
	if err != nil {
		log.Fatal(err)
	}

	err = a.Run()
	if err != nil {
		log.Fatal(err)
	}
}
