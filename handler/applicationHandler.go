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
	"golang.org/x/crypto/bcrypt"
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
	router.HandleFunc("/logged-in", h.loggedIn).Methods("POST")
	router.HandleFunc("/logout", h.logout).Methods("GET")
	router.HandleFunc(config.RootPath+"/applications/{id}", h.getApplication).Methods("GET")
	router.HandleFunc(config.RootPath+"/applications/{id}/change-status", h.changeStatus).Methods("PUT")
	router.HandleFunc(config.RootPath+"/applications/", h.saveApplication).Methods("POST")
	router.HandleFunc(config.RootPath+"/applications/get/filter-info", h.getFilterInfo).Methods("GET")

	return router
}

func (h *applicationHandler) logout(w http.ResponseWriter, r *http.Request) {
	if !h.UserService.AlreadyLoggedIn(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	c, _ := r.Cookie("session-id")

	err := h.UserService.DeleteSessionBySessionId(r.Context(), c.Value)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//remove cookie
	c = &http.Cookie{
		Name:   "session-id",
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(w, c)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	http.Redirect(w, r, "/logged-in", http.StatusSeeOther)
}

func (h *applicationHandler) loggedIn(w http.ResponseWriter, r *http.Request) {
	if h.UserService.AlreadyLoggedIn(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	var u *model.UserRegister
	err := util.DecodeBody(w, r, &u)
	if err != nil {
		return
	}

	user, err2 := h.UserService.GetUserByUsername(r.Context(), u.UserName)

	if user == nil {
		http.Error(w, "Username and/or password doesn't match", http.StatusForbidden)
		return
	}

	if err2 != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err3 := bcrypt.CompareHashAndPassword(user.Password, []byte(u.Password))

	if err3 != nil {
		http.Error(w, "Username and/or password doesn't match", http.StatusForbidden)
		return
	}

	id := uuid.NewV4()
	cookie := &http.Cookie{
		Name:  "session-id",
		Value: id.String(),
	}
	http.SetCookie(w, cookie)

	err = h.UserService.SaveSession(r.Context(), cookie.Value, user.Id)
	if err != nil {
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&user)
}

func (h *applicationHandler) bar(w http.ResponseWriter, r *http.Request) {
	if !h.UserService.AlreadyLoggedIn(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	cookie, _ := r.Cookie("session-id")

	userAddress, err := h.UserService.GetUserIfExist(r.Context(), cookie.Value)
	var user = *userAddress

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if userAddress != nil && user.Role != "Admin" {
		http.Error(w, "Tou don't have a permission to enter this", http.StatusUnauthorized)
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
	cookie, _ := r.Cookie("session-id")
	userAddress, err := h.UserService.GetUserIfExist(r.Context(), cookie.Value)
	var user = *userAddress

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if userAddress != nil && user.Role != "USER" {
		http.Error(w, "You don't have a permission to enter this", http.StatusUnauthorized)
		return
	}

	var application *model.Application
	err = util.DecodeBody(w, r, &application)
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
