package api

import (
	"encoding/json"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/internal/domain"
	"net/http"
	"strconv"
)

type SigningInputDTO struct {
	DeviceID string `json:"device_id"`
	Data     string `json:"data"`
}

type SigningResultDTO struct {
	Signature  string `json:"signature"`
	SignedData string `json:"signed_data"`
}

func (s *Server) CreateSigning(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		WriteErrorResponse(response, http.StatusMethodNotAllowed, nil, http.StatusText(http.StatusMethodNotAllowed))
		return
	}

	var input SigningInputDTO
	if err := json.NewDecoder(request.Body).Decode(&input); err != nil {
		WriteErrorResponse(response, http.StatusBadRequest, err, "Invalid request payload")
		return
	}

	signature, signedData, err := s.signatureService.Sign(input.DeviceID, []byte(input.Data))
	if err != nil {
		WriteErrorResponse(response, http.StatusInternalServerError, err, http.StatusText(http.StatusInternalServerError))
		return
	}

	output := SigningResultDTO{
		Signature:  signature,
		SignedData: signedData,
	}
	WriteAPIResponse(response, http.StatusCreated, output)
}

func (s *Server) GetAllSignings(response http.ResponseWriter, request *http.Request) {
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

	list, totalCount, err := s.signatureService.GetAllSignings(deviceId, pageNr, pageSize)
	if err != nil {
		WriteErrorResponse(response, http.StatusInternalServerError, err, "Failed to retrieve devices")
		return
	}

	output := convertSigningListDomainModelToDTO(&list, pageNr, pageSize, totalCount)
	WriteAPIResponse(response, http.StatusOK, output)
}

func convertSigningListDomainModelToDTO(i *[]*domain.Signings, page int, size int, total int) *PaginatedResponse[SigningResultDTO] {
	if i == nil {
		return &PaginatedResponse[SigningResultDTO]{
			Items: nil,
			Total: total,
		}
	}

	var items []SigningResultDTO
	for _, item := range *i {
		if item != nil {
			items = append(items, SigningResultDTO{
				Signature:  item.Signature,
				SignedData: item.SignedData,
			})
		}
	}
	return &PaginatedResponse[SigningResultDTO]{
		Items:      items,
		Total:      total,
		PageNumber: page,
		PageSize:   size,
	}

}
