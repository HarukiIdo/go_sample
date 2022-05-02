package middleware

import (
	"log"
	"net/http"
)

func MiddlewareLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("start%s\n", r.URL)
		next.ServeHTTP(w, r)
		log.Printf("finish%s\n", r.URL)
	})
}
