package middleware

import (
	"log"
	"net/http"
)

// func someHandler(w http.ResponseWriter, r *http.Request) {
// 	log.Println("from some handlers")
// 	w.WriteHeader(http.StatusOK)
// }

func SomeMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("some logig from middleware")
		next(w, r)
	}
}
