package main

import (
	"fmt"
	"log"

	"github.com/akh-dev/encrypt/encryption-service/engine"

	"github.com/akh-dev/encrypt/encryption-service/config"
	"github.com/akh-dev/encrypt/encryption-service/service"
)

func main() {

	cfg, err := config.Get()
	if err != nil {
		log.Fatalf("Failed to load config: %+v", err)
	}

	encryptionEngine, err := engine.NewAESEngine()
	if err != nil {
		log.Fatalf("Failed to initialise encryption service: %+v", err)
	}

	encryptionService, err := service.New(cfg, encryptionEngine)
	if err != nil {
		log.Fatalf("Failed to initialise encryption service: %+v", err)
	}

	encryptionService.ListenAndServe()

	log.Println("Encryption-Service started, press <ENTER> to exit")
	fmt.Scanln()

}
