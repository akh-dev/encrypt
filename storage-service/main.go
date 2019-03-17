package main

import (
	"fmt"
	"log"

	"github.com/akh-dev/encrypt/storage-service/config"
	"github.com/akh-dev/encrypt/storage-service/service"
)

func main() {

	cfg, err := config.Get()
	if err != nil {
		log.Fatalf("Failed to load config: %+v", err)
	}

	storageService, err := service.New(cfg)
	if err != nil {
		log.Fatalf("Failed to initialise storage service: %+v", err)
	}

	storageService.ListenAndServe()

	log.Println("Storage-Service started, press <ENTER> to exit")
	fmt.Scanln()

}
