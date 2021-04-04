package accounts

import (
	"context"
)

const (
	userKey privateKey = "user"
)

type privateKey string

func WithUserContext(ctx context.Context, user *UserModel) context.Context {
	return context.WithValue(ctx, userKey, user)
}

func UserContext(ctx context.Context) *UserModel {
	if tmp := ctx.Value(userKey); tmp != nil {
		if user, ok := tmp.(*UserModel); ok {
			return user
		}
	}
	return nil
}
