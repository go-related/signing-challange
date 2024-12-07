package api

import (
	"encoding/json"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/internal/domain"
	"net/http"
	"strconv"
)

type SigningCreationInputDTO struct {
	DeviceID string `json:"device_id"`
	Data     string `json:"data"`
}

type SigningCreationResultDTO struct {
	Signature  string `json:"signature"`
	SignedData string `json:"signed_data"`
}

func (s *Server) CreateSigning(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		WriteErrorResponse(response, http.StatusMethodNotAllowed, nil, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	var input SigningCreationInputDTO
	if err := json.NewDecoder(request.Body).Decode(&input); err != nil {
		WriteErrorResponse(response, http.StatusBadRequest, err, "Invalid request payload")
		return
	}

	signature, signedData, err := s.signatureService.SignTransaction(input.DeviceID, []byte(input.Data))
	if err != nil {
		WriteErrorResponse(response, http.StatusInternalServerError, err, http.StatusText(http.StatusInternalServerError))
		return
	}

	output := SigningCreationResultDTO{
		Signature:  signature,
		SignedData: signedData,
	}
	WriteAPIResponse(response, http.StatusCreated, output)
}

func (s *Server) GetAllSigningCreations(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		WriteErrorResponse(response, http.StatusMethodNotAllowed, nil, http.StatusText(http.StatusMethodNotAllowed))
		return
	}
	deviceId := request.URL.Query().Get("deviceId")
	if deviceId == "" {
		WriteErrorResponse(response, http.StatusBadRequest, nil, "Invalid or missing deviceId")
		return
	}

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

	list, totalCount, err := s.signatureService.GetAllSigningCreations(deviceId, pageNr, pageSize)
	if err != nil {
		WriteErrorResponse(response, http.StatusInternalServerError, err, "Failed to retrieve devices")
		return
	}

	output := convertSignatureListDomainModelToDTO(&list, pageNr, pageSize, totalCount)
	WriteAPIResponse(response, http.StatusOK, output)
}

func convertSignatureListDomainModelToDTO(i *[]*domain.SignedCreations, page int, size int, total int) *PaginatedResponse[SigningCreationResultDTO] {
	if i == nil {
		return &PaginatedResponse[SigningCreationResultDTO]{
			Items: nil,
			Total: total,
		}
	}

	var items []SigningCreationResultDTO
	for _, item := range *i {
		if item != nil {
			items = append(items, SigningCreationResultDTO{
				Signature:  item.Signature,
				SignedData: item.SignedData,
			})
		}
	}
	return &PaginatedResponse[SigningCreationResultDTO]{
		Items:      items,
		Total:      total,
		PageNumber: page,
		PageSize:   size,
	}

}
