package repo

import (
	"github.com/aristat/golang-example-app/app/db/domain"
	"github.com/google/wire"
)

// Repo for all records
type Repo struct {
	Users domain.UsersRepo
}

// Provider
func Provider(userDomain domain.UsersRepo) (*Repo, func(), error) {
	repo := &Repo{
		Users: userDomain,
	}

	return repo, func() {}, nil
}

var (
	ProviderProductionSet = wire.NewSet(Provider, NewUsersRepo)
	ProviderTestSet       = wire.NewSet(Provider, NewUsersRepo)
)
