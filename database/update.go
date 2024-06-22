package database

import (
	"net/http"

	"github.com/nutmeg-co/common-backend/claims"
	"github.com/nutmeg-co/common-backend/web"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ObjectPatcher[T any] func(w http.ResponseWriter, r *http.Request, claims *claims.Claims, data *T) (*T, error)

func Update[T any](db *gorm.DB, sel ObjectSelector[T], patch ObjectPatcher[T]) web.DataWithClaims {
	return func(w http.ResponseWriter, r *http.Request, claims *claims.Claims) (interface{}, error) {
		var (
			out   T
			data  T
			where = sel(r, claims)
		)
		if err := web.ParseJson(w, r, &data); err != nil {
			return nil, err
		}
		patch, err := patch(w, r, claims, &data)
		if err != nil {
			return nil, err
		}
		if err := db.Model(&out).Where(where).Clauses(clause.Returning{}).Updates(patch).Error; err != nil {
			return nil, err
		}
		return &out, nil
	}
}
