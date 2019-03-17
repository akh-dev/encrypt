package service

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/pkg/errors"

	"github.com/akh-dev/encrypt/encryption-service/api"
	"github.com/akh-dev/encrypt/encryption-service/config"
	"github.com/akh-dev/encrypt/encryption-service/engine"
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

func (s *service) defaultHandler(w http.ResponseWriter, r *http.Request) {
	writeCommonHeaders(w)
	respondBadRequest(w, "unknown request", []string{})
}

func (s *service) handleStoreRequest(w http.ResponseWriter, r *http.Request) {
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

	newKey, err := s.engine.GenerateNewKey()
	if err != nil {
		log.Println(errors.Wrap(err, "failed to generate a new key during processing a store request"))
		respondInternalServerError(w, "internal server error", []string{})
		return
	}

	cipherText, err := s.engine.Encrypt([]byte(storeReq.Payload), newKey)
	if err != nil {
		log.Println(errors.Wrap(err, "failed to encrypt"))
		respondInternalServerError(w, "internal server error", []string{})
		return
	}

	if err := s.store(storeReq.Id, cipherText); err != nil {
		log.Println(errors.Wrap(err, "failed to store encoded text"))
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

func (s *service) handleRetrieveRequest(w http.ResponseWriter, r *http.Request) {
	writeCommonHeaders(w)

}

func (s *service) store(id string, ciphertext []byte) error {
	//TODO
	textB64 := base64.StdEncoding.EncodeToString(ciphertext)
	log.Printf("storing encoded text:\n%s\n", textB64)
	return nil
}

func respondBadRequest(w http.ResponseWriter, msg string, errors []string) {
	w.WriteHeader(http.StatusBadRequest)
	respObj := &api.Response{
		StatusCode:    http.StatusBadRequest,
		StatusMessage: msg,
		Errors:        errors,
	}
	writeResponse(w, respObj)
}

func respondInternalServerError(w http.ResponseWriter, msg string, errors []string) {
	w.WriteHeader(http.StatusInternalServerError)
	respObj := &api.Response{
		StatusCode:    http.StatusInternalServerError,
		StatusMessage: msg,
		Errors:        errors,
	}
	writeResponse(w, respObj)
}

func writeResponse(w http.ResponseWriter, respObj *api.Response) {
	response, err := json.Marshal(respObj)
	if err != nil {
		log.Println(err.Error())
		return
	}

	n, err := w.Write(response)
	if err != nil {
		log.Println(err.Error())
		return
	}
	log.Printf("wrote %d bytes in the response", n)

}

func writeCommonHeaders(w http.ResponseWriter) {
	w.Header().Add("Content-Type", "application/json")
}

func parseStoreRequest(r *http.Request) (*api.StoreRequest, error) {
	dec := json.NewDecoder(r.Body)
	storeReq := &api.StoreRequest{}
	if err := dec.Decode(storeReq); err != nil {
		err = errors.Wrap(err, "failed to parse Store request")
		log.Println(err.Error())
		return nil, err
	}

	return storeReq, nil
}
