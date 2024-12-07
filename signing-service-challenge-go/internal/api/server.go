package api

import (
	"encoding/json"
	"errors"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/internal/services"
	deviceService "github.com/fiskaly/coding-challenges/signing-service-challenge/internal/services/device"
	signService "github.com/fiskaly/coding-challenges/signing-service-challenge/internal/services/sign"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Response[T any] struct {
	Data T      `json:"data"`
	Err  string `json:"error_message"`
}

type PaginatedResponse[T any] struct {
	PageNumber int `json:"page_number"`
	PageSize   int `json:"page_size"`
	Total      int `json:"total"`
	Items      []T `json:"items"`
}

// Server manages HTTP requests and dispatches them to the appropriate services.
type Server struct {
	listenAddress    string
	deviceService    deviceService.DeviceService
	signatureService signService.SignService
}

// NewServer is a factory to instantiate a new Server.
func NewServer(listenAddress string, deviceService deviceService.DeviceService, signatureService signService.SignService) *Server {
	return &Server{
		listenAddress:    listenAddress,
		deviceService:    deviceService,
		signatureService: signatureService,
	}
}

// Run registers all HandlerFuncs for the existing HTTP routes and starts the Server.
func (s *Server) Run() error {
	mux := http.NewServeMux()

	mux.Handle("/api/v0/health", http.HandlerFunc(s.Health))

	// signature-devices
	mux.Handle("/api/v0/device", http.HandlerFunc(s.CreateDevice))
	mux.Handle("/api/v0/device/", http.HandlerFunc(s.GetDeviceById))
	mux.Handle("/api/v0/devices", http.HandlerFunc(s.GetAllDevices))

	// signing-creation
	mux.Handle("/api/v0/sign", http.HandlerFunc(s.CreateSigning))
	mux.Handle("/api/v0/signings", http.HandlerFunc(s.GetAllSignings))

	return http.ListenAndServe(s.listenAddress, mux)
}

// WriteErrorResponse takes an HTTP status code and a slice of errors
// and writes those as an HTTP error response in a structured format.
func WriteErrorResponse(w http.ResponseWriter, status int, err error, message string) {

	if err != nil {
		logrus.Error(err)

		// validation error
		var badRequest *services.ServiceError
		if errors.As(err, &badRequest) {
			status = badRequest.Status
			message = err.Error()
		}

		// custom db-level error
		var dbError *services.DBError
		if errors.As(err, &dbError) {
			status = http.StatusBadRequest
			message = err.Error()
		}
	}

	w.WriteHeader(status)

	bytes, err := json.Marshal(Response[string]{
		Err: message,
	})
	if err != nil {
		logrus.WithError(err).Error("error marshalling error response")
	}

	_, err = w.Write(bytes)
	if err != nil {
		logrus.WithError(err).Error("failed to write error response")
	}
}

// WriteAPIResponse takes an HTTP status code and a generic data struct
// and writes those as an HTTP response in a structured format.
func WriteAPIResponse[T any](w http.ResponseWriter, statusCode int, data T) {
	w.WriteHeader(statusCode)

	response := Response[T]{
		Data: data,
	}

	bytes, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		logrus.WithError(err).Error("failed to marshal response")
		return
	}

	_, err = w.Write(bytes)
	if err != nil {
		logrus.WithError(err).Error("failed to write api response")
	}
}
