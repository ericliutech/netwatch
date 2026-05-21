package main

import (
	"context"
	"log"

	"github.com/ericliutech/netwatch/internal/discovery"
)

func runTCPProbe() error {
	log.Println("TCP probing local subnets...")

	if err := discovery.ProbeLocalSubnetsTCP(context.Background(), []uint16{80}); err != nil {
		return err
	}

	log.Println("TCP probe complete")
	return nil
}
