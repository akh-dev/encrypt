package service

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	storageApi "github.com/akh-dev/encrypt/storage-service/api"

	"github.com/pkg/errors"

	"github.com/akh-dev/encrypt/encryption-service/api"
	"github.com/akh-dev/encrypt/encryption-service/config"
	"github.com/akh-dev/encrypt/encryption-service/engine"
)

var NotFoundError = errors.New("text not found")

type Service struct {
	config *config.Config
	engine engine.Interface
	client *http.Client
}

func New(cfg *config.Config, engine engine.Interface) (*Service, error) {
	svc := &Service{
		config: cfg,
		engine: engine,
		client: http.DefaultClient,
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
		log.Printf("failed to parse request data : %s", err.Error())
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

	if r.Method != http.MethodGet {
		respondBadRequest(w, "unknown request", []string{})
		return
	}

	retrieveReq, err := parseRetrieveRequest(r)
	if err != nil {
		log.Printf("failed to parse request data : %s", err.Error())
		respondBadRequest(w, "bad request", []string{})
		return
	}
	if s.config.Service.Debug {
		log.Printf("handleRetrieveRequest: request data: id:[%s], key:[%s]", retrieveReq.Id, retrieveReq.Key)
	}

	key, err := base64.StdEncoding.DecodeString(retrieveReq.Key)
	if err != nil {
		log.Printf("malformed key, failed to decode from base64 : %s", err.Error())
		respondInternalServerError(w, "internal server", []string{})
		return
	}

	payload, err := s.ProcessRetrieve([]byte(retrieveReq.Id), key)
	if err != nil {
		if err == NotFoundError {
			respondNotFound(w, []string{fmt.Sprintf("text with id %s not found", retrieveReq.Id)})
		} else {
			log.Printf("failed to process retrieve request: %s", err.Error())
			respondInternalServerError(w, "internal server error", []string{})
		}
		return
	}

	respObj := &api.Response{
		StatusCode:    0,
		StatusMessage: "Success",
		Result: api.IdMessage{
			Id:      retrieveReq.Id,
			Payload: string(payload),
		},
		Errors: []string{},
	}

	writeResponse(w, respObj)
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

	textB64 := base64.StdEncoding.EncodeToString(ciphertext)
	jsonreq := &storageApi.IdMessage{
		Id:      id,
		Payload: textB64,
	}
	buf, err := json.Marshal(jsonreq)
	if err != nil {
		return errors.Wrap(err, "failed to marshal storage request")
	}
	body := bytes.NewBuffer(buf)

	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("http://%s:%s%s", s.config.Storage.Host, s.config.Storage.Port, s.config.Storage.StoreUri),
		body,
	)
	if err != nil {
		return errors.Wrap(err, "failed to create storage request")
	}
	req.Header.Add("Content-Type", "application/json")

	timeout := time.Duration(s.config.Service.CtxTimeout) * time.Second
	ctx, _ := context.WithTimeout(context.Background(), timeout)
	r, err := s.client.Do(req.WithContext(ctx))
	if err != nil {
		return errors.Wrap(err, "failed do perform store request")
	}
	defer func() {
		log.Println(r.Body.Close())
	}()

	response, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return errors.Wrap(err, "failed to read response body")
	}

	parsed := &storageApi.Response{}
	err = json.Unmarshal(response, parsed)
	if err != nil {
		return errors.Wrap(err, "failed to parse response body")
	}

	if parsed.StatusCode != 0 {
		return errors.Wrap(err, fmt.Sprintf("unexpected return from the storage service: %d - %s, %s", parsed.StatusCode, parsed.StatusMessage, strings.Join(parsed.Errors, ":")))
	}

	return nil
}

func (s *Service) getFromStorage(id string) ([]byte, error) {
	if s.config.Service.Debug {
		log.Printf("getFromStorage id: [%s]\n", id)
	}

	jsonreq := &storageApi.Id{
		Id: id,
	}
	buf, err := json.Marshal(jsonreq)
	if err != nil {
		log.Printf("failed to marshal retrieve request: %s", err.Error())
		return nil, errors.Wrap(err, "failed to marshal retrieve request")
	}
	body := bytes.NewBuffer(buf)

	req, err := http.NewRequest(
		http.MethodGet,
		fmt.Sprintf("http://%s:%s%s", s.config.Storage.Host, s.config.Storage.Port, s.config.Storage.RetrieveUri),
		body,
	)
	if err != nil {
		log.Printf("failed to create retrieve request: %s", err.Error())
		return nil, errors.Wrap(err, "failed to create retrieve request")
	}

	req.Header.Add("Content-Type", "application/json")
	timeout := time.Duration(s.config.Service.CtxTimeout) * time.Second
	ctx, _ := context.WithTimeout(context.Background(), timeout)
	r, err := s.client.Do(req.WithContext(ctx))
	if err != nil {
		log.Printf("failed do perform retrieve request: %s", err.Error())
		return nil, errors.Wrap(err, "failed do perform retrieve request")
	}
	defer func() {
		log.Println(r.Body.Close())
	}()

	response, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("failed to read response body: %s", err.Error())
		return nil, errors.Wrap(err, "failed to read response body")
	}

	parsed := &storageApi.Response{}
	err = json.Unmarshal(response, parsed)
	if err != nil {
		log.Printf("failed to parse response body: %s", err.Error())
		return nil, errors.Wrap(err, "failed to parse response body")
	}

	if parsed.StatusCode == http.StatusNotFound {
		log.Printf("no text found for id %s", id)
		return nil, NotFoundError
	}

	if parsed.StatusCode != 0 {
		log.Printf("unexpected status code returned from the storage service: %d - %s, %s", parsed.StatusCode, parsed.StatusMessage, strings.Join(parsed.Errors, ":"))
		return nil, errors.Errorf("unexpected return from the storage service: %d - %s, %s", parsed.StatusCode, parsed.StatusMessage, strings.Join(parsed.Errors, ":"))
	}

	idAndText := &storageApi.IdMessage{}
	err = json.Unmarshal(parsed.Result, idAndText)
	if err != nil {
		log.Printf("unexpected return from the storage %v", parsed.Result)
		return nil, errors.Wrap(err, "unexpected return from the storage")
	}

	if s.config.Service.Debug {
		log.Printf("got id: [%s]\n    payload text: [%s]\n", idAndText.Id, idAndText.Payload)
	}

	payloadBytes, err := base64.StdEncoding.DecodeString(idAndText.Payload)
	if err != nil {
		log.Printf("malformed text, failed to decode from base64: %s", err.Error())
		return nil, errors.Wrap(err, "malformed text, failed to decode from base64")
	}

	return payloadBytes, nil
}

func (s *Service) ProcessRetrieve(id, aesKey []byte) (payload []byte, err error) {

	if len(aesKey) != 32 {
		return nil, errors.New("invalid key")
	}

	log.Printf("ProcessRetrieve: aesKey:[%s]", base64.StdEncoding.EncodeToString(aesKey[:]))
	cipherText, err := s.getFromStorage(string(id))
	if err != nil {
		if err == NotFoundError {
			return nil, err
		} else {
			return nil, errors.Wrap(err, "failed to retrieve text from storage")
		}
	}

	key := [32]byte{}
	for i, b := range aesKey {
		key[i] = b
	}

	log.Printf("ProcessRetrieve: key(array):[%s]", base64.StdEncoding.EncodeToString(key[:]))
	plaintext, err := s.engine.Decrypt(cipherText, &key)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decrypt")
	}

	if s.config.Service.Debug {
		log.Printf("decrypted message: \n%s\n", string(plaintext))
	}

	return plaintext, nil
}
