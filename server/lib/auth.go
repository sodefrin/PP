package lib

import (
	"context"
	"errors"

	"github.com/sodefrin/PP/server/db"
)

func SetUserContext(ctx context.Context, user db.User) context.Context {
	return context.WithValue(ctx, userContextKey, user)
}

func GetUserContext(ctx context.Context) (db.User, error) {
	user, ok := ctx.Value(userContextKey).(db.User)
	if !ok {
		return db.User{}, errors.New("user not found")
	}
	return user, nil
}
