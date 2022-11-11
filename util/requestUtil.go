package util

import (
	"encoding/json"
	"github.com/namrahov/hesen-go/model"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

func DecodeBody(w http.ResponseWriter, r *http.Request, data interface{}) error {
	logger := r.Context().Value(model.ContextLogger).(*log.Entry)

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		logger.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		logger.Error(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return err
	}

	return nil
}
