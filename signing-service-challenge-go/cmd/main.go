package main

import (
	"github.com/sirupsen/logrus"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/internal/api"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/internal/configuration"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/internal/crypto"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/internal/persistence"
	deviceService "github.com/fiskaly/coding-challenges/signing-service-challenge/internal/services/device"
	signService "github.com/fiskaly/coding-challenges/signing-service-challenge/internal/services/sign"
)

func main() {
	config, err := configuration.LoadConfiguration()
	if err != nil {
		logrus.Fatal(err)
	}

	err = runServer(config)
	if err != nil {
		logrus.Fatal(err)
	}
}

func runServer(config *configuration.Configuration) error {
	// repositories
	storage := persistence.NewInMemoryStorage()

	//services
	factory := crypto.NewFactory()
	deviceSrv := deviceService.NewDeviceService(storage, factory)
	signSrv := signService.NewSignService(storage, factory)

	server := api.NewServer(config.ListenAddress, deviceSrv, signSrv)

	logrus.Info("starting server on port " + config.ListenAddress)
	err := server.Run()
	if err != nil {
		return err
	}

	return nil
}
