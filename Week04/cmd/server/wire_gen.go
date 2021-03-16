package main

import (
	"Week04/internal/biz"
	"Week04/internal/data"
)

// Injectors from wire.go:

func InitUserUsecase() *biz.UserUsecase {
	userRepo := data.NewUserRepo()
	userUsecase := biz.NewUserUsecase(userRepo)
	return userUsecase
}
