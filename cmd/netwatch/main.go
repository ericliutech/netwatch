package main

import (
	"fmt"
	"log"

	"github.com/ericliutech/netwatch/internal/config"
	"github.com/ericliutech/netwatch/internal/httpapi"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	router := httpapi.NewRouter(cfg)

	addr := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("netwatch listening on %s", addr)

	if err := router.Run(addr); err != nil {
		log.Fatal(err)
	}
}
