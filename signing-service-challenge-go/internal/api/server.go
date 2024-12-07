package api

import (
	"encoding/json"
	"errors"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/internal/services"
	signaturecreation "github.com/fiskaly/coding-challenges/signing-service-challenge/internal/services/signature-creation"
	signaturedevice "github.com/fiskaly/coding-challenges/signing-service-challenge/internal/services/signature-device"
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
	deviceService    signaturedevice.SignatureDeviceService
	signatureService signaturecreation.SignatureService
}

// NewServer is a factory to instantiate a new Server.
func NewServer(listenAddress string, deviceService signaturedevice.SignatureDeviceService, signatureService signaturecreation.SignatureService) *Server {
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

	//signature-devices
	mux.Handle("/api/v0/signature-device", http.HandlerFunc(s.CreateSigningDevice))
	mux.Handle("/api/v0/signature-device/", http.HandlerFunc(s.GetSigningDeviceById))
	mux.Handle("/api/v0/signature-devices", http.HandlerFunc(s.GetAllDevices))

	return http.ListenAndServe(s.listenAddress, mux)
}

// WriteErrorResponse takes an HTTP status code and a slice of errors
// and writes those as an HTTP error response in a structured format.
func WriteErrorResponse(w http.ResponseWriter, status int, err error, message string) {

	// we check if the error type is like this than meas it's a bad request
	var badRequest *services.ServiceError
	if errors.As(err, &badRequest) {
		status = badRequest.Status
		message = err.Error()
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
		WriteErrorResponse(w, http.StatusInternalServerError, err, "failed to marshal response")
	}

	_, err = w.Write(bytes)
	if err != nil {
		logrus.WithError(err).Error("failed to write api response")
	}
}
