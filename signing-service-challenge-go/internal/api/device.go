package api

import (
	"encoding/json"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/internal/domain"
	"net/http"
	"strconv"
)

type SigningDeviceDTO struct {
	Id        string  `json:"id"`
	Algorithm string  `json:"algorithm"` // the validation is done on the service level, so we delegate the check there
	Label     *string `json:"label,omitempty"`
	Counter   int     `json:"counter"`
}

func (s *Server) CreateSigningDevice(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		WriteErrorResponse(response, http.StatusMethodNotAllowed, nil, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	var device SigningDeviceDTO
	if err := json.NewDecoder(request.Body).Decode(&device); err != nil {
		WriteErrorResponse(response, http.StatusBadRequest, err, "Invalid request payload")
		return
	}

	input := convertDeviceDTOtoDomainModel(&device)
	if err := s.deviceService.Save(input); err != nil {
		WriteErrorResponse(response, http.StatusInternalServerError, err, http.StatusText(http.StatusInternalServerError))
		return
	}

	result, err := s.deviceService.GetById(device.Id)
	if err != nil {
		WriteErrorResponse(response, http.StatusInternalServerError, err, http.StatusText(http.StatusInternalServerError))
		return
	}
	WriteAPIResponse(response, http.StatusCreated, result)
}

func (s *Server) GetSigningDeviceById(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		WriteErrorResponse(response, http.StatusMethodNotAllowed, nil, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	id := request.URL.Query().Get("id")
	device, err := s.deviceService.GetById(id)
	if err != nil {
		WriteErrorResponse(response, http.StatusInternalServerError, err, http.StatusText(http.StatusInternalServerError))
		return
	}

	result := convertDeviceDomainModelToDTO(device)
	WriteAPIResponse(response, http.StatusOK, result)
}

func (s *Server) GetAllDevices(response http.ResponseWriter, request *http.Request) {
	pageNr, err := strconv.Atoi(request.URL.Query().Get("pageNr"))
	if err != nil || pageNr < 1 {
		WriteErrorResponse(response, http.StatusBadRequest, err, "Invalid or missing pageNr")
		return
	}

	pageSize, err := strconv.Atoi(request.URL.Query().Get("pageSize"))
	if err != nil || pageSize < 1 {
		WriteErrorResponse(response, http.StatusBadRequest, err, "Invalid or missing pageSize")
		return
	}

	devices, totalCount, err := s.deviceService.GetAll(pageNr, pageSize)
	if err != nil {
		WriteErrorResponse(response, http.StatusInternalServerError, err, "Failed to retrieve devices")
		return
	}
	result := convertDeviceListDomainModelToDTO(&devices, pageNr, pageSize, totalCount)

	WriteAPIResponse(response, http.StatusOK, result)
}

func convertDeviceDTOtoDomainModel(input *SigningDeviceDTO) *domain.SignatureDevice {
	if input == nil {
		return nil
	}
	return &domain.SignatureDevice{
		ID:            input.Id,
		Label:         input.Label,
		Counter:       int64(input.Counter),
		AlgorithmType: domain.ConvertStringToAlgorithmType(input.Algorithm),
	}
}

func convertDeviceDomainModelToDTO(input *domain.SignatureDevice) *SigningDeviceDTO {
	if input == nil {
		return nil
	}
	return &SigningDeviceDTO{
		Id:        input.ID,
		Label:     input.Label,
		Counter:   int(input.Counter),
		Algorithm: string(input.AlgorithmType),
	}
}

func convertDeviceListDomainModelToDTO(input *[]*domain.SignatureDevice, page, pageSize, total int) *PaginatedResponse[SigningDeviceDTO] {
	if input == nil {
		return nil
	}
	var results []SigningDeviceDTO
	for _, device := range *input {
		results = append(results, *convertDeviceDomainModelToDTO(device))
	}
	return &PaginatedResponse[SigningDeviceDTO]{
		Items:      results,
		Total:      total,
		PageNumber: page,
		PageSize:   pageSize,
	}
}
