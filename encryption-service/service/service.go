package service

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"

	"github.com/pkg/errors"

	"github.com/akh-dev/encrypt/encryption-service/api"
	"github.com/akh-dev/encrypt/encryption-service/config"
	"github.com/akh-dev/encrypt/encryption-service/engine"
)

type Service struct {
	config *config.Config
	engine engine.Interface
}

func New(cfg *config.Config, engine engine.Interface) (*Service, error) {
	svc := &Service{
		config: cfg,
		engine: engine,
	}

	return svc, nil
}

func (s *Service) ListenAndServe() {
	http.HandleFunc("/", s.defaultHandler)
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

func (s *Service) defaultHandler(w http.ResponseWriter, r *http.Request) {
	writeCommonHeaders(w)
	respondBadRequest(w, "unknown request", []string{})
}

func (s *Service) handleStoreRequest(w http.ResponseWriter, r *http.Request) {
	writeCommonHeaders(w)

	if r.Method != http.MethodPost {
		respondBadRequest(w, "unknown request", []string{})
		return
	}

	storeReq, err := parseStoreRequest(r)
	if err != nil {
		log.Println(errors.Wrap(err, "failed to parse request data"))
		respondBadRequest(w, "bad request", []string{})
		return
	}
	if s.config.Service.Debug {
		log.Printf("handleStoreRequest: request data: %s, %s", storeReq.Id, storeReq.Payload)
	}

	newKey, err := s.ProcessStore([]byte(storeReq.Id), []byte(storeReq.Payload))
	if err != nil {
		log.Printf("failed to process Store request: %s", err.Error())
		respondInternalServerError(w, "internal server error", []string{})
		return
	}

	newKeyB64 := base64.StdEncoding.EncodeToString(newKey[:])
	respObj := &api.Response{
		StatusCode:    0,
		StatusMessage: "Success",
		Result: api.IdKeyPair{
			Id:  storeReq.Id,
			Key: newKeyB64,
		},
		Errors: []string{},
	}

	writeResponse(w, respObj)
}

func (s *Service) handleRetrieveRequest(w http.ResponseWriter, r *http.Request) {
	writeCommonHeaders(w)
}

func (s *Service) ProcessStore(id, payload []byte) (aesKey []byte, err error) {

	newKey, err := s.engine.GenerateNewKey()
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate a new key during processing a store request")
	}

	cipherText, err := s.engine.Encrypt(payload, newKey)
	if err != nil {
		return nil, errors.Wrap(err, "failed to encrypt")
	}

	if err := s.sendToStorage(string(id), cipherText); err != nil {
		return nil, errors.Wrap(err, "failed to store encoded text")
	}

	return newKey[:], nil
}

func (s *Service) sendToStorage(id string, ciphertext []byte) error {
	//TODO
	textB64 := base64.StdEncoding.EncodeToString(ciphertext)
	log.Printf("storing encoded text:\n%s\n", textB64)
	return nil
}

func (s *Service) ProcessRetrieve(id, aesKey []byte) (payload []byte, err error) {
	return []byte{}, nil
	//TODO
}
