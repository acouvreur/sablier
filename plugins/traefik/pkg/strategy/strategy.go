package strategy

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

// Net client is a custom client to timeout after 2 seconds if the service is not ready
var netClient = &http.Client{
	Timeout: time.Second * 2,
}

type SablierResponse struct {
	State string `json:"state"`
	Error string `json:"error"`
}

type Strategy interface {
	ServeHTTP(rw http.ResponseWriter, req *http.Request)
}

func getServiceStatus(request string) (string, error) {

	// This request wakes up the service if he's scaled to 0
	resp, err := netClient.Get(request)
	if err != nil {
		return "error", err
	}

	decoder := json.NewDecoder(resp.Body)
	var response SablierResponse
	err = decoder.Decode(&response)
	if err != nil {
		return "error from ondemand service", err
	}

	if resp.StatusCode >= 400 {
		return "error from ondemand service", errors.New(response.Error)
	}

	return response.State, nil
}
