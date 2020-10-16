package manager

import (
	"tweetlater/infra"
	"tweetlater/repository"
)

type RepoManager interface {
	UserRepo() repository.IUserRepository
	UserAuthRepo() repository.IUserAuthRepository
	AppRepo() repository.IAppRepository
}
type repoManager struct {
	infra infra.Infra
}

func (rm *repoManager) UserRepo() repository.IUserRepository {
	return repository.NewUserRepository(rm.infra.SqlDb())
}

func (rm *repoManager) AppRepo() repository.IAppRepository {
	return repository.NewAppRepository(rm.infra.SqlDb())
}

func (rm *repoManager) UserAuthRepo() repository.IUserAuthRepository {
	return repository.NewUserAuthRepository(rm.infra.SqlDb())
}

func NewRepoManager(infra infra.Infra) RepoManager {
	return &repoManager{infra}
}
