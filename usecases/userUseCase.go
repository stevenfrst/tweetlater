package usecase

import (
	"tweetlater/models"
	"tweetlater/repository"
)

type IUserUseCase interface {
	Register(newUser *models.User) error
	GetUserInfo(id string) *models.User
	Unregister(id string) error
	UpdateInfo(id string, newUser *models.User) error
	UpgradeAccount(user *models.User) error
}
type UserUseCase struct {
	userRepo repository.IUserRepository
}

func NewUserUseCase(userRepo repository.IUserRepository) IUserUseCase {
	return &UserUseCase{
		userRepo,
	}
}

func (uc *UserUseCase) Register(newUser *models.User) error {
	return uc.userRepo.Create(newUser)
}

func (uc *UserUseCase) GetUserInfo(id string) *models.User {
	user, err := uc.userRepo.FindOneById(id)
	if err != nil {
		return nil
	}
	return user
}

func (uc *UserUseCase) UpdateInfo(id string, newUser *models.User) error {
	return uc.userRepo.Update(id, newUser)
}

func (uc *UserUseCase) UpgradeAccount(user *models.User) error {
	return uc.userRepo.Upgrade(user)
}


func (uc *UserUseCase) Unregister(id string) error {
	return uc.userRepo.Delete(id)
}