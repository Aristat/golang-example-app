package casbin

import (
	"github.com/aristat/golang-example-app/app/entrypoint"
	"github.com/casbin/casbin"
	"github.com/google/wire"
)

func Provider() (*casbin.Enforcer, func(), error) {
	wd := entrypoint.WorkDir()
	enf := casbin.NewEnforcer(wd+"/casbin/model.conf", wd+"/casbin/policy.csv")
	return enf, func() {}, nil
}

func ProviderTest() (*casbin.Enforcer, func(), error) {
	wd := entrypoint.WorkDir()
	enf := casbin.NewEnforcer(wd+"/casbin/model.conf", wd+"/casbin/policy.csv")
	return enf, func() {}, nil
}

var (
	ProviderProductionSet = wire.NewSet(Provider)
	ProviderTestSet       = wire.NewSet(ProviderTest)
)
