package main

import (
	"larsa-e-commerce-microservice-v2/internal/api"
	"log"
)

func main() {

	if err := api.Start(); err != nil {
		log.Fatal(err)
	}

}
