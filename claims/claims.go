package claims

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lucsky/cuid"
)

type Claims struct {
	jwt.RegisteredClaims
	OutletID *string `json:"outlet_id,omitempty"`
	Role     *string `json:"role,omitempty"`
}

func NewClaims(accountID string) *Claims {
	now := time.Now()
	return &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        cuid.New(),
			Subject:   accountID,
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(LoginExpiration)),
		},
	}
}

func NewClaimsWithScope(accountID, role, outletID string) *Claims {
	now := time.Now()
	return &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        cuid.New(),
			Subject:   accountID,
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(LoginExpiration)),
		},
		OutletID: &outletID,
		Role:     &role,
	}
}

func (claims *Claims) Authorize(w http.ResponseWriter) error {
	token, err := ClaimsFactory.Token(claims)
	if err != nil {
		return err
	}
	http.SetCookie(w, &http.Cookie{
		Path:    "/",
		Name:    "token",
		Value:   token,
		Expires: claims.ExpiresAt.Time,
	})
	return nil
}

func Unauthorize(w http.ResponseWriter) error {
	http.SetCookie(w, &http.Cookie{
		Path:  "/",
		Name:  "token",
		Value: "",
	})
	return nil
}
