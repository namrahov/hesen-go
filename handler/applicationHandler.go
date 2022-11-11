package handler

import (
	"encoding/json"
	mid "github.com/go-chi/chi/middleware"
	"github.com/gorilla/mux"
	"github.com/namrahov/hesen-go/config"
	"github.com/namrahov/hesen-go/middleware"
	"github.com/namrahov/hesen-go/repo"
	"github.com/namrahov/hesen-go/service"
	"net/http"
	"strconv"
)

type applicationHandler struct {
	Service service.IService
}

func ApplicationHandler(router *mux.Router) *mux.Router {

	router.Use(mid.Recoverer)
	router.Use(middleware.RequestParamsMiddleware)

	h := &applicationHandler{
		Service: &service.Service{
			ApplicationRepo: &repo.ApplicationRepo{},
		},
	}

	router.HandleFunc(config.RootPath+"/applications/{id}", h.getApplication).Methods("GET")

	return router
}

func (h *applicationHandler) getApplication(w http.ResponseWriter, r *http.Request) {

	idStr := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := h.Service.GetApplication(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}
