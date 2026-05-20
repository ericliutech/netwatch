package main

import (
	"log"

	"github.com/ericliutech/netwatch/internal/httpapi"
)

func main() {
	router := httpapi.NewRouter()

	log.Println("netwatch listening on :8080")

	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
