package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("senhaforte123")

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			http.Error(w, "Token ausente ou inválido", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(auth, "Bearer ")

		claims := &jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Token inválido", http.StatusUnauthorized)
			return
		}

		// Você pode adicionar informações ao contexto, se quiser
		ctx := context.WithValue(r.Context(), "usuario", (*claims)["cpf"])
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

