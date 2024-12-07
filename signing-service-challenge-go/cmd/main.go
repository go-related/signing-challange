package main

import (
	"github.com/fiskaly/coding-challenges/signing-service-challenge/internal/api"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/internal/configuration"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/internal/persistence"
	signaturecreation "github.com/fiskaly/coding-challenges/signing-service-challenge/internal/services/signature-creation"
	signaturedevice "github.com/fiskaly/coding-challenges/signing-service-challenge/internal/services/signature-device"
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
	signatureService := signaturedevice.NewDeviceService(storage)
	signatureCreationService := signaturecreation.NewSignatureCreation(storage)

	server := api.NewServer(config.ListenAddress, signatureService, signatureCreationService)

	logrus.Info("starting server on port " + config.ListenAddress)
	err := server.Run()
	if err != nil {
		return err
	}

	return nil
}
