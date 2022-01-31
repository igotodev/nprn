package handler

import (
	"context"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"nprn/internal/customerr"
	"nprn/internal/entity/sale/salemodel"
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

type answer struct {
	ID string `json:"id"`
}

func NewHandler(service *service.Service, logger *logging.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

func (h *Handler) RegisterRouting(router *httprouter.Router) {

	router.GET("/auth/sign-in", h.CheckErrorMiddleware(h.SignIn))
	router.POST("/auth/sign-up", h.CheckErrorMiddleware(h.SignUp))
	{
		//router.PUT("/user/:id", h.CheckAuthorizationMiddleware(h.Update))
		//router.DELETE("/user/:id", h.CheckAuthorizationMiddleware(h.Delete))
	}

	{
		router.GET("/api/v1/sale/", h.CheckAuthorizationMiddleware(h.GetAllSales))
		router.GET("/api/v1/sale/:id", h.CheckAuthorizationMiddleware(h.GetSale))
		router.POST("/api/v1/sale/", h.CheckAuthorizationMiddleware(h.CreateSale))
		router.PUT("/api/v1/sale/:id", h.CheckAuthorizationMiddleware(h.UpdateSale))
		router.DELETE("/api/v1/sale/:id", h.CheckAuthorizationMiddleware(h.DeleteSale))
	}

	h.logger.Info("routing is registered")
}

func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request, _ httprouter.Params) error {
	var signReq signInRequest

	err := json.NewDecoder(r.Body).Decode(&signReq)
	if err != nil {
		return customerr.NewCustomError(err, "error with decode body")
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
		return customerr.NewCustomError(err, "error with marshal json answer")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(marshal)

	return nil
}

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request, _ httprouter.Params) error {
	var usr usermodel.UserInternal

	err := json.NewDecoder(r.Body).Decode(&usr)
	if err != nil {
		return customerr.NewCustomError(err, "error with decode body")
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
		return customerr.NewCustomError(err, "error with marshal json answer")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(marshal)

	return nil
}

func (h *Handler) CreateSale(w http.ResponseWriter, r *http.Request, _ httprouter.Params) error {
	var sale salemodel.Sale

	err := json.NewDecoder(r.Body).Decode(&sale)
	if err != nil {
		return customerr.NewCustomError(err, "error with decode body")
	}

	id, err := h.service.CreateSale(context.Background(), sale)
	if err != nil {
		h.logger.Info(err)
		return err
	}

	marshal, err := json.Marshal(&answer{ID: id})
	if err != nil {
		return customerr.NewCustomError(err, "error with marshal json answer")
	}
	w.WriteHeader(200)
	w.Write(marshal)

	return nil
}

func (h *Handler) GetSale(w http.ResponseWriter, _ *http.Request, params httprouter.Params) error {
	idStr := params.ByName("id")

	result, err := h.service.GetSale(context.Background(), idStr)
	if err != nil {
		h.logger.Info(err)
		return err
	}

	marshal, err := json.Marshal(result)
	if err != nil {
		return customerr.NewCustomError(err, "error with marshal json answer")
	}

	w.WriteHeader(200)
	w.Write(marshal)

	return nil
}

func (h *Handler) GetAllSales(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) error {

	result, err := h.service.GetAllSales(context.Background())
	if err != nil {
		h.logger.Info(err)
		return err
	}

	marshal, err := json.Marshal(result)
	if err != nil {
		return customerr.NewCustomError(err, "error with marshal json answer")
	}

	w.WriteHeader(200)
	w.Write(marshal)

	return nil
}

func (h *Handler) UpdateSale(w http.ResponseWriter, r *http.Request, params httprouter.Params) error {
	idStr := params.ByName("id")

	var saleUpdate salemodel.Sale

	err := json.NewDecoder(r.Body).Decode(&saleUpdate)
	if err != nil {
		return customerr.NewCustomError(err, "error with decode body")
	}

	saleUpdate.ID = idStr

	err = h.service.UpdateSale(context.Background(), saleUpdate)
	if err != nil {
		h.logger.Info(err)
		return err
	}

	marshal, err := json.Marshal(&answer{ID: saleUpdate.ID})
	if err != nil {
		return customerr.NewCustomError(err, "error with marshal json answer")
	}
	w.WriteHeader(200)
	w.Write(marshal)

	return nil
}

func (h *Handler) DeleteSale(w http.ResponseWriter, _ *http.Request, params httprouter.Params) error {
	idStr := params.ByName("id")

	err := h.service.DeleteSale(context.Background(), idStr)
	if err != nil {
		h.logger.Info(err)
		return err
	}

	marshal, err := json.Marshal(&answer{ID: idStr})
	if err != nil {
		return customerr.NewCustomError(err, "error with marshal json answer")
	}
	w.WriteHeader(200)
	w.Write(marshal)

	return nil
}

//func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) error {
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
//	err = h.service.UpdateUser(context.Background(), usr)
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
//func (h *Handler) DeleteUser(w http.ResponseWriter, r *http.Request, params httprouter.Params) error {
//	idStr := params.ByName("id")
//
//	err := h.service.DeleteUser(context.Background(), idStr)
//	if err != nil {
//		h.logger.Info(err)
//		return err
//	}
//
//	w.WriteHeader(204)
//
//	return nil
//}
