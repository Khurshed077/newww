package jwt

import "context"

func SetClaimsToContext(ctx context.Context, claims *Claims) context.Context {
	return context.WithValue(ctx, ClaimsKey, claims)
}
