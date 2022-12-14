package middlewares

import (
	"gitlab.com/krespix/gamification-api/internal/models"
	"gitlab.com/krespix/gamification-api/internal/services/auth"
	"gitlab.com/krespix/gamification-api/pkg/utils"
	"net/http"
	"strconv"
)

type Auth struct {
	authService    auth.Service
	fakeAuth       bool
	allowedMethods string
	allowedHeaders string
}

func (a *Auth) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var claims *models.Claims
		claims, err := a.tryGetClaims(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		if a.fakeAuth && claims == nil {
			claims = a.tryGetFakeClaims(r)
		}
		if claims == nil {
			next.ServeHTTP(w, r)
			return
		}

		ctx := utils.PutClaimsToContext(r.Context(), claims)

		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func (a *Auth) tryGetClaims(r *http.Request) (*models.Claims, error) {
	authHeader := r.Header.Get("Authorization")

	if authHeader == "" {
		return nil, nil
	}

	bearer := "Bearer "
	token := authHeader[len(bearer):]

	claims, err := a.authService.ValidateToken(r.Context(), token)
	if err != nil {
		return nil, err
	}
	return claims, nil
}

func (a *Auth) tryGetFakeClaims(r *http.Request) *models.Claims {
	userID := r.Header.Get("X-Auth-User-ID")
	if userID == "" {
		return nil
	}
	role := r.Header.Get("X-Auth-Role")
	if role == "" {
		return nil
	}

	id, err := strconv.Atoi(userID)
	if err != nil {
		return nil
	}

	return &models.Claims{
		ID:   int64(id),
		Role: models.Role(role),
	}
}

func (a *Auth) EnableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", a.allowedMethods)
		w.Header().Set("Access-Control-Allow-Headers", a.allowedHeaders)
		next.ServeHTTP(w, r)
	})
}

func NewAuth(service auth.Service, allowedHeaders, allowedMethods string, fakeAuth bool) *Auth {
	return &Auth{
		authService:    service,
		allowedMethods: allowedMethods,
		allowedHeaders: allowedHeaders,
		fakeAuth:       fakeAuth,
	}
}
