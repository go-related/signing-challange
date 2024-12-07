package main

import (
	"github.com/fiskaly/coding-challenges/signing-service-challenge/internal/api"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/internal/configuration"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/internal/persistence"
	deviceService "github.com/fiskaly/coding-challenges/signing-service-challenge/internal/services/device"
	signService "github.com/fiskaly/coding-challenges/signing-service-challenge/internal/services/sign"
	"github.com/sirupsen/logrus"
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
	deviceSrv := deviceService.NewDeviceService(storage)
	signSrv := signService.NewSignService(storage)

	server := api.NewServer(config.ListenAddress, deviceSrv, signSrv)

	logrus.Info("starting server on port " + config.ListenAddress)
	err := server.Run()
	if err != nil {
		return err
	}

	return nil
}
