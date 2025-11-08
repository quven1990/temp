package service

import "context"

type IUser interface {
	Login(ctx context.Context, username, password string) (string, error)
}

var User IUser
