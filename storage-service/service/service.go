package service

import (
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/pkg/errors"

	"github.com/akh-dev/encrypt/storage-service/api"
	"github.com/akh-dev/encrypt/storage-service/config"
)

var (
	lock          sync.RWMutex
	NotFoundError = errors.New("text not found")
)

type Service struct {
	config  *config.Config
	storage map[string]string
}

func New(cfg *config.Config) (*Service, error) {
	svc := &Service{
		config:  cfg,
		storage: map[string]string{},
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

	err = s.store(storeReq.Id, storeReq.Payload)
	if err != nil {
		log.Println(errors.Wrap(err, "failed to parse request data"))
		respondInternalServerError(w, "internal server error", []string{})
		return
	}

	result, err := json.Marshal(api.Id{
		Id: storeReq.Id,
	})
	if err != nil {
		log.Println("error marshaling response: %s", err.Error())
		respondInternalServerError(w, "internal server error", []string{})
		return
	}

	respObj := &api.Response{
		StatusCode:    0,
		StatusMessage: "Success",
		Result:        result,
		Errors:        []string{},
	}

	writeResponse(w, respObj)
}

func (s *Service) handleRetrieveRequest(w http.ResponseWriter, r *http.Request) {
	writeCommonHeaders(w)

	if r.Method != http.MethodGet {
		respondBadRequest(w, "unknown request", []string{})
		return
	}

	retrieveReq, err := parseRetrieveRequest(r)
	if err != nil {
		log.Printf("failed to parse request data: %s", err.Error())
		respondBadRequest(w, "bad request", []string{})
		return
	}

	plaintext, err := s.retrieve(retrieveReq.Id)
	if err != nil {
		if err == NotFoundError {
			log.Printf("not found by id %s", retrieveReq.Id)
			respondNotFound(w, []string{fmt.Sprintf("text with id %s not found", retrieveReq.Id)})
		} else {
			log.Printf("error while retrieving text with id %s : %s", retrieveReq.Id, err.Error())
			respondInternalServerError(w, "internal server error", []string{})
		}
		return
	}

	result, err := json.Marshal(api.IdMessage{
		Id:      retrieveReq.Id,
		Payload: plaintext,
	})
	if err != nil {
		log.Println("error marshaling response: %s", err.Error())
		respondInternalServerError(w, "internal server error", []string{})
		return
	}

	respObj := &api.Response{
		StatusCode:    0,
		StatusMessage: "Success",
		Result:        result,
		Errors:        []string{},
	}

	writeResponse(w, respObj)
}

func (s *Service) keyHash(key string) string {
	hasher := sha512.New()
	saltedKey := fmt.Sprintf("%sx%s", key, s.config.Service.Salt)
	hasher.Write([]byte(saltedKey))
	sum := hasher.Sum(nil)
	return base64.URLEncoding.EncodeToString(sum)
}

func (s *Service) store(id, plaintext string) error {
	hash := s.keyHash(id)

	lock.Lock()
	defer lock.Unlock()

	if s.config.Service.Debug {
		log.Printf("storing\nid: %s\ntext: %s\n", hash, plaintext)
	}
	s.storage[hash] = plaintext

	return nil
}

func (s *Service) retrieve(id string) (string, error) {
	hash := s.keyHash(id)

	lock.RLock()
	defer lock.RUnlock()

	plaintext, ok := s.storage[hash]
	if !ok {
		return "", NotFoundError
	}

	if s.config.Service.Debug {
		log.Printf("ratriving\nid: %s\ntext: %s\n", id, plaintext)
	}

	return plaintext, nil
}
