package api

import (
	"crypto/subtle"
	"net/http"
)

const authHeader = "X-API-Key"

// apiKeyAuth returns a middleware that enforces a static API key check.
// If the configured key is empty, authentication is skipped.
func apiKeyAuth(key string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if key == "" {
				next.ServeHTTP(w, r)
				return
			}

			provided := r.Header.Get(authHeader)
			if provided == "" {
				writeJSON(w, http.StatusUnauthorized, map[string]string{
					"error": "missing API key",
				})
				return
			}

			if subtle.ConstantTimeCompare([]byte(provided), []byte(key)) != 1 {
				writeJSON(w, http.StatusForbidden, map[string]string{
					"error": "invalid API key",
				})
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
