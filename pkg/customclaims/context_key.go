package customclaims

type ContextKey string

var (
	ContextKeyUserId = ContextKey("authful_user_id")
	ContextJwt       = ContextKey("authful_jwt")
	ContextTraceId   = ContextKey("authful_trace_id")
)
