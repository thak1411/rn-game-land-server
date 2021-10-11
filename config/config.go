package config

import "time"

const (
	// Server Port //
	Port = ":8192"

	// WebSocket //
	WriteWait      = 10 * time.Second
	PongWait       = 60 * time.Second
	PingPeriod     = (PongWait * 9) / 10
	MaxMessageSize = 512

	// JWT //
	JwtSecretKey = "hi_this_is_secret_key_for_rn_jwt"
	Session      = "RN_SESSION"

	// User Role //
	RoleAdmin = "RN_ROLE_ADMIN"
	RoleBasic = "RN_ROLE_BASIC"
)
