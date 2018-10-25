package api

import (
	"io"
	"net/http"
)

func wrapper(f func(io.Writer, *http.Request) (interface{}, int, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, status, err := f(w, r)
		if err != nil {
			http.Error(w, err.Error(), status)
			return
		}

		w.Header().Set("Content-Type", " application/json")
		w.WriteHeader(status)

		_, err = io.WriteString(w, data.(string))
		if err != nil {
			panic(err)
		}
	}
}
