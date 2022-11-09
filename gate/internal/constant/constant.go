package constant

import "errors"

const (
	SessionAttrKey                   = "SESSION_KEY"
	AuthenticateTimeout              = 30
	AuthenticateTokenDiff            = 3600
	AuthenticateTokenRefreshInterval = 600
)

var (
	ErrWrongMessageAssert = errors.New("message assert failure")
	ErrWrongAuthFailure = errors.New("session auth failure")
)