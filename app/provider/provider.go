package provider

import (
	"github.com/aristat/golang-gin-oauth2-example-app/app/config"
	"github.com/aristat/golang-gin-oauth2-example-app/app/entrypoint"
	"github.com/aristat/golang-gin-oauth2-example-app/app/logger"
	"github.com/google/wire"
)

var AwareProductionSet = wire.NewSet(
	entrypoint.ProviderProductionSet,
	logger.ProviderProductionSet,
	config.ProviderSet,
)

var AwareTestSet = wire.NewSet(
	entrypoint.ProviderTestSet,
	logger.ProviderTestSet,
)
