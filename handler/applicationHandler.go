package handler

import (
	"encoding/json"
	mid "github.com/go-chi/chi/middleware"
	"github.com/gorilla/mux"
	"github.com/namrahov/hesen-go/config"
	"github.com/namrahov/hesen-go/middleware"
	"github.com/namrahov/hesen-go/model"
	"github.com/namrahov/hesen-go/repo"
	"github.com/namrahov/hesen-go/service"
	"github.com/namrahov/hesen-go/util"
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
			ValidationUtil:  &util.ValidationUtil{},
			CommentRepo:     &repo.CommentRepo{},
		},
	}

	router.HandleFunc(config.RootPath+"/applications/{id}", h.getApplication).Methods("GET")
	router.HandleFunc(config.RootPath+"/applications/{id}/change-status", h.changeStatus).Methods("PUT")
	router.HandleFunc(config.RootPath+"/applications/", h.saveApplication).Methods("POST")
	router.HandleFunc(config.RootPath+"/applications/get/filter-info", h.getFilterInfo).Methods("GET")

	return router
}

func (h *applicationHandler) saveApplication(w http.ResponseWriter, r *http.Request) {
	var application *model.Application
	err := util.DecodeBody(w, r, &application)
	if err != nil {
		return
	}

	result, err := h.Service.SaveApplication(r.Context(), application)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
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

func (h *applicationHandler) changeStatus(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var request model.ChangeStatusRequest
	err = util.DecodeBody(w, r, &request)
	if err != nil {
		return
	}

	errorResponse := h.Service.ChangeStatus(r.Context(), id, request)
	if errorResponse != nil {
		w.WriteHeader(errorResponse.Status)
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (h *applicationHandler) getFilterInfo(w http.ResponseWriter, r *http.Request) {

	result, err := h.Service.GetFilterInfo(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}
