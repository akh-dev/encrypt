package service

import (
	"fmt"
	"log"
	"net/http"

	"github.com/akh-dev/encrypt/encryption-service/engine"

	"github.com/akh-dev/encrypt/encryption-service/config"
)

type service struct {
	config *config.Config
	engine engine.Interface
}

func New(cfg *config.Config, engine engine.Interface) (*service, error) {
	svc := &service{
		config: cfg,
		engine: engine,
	}

	return svc, nil
}

func (s *service) ListenAndServe() {
	http.HandleFunc("/store", s.handleStoreRequest)
	http.HandleFunc("/retrieve", s.handleRetrieveRequest)

	go func() {
		//err := http.ListenAndServeTLS(fmt.Sprintf(":%s", s.port), "cert.pem", "key.pem")
		err := http.ListenAndServe(fmt.Sprintf(":%s", s.config.Service.Port), nil)
		if err != nil {
			log.Fatal(err.Error())
		}
	}()
}

func (s *service) handleStoreRequest(w http.ResponseWriter, r *http.Request) {
}

func (s *service) handleRetrieveRequest(w http.ResponseWriter, r *http.Request) {
}
