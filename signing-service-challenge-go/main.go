package main

import (
	"github.com/DrMonez/coding-challenges/signing-service-challenge/api"
	"github.com/DrMonez/coding-challenges/signing-service-challenge/domain"
	"github.com/DrMonez/coding-challenges/signing-service-challenge/persistence"
	"log"
)

const (
	ListenAddress = ":8080"
)

func main() {
	var storage persistence.Storage = &persistence.LocalStorage{
		UserDevices: make(map[string]map[string]struct{}),
		Devices:     make(map[string]*domain.Device),
		Signatures:  make(map[string]map[int]*domain.Signature),
	}

	server := api.NewServer(
		ListenAddress,
		storage,
	)

	if err := server.Run(); err != nil {
		log.Fatal("Could not start server on ", ListenAddress)
	}
}
