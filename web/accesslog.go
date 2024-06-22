package web

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/nutmeg-co/common-backend/claims"
)

type accessLogWrapper struct {
	http.ResponseWriter
	code   int
	log    bool
	r      *http.Request
	claims *claims.Claims
	data   interface{}
	error  string
}

func AccessLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		wrapped := &accessLogWrapper{
			ResponseWriter: w,
			code:           200,
			log:            true,
			r:              r,
		}
		next.ServeHTTP(wrapped, r)
		if !wrapped.log {
			return
		}
		if wrapped.code < 400 {
			if wrapped.claims == nil {
				log.Printf(
					"%d %s %s",
					wrapped.code, r.Method, r.URL.Path,
				)
			} else if wrapped.claims.OutletID == nil {
				log.Printf(
					"%d %s %s (account_id:%s)",
					wrapped.code, r.Method, r.URL.Path,
					wrapped.claims.Subject,
				)
			} else {
				log.Printf(
					"%d %s %s (account_id:%s outlet_id:%s)",
					wrapped.code, r.Method, r.URL.Path,
					wrapped.claims.Subject,
					*wrapped.claims.OutletID,
				)
			}
			return
		}
		claims, _ := json.Marshal(wrapped.claims)
		data, _ := json.Marshal(wrapped.data)
		log.Printf(
			"%d %s %s\nclaims: %s\ndata:   %s\nerror:  %s",
			wrapped.code, r.Method, r.URL.Path,
			claims, data, wrapped.error,
		)
	})
}

func (w *accessLogWrapper) WriteHeader(code int) {
	w.code = code
	w.ResponseWriter.WriteHeader(code)
}

func (c *accessLogWrapper) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	h, ok := c.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, fmt.Errorf("not supported")
	}
	conn, buf, err := h.Hijack()
	if err == nil {
		c.log = false
		log.Printf("HIJACK %s %s", c.r.Method, c.r.URL.Path)
	}
	return conn, buf, err
}
