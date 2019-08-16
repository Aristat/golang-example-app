package config

import (
	"github.com/aristat/golang-oauth2-example-app/app/entrypoint"
	"github.com/google/wire"
	"github.com/spf13/viper"
)

// Provider returns spf13/viper instance with resolved dependencies
func Provider() (*viper.Viper, func(), error) {
	v := entrypoint.Viper()
	for _, key := range v.AllKeys() {
		val := v.Get(key)
		v.Set(key, val)
	}
	return v, func() {}, nil
}

var ProviderSet = wire.NewSet(Provider)
