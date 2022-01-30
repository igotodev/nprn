package handler

import (
	"context"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"net/http"
	"nprn/internal/entity/user/usermodel"
	"nprn/internal/service"
	"nprn/pkg/logging"
)

type Handler struct {
	service *service.Service
	logger  *logging.Logger
}

type signInRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type tokenResponse struct {
	Token string `json:"token"`
}

func NewHandler(service *service.Service, logger *logging.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

func (h *Handler) RegisterRouting(router *httprouter.Router) {
	{
		router.POST("/auth/sign-in", h.CheckErrorMiddleware(h.SignIn))
		router.POST("/auth/sign-up", h.CheckErrorMiddleware(h.SignUp))

		{
			router.GET("/user/", h.CheckAuthorizationMiddleware(h.Test))
			//router.PUT("/user/:id", middleware.CheckAuthorizationMiddleware(h.Update))
			//router.DELETE("/user/:id", middleware.CheckAuthorizationMiddleware(h.Delete))
		}

		{
			router.GET("/api/v1/sale/", nil)
			router.GET("/api/v1/sale/:id", nil)
			router.POST("/api/v1/sale/", nil)
			router.PUT("/api/v1/sale/:id", nil)
			router.DELETE("/api/v1/sale/:id", nil)
		}
	}

	h.logger.Info("routing is registered")
}

func (h *Handler) Test(w http.ResponseWriter, r *http.Request, params httprouter.Params) error {
	w.WriteHeader(200)
	w.Write([]byte("hello world"))
	return nil
}

func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request, params httprouter.Params) error {

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		h.logger.Info(err)
		return err
	}

	defer r.Body.Close()

	var signReq signInRequest

	err = json.Unmarshal(body, &signReq)
	if err != nil {
		h.logger.Info(err)
	}

	token, err := h.service.SignIn(context.Background(), signReq.Username, signReq.Password)
	if err != nil {
		h.logger.Info(err)
		return err
	}

	tr := tokenResponse{
		Token: token,
	}

	marshal, err := json.Marshal(tr)
	if err != nil {
		h.logger.Info(err)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(marshal)

	return nil
}

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request, _ httprouter.Params) error {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		h.logger.Info(err)
	}

	defer r.Body.Close()

	var usr usermodel.UserInternal

	err = json.Unmarshal(body, &usr)
	if err != nil {
		h.logger.Info(err)
	}

	token, err := h.service.SignUp(context.Background(), usr)
	if err != nil {
		h.logger.Info(err)
		return err
	}

	tr := tokenResponse{
		Token: token,
	}

	marshal, err := json.Marshal(tr)
	if err != nil {
		h.logger.Info(err)
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(marshal)

	return nil
}

//func (h *Handler) Update(w http.ResponseWriter, r *http.Request, params httprouter.Params) error {
//
//	by, err := ioutil.ReadAll(r.Body)
//	if err != nil {
//		h.logger.Info(err)
//	}
//
//	var usr usermodel.UserInternal
//
//	err = json.Unmarshal(by, &usr)
//	if err != nil {
//		h.logger.Info(err)
//	}
//
//	idStr := params.ByName("id")
//	usr.ID = idStr
//
//	err = h.service.Update(context.Background(), usr)
//	if err != nil {
//		h.logger.Info(err)
//		return err
//	}
//
//	w.WriteHeader(204)
//
//	return nil
//}
//
//func (h *Handler) Delete(w http.ResponseWriter, r *http.Request, params httprouter.Params) error {
//	idStr := params.ByName("id")
//
//	err := h.service.Delete(context.Background(), idStr)
//	if err != nil {
//		h.logger.Info(err)
//		return err
//	}
//
//	w.WriteHeader(204)
//
//	return nil
//}
