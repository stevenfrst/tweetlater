package usecase

import (
	"tweetlater/models"
	"tweetlater/repository"
)



type IUserAuthUseCase interface {
	UserNamePasswordValidation(userName string, password string) *models.User
}

type UserAuthUseCase struct {
	userAuthRepo repository.IUserAuthRepository
	userRepo     repository.IUserRepository
}

func NewUserAuthUseCase(userAuthRepo repository.IUserAuthRepository, userRepo repository.IUserRepository) IUserAuthUseCase {
	return &UserAuthUseCase{
		userAuthRepo, userRepo,
	}
}

func (uc *UserAuthUseCase) UserNamePasswordValidation(userName string, password string) *models.User {
	userAuth := uc.userAuthRepo.FindOneByUserNameAndPassword(userName, password)
	if userAuth != nil {
		userInfo, err := uc.userRepo.FindOneById(userAuth.Id)
		if err != nil {
			return nil
		}
		return userInfo
	}
	return nil
}
