package sign

import (
	"errors"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/internal/domain"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/internal/services/sign/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
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
