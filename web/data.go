package web

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/nutmeg-co/common-backend/errors"
	"gorm.io/gorm"
)

type DataHandler func(w http.ResponseWriter, r *http.Request) (interface{}, error)

func (h DataHandler) Json() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		result, err := h(w, r)
		if err == nil {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(result)
			return
		}
		wrapped, ok := w.(*accessLogWrapper)
		if ok {
			wrapped.error = err.Error()
		}
		switch err {
		default:
			w.WriteHeader(http.StatusBadRequest)
		case errors.ErrForbidden:
			w.WriteHeader(http.StatusForbidden)
		case errors.ErrUnauthorized:
			w.WriteHeader(http.StatusUnauthorized)
		case gorm.ErrRecordNotFound:
			w.WriteHeader(http.StatusNotFound)
		case errors.ErrWebsocket:
			// do nothing
			return
		}
		fmt.Fprint(w, err.Error())
	}
}
