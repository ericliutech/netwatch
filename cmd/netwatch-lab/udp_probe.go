package main

import (
	"context"
	"log"

	"github.com/ericliutech/netwatch/internal/discovery"
)

func runUDPProbe() error {
	log.Println("UDP probing local subnets...")

	if err := discovery.ProbeLocalSubnetsUDP(context.Background(), []uint16{9}); err != nil {
		return err
	}

	log.Println("UDP probe complete")
	return nil
}
