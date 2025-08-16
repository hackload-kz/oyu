package entity

import "context"

type userCtxKeyType string

const UserCtxKey userCtxKeyType = "auth.user"

func UserFromContext(ctx context.Context) (*AuthUser, bool) {
	u, ok := ctx.Value(UserCtxKey).(*AuthUser)
	return u, ok
}
