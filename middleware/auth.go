package middleware

import (
	"context"
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
		ctx := context.WithValue(r.Context(), config.Session, claims)
		f(w, r.WithContext(ctx))
	}
}

/**
 * Same Logic With TokenDecode Function
 *
 * Not Abort to Next Handler At Tokenless State / Insert Guest Token
 */
func TokenParse(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := r.Cookie(config.Session)
		guestCtx := context.WithValue(
			r.Context(),
			config.Session,
			model.AuthTokenClaims{
				Id:       -1,
				Role:     config.RoleGuest,
				Name:     util.GenGuestName(),
				Username: "",
			},
		)
		if err != nil {
			f(w, r.WithContext(guestCtx))
			return
		}
		tok, claims, err := util.AuthToken(token.Value)
		if err != nil || !tok.Valid {
			f(w, r.WithContext(guestCtx))
			return
		}
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
		if ok && claims.Role == config.RoleAdmin {
			f(w, r)
			return
		}
		http.Error(w, "unahthorized token", http.StatusUnauthorized)
	}
}
