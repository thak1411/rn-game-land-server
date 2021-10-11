package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/thak1411/rn-game-land-server/config"
	"github.com/thak1411/rn-game-land-server/model"
	"github.com/thak1411/rn-game-land-server/util"
)

func TokenDecode(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := r.Cookie(config.Session)
		if err != nil {
			http.Error(w, "unahthorized token", http.StatusUnauthorized)
			return
		}
		tok, claims, err := util.AuthToken(token.Value)
		if err != nil || !tok.Valid {
			http.Error(w, "unahthorized token", http.StatusUnauthorized)
			return
		}
		fmt.Println("ABC", claims)
		ctx := context.WithValue(r.Context(), config.Session, claims)
		f(w, r.WithContext(ctx))
	}
}

func AuthAdmin(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		iClaims := r.Context().Value(config.Session)
		if iClaims == nil {
			http.Error(w, "unahthorized token", http.StatusUnauthorized)
			return
		}
		claims, ok := iClaims.(model.AuthTokenClaims)
		fmt.Println("CBA", iClaims, claims)
		if ok && claims.Role == config.RoleAdmin {
			f(w, r)
			return
		}
		http.Error(w, "unahthorized token", http.StatusUnauthorized)
	}
}
