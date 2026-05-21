package main

import (
	"encoding/json"
	"fmt"

	"github.com/ericliutech/netwatch/internal/discovery"
)

func runARP() error {
	observations, err := discovery.ReadARPObservations()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(observations, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(data))

	return nil
}
