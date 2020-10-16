package appMiddleware

import (
	"github.com/gorilla/sessions"
	"net/http"
)

type TokenValidationMiddleware struct {
	Store *sessions.CookieStore
}

func NewTokenValidationMiddleware(cs *sessions.CookieStore) *TokenValidationMiddleware {
	return &TokenValidationMiddleware{cs}
}

func (v *TokenValidationMiddleware) Validate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, _ := v.Store.Get(r, "app-cookie")
		if isAuth, ok := session.Values["authenticated"].(bool); !isAuth || !ok {
			http.Error(w, "", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (v *TokenValidationMiddleware) GetCookieStore() *sessions.CookieStore {
	return v.Store
}
