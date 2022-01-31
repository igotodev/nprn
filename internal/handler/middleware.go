package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"nprn/internal/customerr"
	"strings"
)

type authError struct {
	Message string `json:"message"`
}

type CustomHandlerFunc func(w http.ResponseWriter, r *http.Request, params httprouter.Params) error

func (h *Handler) CheckAuthorizationMiddleware(handlerFunc CustomHandlerFunc) httprouter.Handle {

	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")

		header := r.Header.Get("Authorization")
		if len(header) == 0 {

			authErr := authError{
				Message: "unauthorized: header is empty",
			}

			marshal, err := json.Marshal(authErr)
			if err != nil {
				return
			}

			w.WriteHeader(401)
			w.Write(marshal)
			return
		}

		parts := strings.Split(header, " ")

		if len(parts) != 2 || parts[0] != "Bearer" || len(parts[1]) == 0 {
			authErr := authError{
				Message: "unauthorized",
			}

			marshal, err := json.Marshal(authErr)
			if err != nil {
				return
			}

			w.WriteHeader(401)
			w.Write(marshal)
			return
		}

		userID, err := h.service.ParseToken(parts[1])

		if err != nil {
			authErr := authError{
				Message: "unauthorized: not valid token",
			}

			marshal, err := json.Marshal(authErr)
			if err != nil {
				return
			}

			w.WriteHeader(401)
			w.Write(marshal)
			return
		}

		err = handlerFunc(w, r, params)
		if err != nil {
			h.logger.Info(err)
			checkCustomError(err, w)
			return
		}

		h.logger.Logger.Info(fmt.Sprintf("user with id=%v is accepted", userID))
	}
}

func (h *Handler) CheckErrorMiddleware(handlerFunc CustomHandlerFunc) httprouter.Handle {

	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")

		err := handlerFunc(w, r, params)
		if err != nil {
			h.logger.Info(err)
			checkCustomError(err, w)
		}

	}
}

func checkCustomError(err error, w http.ResponseWriter) {
	var customErr *customerr.CustomError
	if err != nil {
		if errors.As(err, &customErr) {
			if errors.Is(err, customerr.NotFoundErr) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(404)

				ce := err.(*customerr.CustomError)
				w.Write(ce.Marshal())

			} else if errors.Is(err, customerr.NotAcceptable) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(406)

				ce := err.(*customerr.CustomError)
				w.Write(ce.Marshal())

			} else {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(418)

				ce := err.(*customerr.CustomError)
				w.Write(ce.Marshal())
			}
		}
	}
}
