package database

import (
	"net/http"

	"github.com/nutmeg-co/common-backend/claims"
	"github.com/nutmeg-co/common-backend/web"
	"gorm.io/gorm"
)

type ObjectSelector[T any] func(r *http.Request, claims *claims.Claims) *T

func First[T any](db *gorm.DB, sel ObjectSelector[T], filters ...Filter) web.DataWithClaims {
	return func(w http.ResponseWriter, r *http.Request, claims *claims.Claims) (interface{}, error) {
		obj := sel(r, claims)
		db := db.Where(obj)
		var err error
		q := r.URL.Query()
		for _, f := range filters {
			db, err = f(db, r, claims, q)
			if err != nil {
				return nil, err
			}
		}
		var out T
		if err := db.First(&out).Error; err != nil {
			return nil, err
		}
		return &out, nil
	}
}

func Find[T any](db *gorm.DB, sel ObjectSelector[T], filters ...Filter) web.DataWithClaims {
	return func(w http.ResponseWriter, r *http.Request, claims *claims.Claims) (interface{}, error) {
		obj := sel(r, claims)
		db := db.Where(obj)
		var err error
		q := r.URL.Query()
		for _, f := range filters {
			db, err = f(db, r, claims, q)
			if err != nil {
				return nil, err
			}
		}
		var out []*T
		if err := db.Find(&out).Error; err != nil {
			return nil, err
		}
		return &out, nil
	}
}
