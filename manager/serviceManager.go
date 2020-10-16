package manager

import (
	"tweetlater/infra"
	usecase "tweetlater/usecases"
)

type ServiceManager interface {
	UserUseCase() usecase.IUserUseCase
	UserAuthUseCase() usecase.IUserAuthUseCase
	AppUseCase() usecase.IAppUseCase
}

type serviceManager struct {
	repo RepoManager
}

func (sm *serviceManager) UserUseCase() usecase.IUserUseCase {
	return usecase.NewUserUseCase(sm.repo.UserRepo())
}

func (sm *serviceManager) AppUseCase() usecase.IAppUseCase {
	return usecase.NewAppUseCase(sm.repo.AppRepo())
}

func (sm *serviceManager) UserAuthUseCase() usecase.IUserAuthUseCase {
	return usecase.NewUserAuthUseCase(sm.repo.UserAuthRepo(), sm.repo.UserRepo())
}
func NewServiceManger(infra infra.Infra) ServiceManager {
	return &serviceManager{repo: NewRepoManager(infra)}
}
