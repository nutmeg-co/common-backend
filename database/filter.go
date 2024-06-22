package database

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/nutmeg-co/common-backend/claims"
	"gorm.io/gorm"
)

type Filter func(db *gorm.DB, r *http.Request, claims *claims.Claims, q url.Values) (*gorm.DB, error)

func AllowPreload(db *gorm.DB, r *http.Request, claims *claims.Claims, q url.Values) (*gorm.DB, error) {
	if q.Has("preload") {
		for _, preload := range strings.Split(q.Get("preload"), ",") {
			db = db.Preload(preload, func(db *gorm.DB) *gorm.DB {
				return db.Unscoped()
			})
		}
	}
	return db, nil
}

func AllowTimeFilter(db *gorm.DB, r *http.Request, claims *claims.Claims, q url.Values) (*gorm.DB, error) {
	for k, v := range r.URL.Query() {
		switch k {
		case "start":
			db = db.Where("created_at >= ?", v)
		case "end":
			db = db.Where("created_at < ?", v)
		}
	}
	return db, nil
}

func AllowQueryFilter(list ...string) Filter {
	return func(db *gorm.DB, r *http.Request, claims *claims.Claims, q url.Values) (*gorm.DB, error) {
		for _, k := range list {
			if q.Has(k) {
				db = db.Where(k+" = ?", q.Get(k))
			}
		}
		return db, nil
	}
}

func AllowUnscoped(db *gorm.DB, r *http.Request, claims *claims.Claims, q url.Values) (*gorm.DB, error) {
	if q.Has("unscoped") {
		db = db.Unscoped()
	}
	return db, nil
}
