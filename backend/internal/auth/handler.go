package auth

import (
	"net/http"
	"strings"
)

func Middleware(service *Service, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")

		if strings.HasPrefix(header, "Bearer ") {
			token := strings.TrimPrefix(header, "Bearer ")

			if userID, err := service.Parse(token); err == nil {
				r = r.WithContext(WithUserID(r.Context(), userID))
			}
		}

		next.ServeHTTP(w, r)
	})
}
