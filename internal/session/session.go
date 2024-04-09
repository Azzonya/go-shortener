package session

import (
	"context"
	"errors"

	"github.com/Azzonya/go-shortener/internal/user"
)

type ctxKey int

const (
	ctxKeyUID ctxKey = iota
)

func SetUserContext(parent context.Context, u *user.User) context.Context {
	return context.WithValue(parent, ctxKeyUID, u)
}

func GetUserFromContext(ctx context.Context) (u *user.User, ok bool) {
	u, ok = ctx.Value(ctxKeyUID).(*user.User)
	return
}

func GetUser(c context.Context) (*user.User, error) {
	u, ok := GetUserFromContext(c)
	if !ok {
		return nil, errors.New("failed to get user from context")
	}

	if u.IsNew() {

		return nil, errors.New("no authorized")
	}

	return u, nil
}
