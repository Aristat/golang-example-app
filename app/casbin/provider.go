package casbin

import (
	"github.com/aristat/golang-example-app/app/entrypoint"
	"github.com/casbin/casbin"
	fileadapter "github.com/casbin/casbin/persist/file-adapter"
	"github.com/google/wire"
)

func Provider() (*casbin.Enforcer, func(), error) {
	wd := entrypoint.WorkDir()
	enf := casbin.NewEnforcer(wd+"/casbin/model.conf", wd+"/casbin/policy.csv")
	return enf, func() {}, nil
}

func ProviderTest() (*casbin.Enforcer, func(), error) {

	model := `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act
`
	policy := `
p, reader, users, read
p, owner, users, write

g, reader, anonymous
g, owner, reader
`
	enf := casbin.NewEnforcer(casbin.NewModel(model), fileadapter.NewAdapter(policy))
	return enf, func() {}, nil
}

var (
	ProviderProductionSet = wire.NewSet(Provider)
	ProviderTestSet       = wire.NewSet(ProviderTest)
)
