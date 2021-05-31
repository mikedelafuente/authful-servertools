package customclaims

type ContextKey string

var (
	ContextKeyUserId = ContextKey("userId")
	ContextJwt       = ContextKey("jwt")
)
