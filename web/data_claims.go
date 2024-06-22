package web

import (
	"net/http"

	"github.com/nutmeg-co/common-backend/claims"
	"github.com/nutmeg-co/common-backend/errors"
)

type DataWithClaims func(w http.ResponseWriter, r *http.Request, claims *claims.Claims) (interface{}, error)

func (h DataWithClaims) Handler() DataHandler {
	return func(w http.ResponseWriter, r *http.Request) (interface{}, error) {
		claims, err := claims.ClaimsFactory.GetClaims(r)
		if err != nil {
			return nil, errors.ErrUnauthorized
		}
		wrapped, ok := w.(*accessLogWrapper)
		if ok {
			wrapped.claims = claims
		}
		return h(w, r, claims)
	}
}

func (h DataWithClaims) Unauthorize() DataWithClaims {
	return func(w http.ResponseWriter, r *http.Request, c *claims.Claims) (interface{}, error) {
		result, err := h(w, r, c)
		if err != nil {
			return nil, err
		}
		claims.Unauthorize(w)
		return result, nil
	}
}

func (h DataWithClaims) RequireRole(roles ...string) DataWithClaims {
	return func(w http.ResponseWriter, r *http.Request, claims *claims.Claims) (interface{}, error) {
		if claims.Role == nil {
			return nil, errors.ErrForbidden
		}
		if len(roles) == 0 {
			return h(w, r, claims)
		}
		rr := *claims.Role
		for _, role := range roles {
			if role == rr {
				return h(w, r, claims)
			}
		}
		return nil, errors.ErrForbidden
	}
}
