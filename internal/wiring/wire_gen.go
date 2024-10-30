// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package wiring

import (
	"GoLoad/internal/configs"
	"GoLoad/internal/dataaccess"
	"GoLoad/internal/dataaccess/database"
	"GoLoad/internal/handler"
	"GoLoad/internal/handler/grpc"
	"GoLoad/internal/logic"
	"github.com/google/wire"
)

// Injectors from wire.go:

func InitializeGRPCServer(configFilePath configs.ConfigFilePath) (grpc.Server, func(), error) {
	config, err := configs.NewConfig(configFilePath)
	if err != nil {
		return nil, nil, err
	}
	configsDatabase := config.Database
	db, cleanup, err := database.InitializeDB(configsDatabase)
	if err != nil {
		return nil, nil, err
	}
	goquDatabase := database.InitializeGoquDB(db)
	accountDataAccessor := database.NewAccountDataAccessor(goquDatabase)
	accountPasswordDataAccessor := database.NewAccountPasswordDataAccessor(goquDatabase)
	auth := config.Auth
	hash := logic.NewHash(auth)
	account := logic.NewAccount(goquDatabase, accountDataAccessor, accountPasswordDataAccessor, hash)
	goLoadServiceServer := grpc.NewHandler(account)
	server := grpc.NewServer(goLoadServiceServer)
	return server, func() {
		cleanup()
	}, nil
}

// wire.go:

var WireSet = wire.NewSet(configs.WireSet, dataaccess.WireSet, logic.WireSet, handler.WireSet)
