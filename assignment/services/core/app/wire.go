//go:build wireinject
// +build wireinject

package main

import (
	"project/internal/auth"
	"project/pkg/casbin"
	"project/pkg/db/cache"
	"project/pkg/db/sqlx"
	"project/services/core/repo"
	"project/services/core/service"
	"project/services/core/web"

	"github.com/google/wire"
)

var CoreModuleSet = wire.NewSet(
	wire.FieldsOf(new(*CoreConfig), "server", "db"),
	NewCacheConfig,
	cache.NewRedisInstance,
	casbin.InitCasbin,
	sqlx.NewDBConn,

	wire.Struct(new(web.CoreHandlerRegistryOptions), "*"),
	web.NewCoreHandlerRegistry,

	auth.NewAuthImpl,

	repo.NewAssignmentUserRepoImpl,
	wire.Bind(new(service.AssignmentUserRepo), new(*repo.AssignmentUserRepoImpl)),

	service.NewAssignmentUserServiceImpl,
	wire.Bind(new(service.AssignmentUserService), new(*service.AssignmentUserServiceImpl)),

// -----==-----==DO NOT ADD CODE BELOW THIS LINE------
)

func initServer(config *CoreConfig) (*web.CoreHandlerRegistry, error) {
	panic(wire.Build(CoreModuleSet))
}
