package database

import (
	"net/http"

	"github.com/nutmeg-co/common-backend/claims"
	"github.com/nutmeg-co/common-backend/web"
	"gorm.io/gorm"
)

func Delete[T any](db *gorm.DB, sel ObjectSelector[T]) web.DataWithClaims {
	return func(w http.ResponseWriter, r *http.Request, claims *claims.Claims) (interface{}, error) {
		db := db
		obj := sel(r, claims)
		if err := db.Delete(obj).Error; err != nil {
			return nil, err
		}
		return true, nil
	}
}
