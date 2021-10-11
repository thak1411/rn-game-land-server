package model

import "github.com/dgrijalva/jwt-go/v4"

/**
 * JWT Token Claims
 *
 * Uuid: Auto Injected - Random UUID String
 * Name: User(Player) Display Name
 * Username: User(Player) ID For LOG-IN
 * Password: User(Player) Pass For LOG-IN
 */
type AuthTokenClaims struct {
	Id       int    `json:"id"`
	Uuid     string `json:"uuid"`
	Role     string `json:"role"`
	Name     string `json:"name"`
	Username string `json:"username"`
	jwt.StandardClaims
}
