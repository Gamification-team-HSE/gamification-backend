package middlewares

import (
	"fmt"
	"gitlab.com/krespix/gamification-api/internal/services/auth"
	"gitlab.com/krespix/gamification-api/pkg/utils"
	"net/http"
)

type Auth struct {
	authService auth.Service
}

func (a *Auth) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			next.ServeHTTP(w, r)
			return
		}

		bearer := "Bearer "
		token := authHeader[len(bearer):]

		claims, err := a.authService.ValidateToken(r.Context(), token)
		if err != nil {
			fmt.Println(err)
			http.Error(w, fmt.Sprintf("Invalid token: %v", err), http.StatusUnauthorized)
			return
		}

		ctx := utils.PutClaimsToContext(r.Context(), claims)

		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func NewAuth(service auth.Service) *Auth {
	return &Auth{authService: service}
}
