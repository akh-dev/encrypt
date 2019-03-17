package service

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/akh-dev/encrypt/storage-service/api"
	"github.com/pkg/errors"
)

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

func respondNotFound(w http.ResponseWriter, errors []string) {
	w.WriteHeader(http.StatusBadRequest)
	respObj := &api.Response{
		StatusCode:    http.StatusNotFound,
		StatusMessage: http.StatusText(http.StatusNotFound),
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

func parseStoreRequest(r *http.Request) (*api.IdMessage, error) {
	dec := json.NewDecoder(r.Body)
	storeReq := &api.IdMessage{}
	if err := dec.Decode(storeReq); err != nil {
		err = errors.Wrap(err, "failed to parse Store request")
		log.Println(err.Error())
		return nil, err
	}

	return storeReq, nil
}

func parseRetrieveRequest(r *http.Request) (*api.Id, error) {
	dec := json.NewDecoder(r.Body)
	retrieveReq := &api.Id{}
	if err := dec.Decode(retrieveReq); err != nil {
		err = errors.Wrap(err, "failed to parse Retrieve request")
		log.Println(err.Error())
		return nil, err
	}

	return retrieveReq, nil
}
