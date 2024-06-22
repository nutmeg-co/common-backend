package claims

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/nutmeg-co/common-backend/errors"
)

type Factory struct {
	Secret []byte
}

func (factory *Factory) Token(claims *Claims) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(factory.Secret))
}

func (factory *Factory) GetClaims(r *http.Request) (*Claims, error) {
	cookie, err := r.Cookie("token")
	if err != nil {
		return nil, err
	}
	if err := cookie.Valid(); err != nil {
		return nil, err
	}
	var claims Claims
	token, err := jwt.ParseWithClaims(cookie.Value, &claims, func(t *jwt.Token) (interface{}, error) {
		return factory.Secret, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.ErrInvalidToken
	}
	return &claims, nil
}
