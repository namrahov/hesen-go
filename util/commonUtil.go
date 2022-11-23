package util

import (
	"encoding/json"
	"github.com/namrahov/hesen-go/model"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func HandleError(w http.ResponseWriter, err *model.ErrorResponse) {
	w.Header().Add(model.ContentTypeString, model.JSONType)
	w.WriteHeader(err.Status)
	encodeErr := json.NewEncoder(w).Encode(err)

	if encodeErr != nil {
		log.Error("ActionLog.HandleError.error happened when encode json")
	}
}
