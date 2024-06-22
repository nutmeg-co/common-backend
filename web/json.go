package web

import (
	"encoding/json"
	"net/http"
)

func ParseJson(w http.ResponseWriter, r *http.Request, data interface{}) error {
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(data); err != nil {
		return err
	}
	wrapped, ok := w.(*accessLogWrapper)
	if ok {
		wrapped.data = data
	}
	return nil
}
