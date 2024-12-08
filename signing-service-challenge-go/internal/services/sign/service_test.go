package sign

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/internal/services"
	"github.com/stretchr/testify/mock"
	"testing"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/internal/crypto"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/internal/domain"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/internal/services/sign/mocks"
	"github.com/stretchr/testify/assert"
)

func TestGetAllSignings(t *testing.T) {
	tests := []struct {
		name           string
		mockDevice     *domain.Device
		inputDeviceId  string
		inputPageNr    int
		inputPageSize  int
		mockSignings   []*domain.Signings
		mockTotalCount int
		mockDbError    error
		expectError    bool
	}{
		{
			name:           "Successful GetAllSignings",
			inputDeviceId:  "device-1",
			inputPageNr:    1,
			inputPageSize:  10,
			mockSignings:   []*domain.Signings{{DeviceId: "device-1"}},
			mockTotalCount: 1,
			expectError:    false,
		},
		{
			name:          "Empty DeviceID",
			inputDeviceId: "",
			inputPageNr:   1,
			inputPageSize: 10,
			expectError:   true,
		},
		{
			name:          "Invalid Page Number",
			inputDeviceId: "device-1",
			inputPageNr:   0,
			inputPageSize: 10,
			expectError:   true,
		},
		{
			name:          "Invalid Page Size",
			inputDeviceId: "device-1",
			inputPageNr:   1,
			inputPageSize: 0,
			expectError:   true,
		},
		{
			name:          "Repository Error",
			inputDeviceId: "device-1",
			inputPageNr:   1,
			inputPageSize: 10,
			mockDbError:   errors.New("repository error"),
			expectError:   true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockRepo := new(mocks.MockSignRepository)
			service := NewSignService(mockRepo, nil)

			if !test.expectError || test.mockDbError != nil {
				mockRepo.On("GetAllSignings", test.inputDeviceId, test.inputPageNr, test.inputPageSize).
					Return(test.mockSignings, test.mockTotalCount, test.mockDbError).
					Once()
			}

			// execute
			signings, totalCount, err := service.GetAllSignings(test.inputDeviceId, test.inputPageNr, test.inputPageSize)

			if test.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.mockTotalCount, totalCount)
				assert.Equal(t, len(test.mockSignings), len(signings))
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestSignTransactionValidity(t *testing.T) {
	tests := []struct {
		name             string
		inputCounter     int64
		tp               domain.AlgorithmType
		inputLastEncoded string
		inputDeviceId    string
		inputData        string
		getDeviceError   error
		saveDeviceError  error
		expectedError    bool
	}{
		{
			name:             "Valid ECC with counter 0",
			inputCounter:     0,
			tp:               domain.AlgorithmTypeECC,
			inputLastEncoded: "",
			inputDeviceId:    "testing1",
			inputData:        "testing---1",
		},
		{
			name:             "Valid ECC with counter 25",
			inputCounter:     25,
			tp:               domain.AlgorithmTypeECC,
			inputLastEncoded: "test",
			inputDeviceId:    "testing1",
			inputData:        "testing---1",
		},
		{
			name:             "Valid RSA with counter 0",
			inputCounter:     0,
			tp:               domain.AlgorithmTypeRSA,
			inputLastEncoded: "",
			inputDeviceId:    "testing1",
			inputData:        "testing---1",
		},
		{
			name:             "Valid RSA with counter 2",
			inputCounter:     2,
			tp:               domain.AlgorithmTypeRSA,
			inputLastEncoded: "test",
			inputDeviceId:    "testing1",
			inputData:        "testing---1",
		},
		{
			name:             "Valid Input with getError",
			inputCounter:     0,
			tp:               domain.AlgorithmTypeECC,
			inputLastEncoded: "",
			inputDeviceId:    "testing1",
			inputData:        "testing---1",
			getDeviceError:   services.NewDBError("get error"),
			expectedError:    true,
		},
		{
			name:             "Valid Input with saveError",
			inputCounter:     0,
			tp:               domain.AlgorithmTypeECC,
			inputLastEncoded: "",
			inputDeviceId:    "testing1",
			inputData:        "testing---1",
			saveDeviceError:  services.NewDBError("get error"),
			expectedError:    true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// setup service and mocks
			mockRepo := new(mocks.MockSignRepository)
			factory := crypto.NewFactory()
			mockDevice, _, expectedData := generateDeviceModel(t, test.inputDeviceId, test.inputCounter, test.tp, test.inputData, test.inputLastEncoded)
			service := NewSignService(mockRepo, factory)

			mockRepo.On("GetDeviceCounterAndLastEncoded", test.inputDeviceId).Return(test.inputCounter, test.inputLastEncoded, test.getDeviceError).Once()
			if test.getDeviceError == nil {
				mockRepo.On("SaveDeviceCounterAndLastEncoded", test.inputDeviceId, test.inputCounter+1, mock.Anything, mock.Anything).Return(test.saveDeviceError).Once()
			}

			// execute
			_, signedData, err := service.signTransaction(mockDevice, []byte(test.inputData))

			// asserts
			if test.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)

				// this can be equal since when signing we are using random reader
				//assert.Equal(t, expectedSignature, signature)
				assert.Equal(t, expectedData, signedData)
			}
			mockRepo.AssertExpectations(t)
		})
	}

}

func generateDeviceModel(t *testing.T, id string, counter int64, tp domain.AlgorithmType, data, lastSignature string) (*domain.Device, string, string) {
	factory := crypto.NewFactory()
	algorithm, err := factory.GenerateAlgorithm(tp)
	if err != nil {
		t.Error(err)
	}
	marshaller, err := factory.CreateMarshaller(tp)
	if err != nil {
		t.Error(err)
	}
	pbytes, privateBytes, err := marshaller.Encode(algorithm)
	if err != nil {
		t.Error(err)
	}

	// signature
	signature, err := algorithm.Sign([]byte(data))
	if err != nil {
		t.Error(err)
	}
	currentSignatureEncoded := base64.StdEncoding.EncodeToString(signature)

	if counter == 0 {
		lastSignature = base64.StdEncoding.EncodeToString([]byte(id))
	}
	signedData := fmt.Sprintf("%d_%s_%s", counter+1, data, lastSignature)

	return &domain.Device{
		ID:            id,
		AlgorithmType: tp,
		Counter:       counter,
		PublicKey:     pbytes,
		PrivateKey:    privateBytes,
	}, currentSignatureEncoded, signedData
}
