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
	"github.com/satori/go.uuid"
	"net/http"
	"strconv"
)

type applicationHandler struct {
	Service     service.IService
	UserService service.IUserService
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
		UserService: &service.UserService{
			UserRepo: &repo.UserRepo{},
		},
	}

	router.HandleFunc("/sign-up", h.signUp).Methods("POST")
	router.HandleFunc("/bar", h.bar).Methods("GET")
	router.HandleFunc(config.RootPath+"/applications/{id}", h.getApplication).Methods("GET")
	router.HandleFunc(config.RootPath+"/applications/{id}/change-status", h.changeStatus).Methods("PUT")
	router.HandleFunc(config.RootPath+"/applications/", h.saveApplication).Methods("POST")
	router.HandleFunc(config.RootPath+"/applications/get/filter-info", h.getFilterInfo).Methods("GET")

	return router
}

func (h *applicationHandler) bar(w http.ResponseWriter, r *http.Request) {
	if !h.UserService.AlreadyLoggedIn(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
}
func (h *applicationHandler) signUp(w http.ResponseWriter, r *http.Request) {
	cookie, err3 := r.Cookie("session-id")

	if err3 != nil {
		id := uuid.NewV4()
		cookie = &http.Cookie{
			Name:     "session-id",
			Value:    id.String(),
			HttpOnly: true,
		}
		http.SetCookie(w, cookie)
	}

	userAddress, err := h.UserService.GetUserIfExist(r.Context(), cookie.Value)

	var user model.User
	if err == nil && userAddress == nil {
		//process form submission
		var u *model.UserRegister
		err := util.DecodeBody(w, r, &u)
		if err != nil {
			return
		}

		userAddress, err := h.UserService.SaveUser(r.Context(), u)
		user = *userAddress
		if err != nil {
			return
		}

		err = h.UserService.SaveSession(r.Context(), cookie.Value, user.Id)
		if err != nil {
			return
		}

	} else if err == nil && userAddress != nil {
		user = *userAddress
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	} else {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&user)

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
