package device

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/internal/domain"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/internal/services/device/mocks"
)

func TestGetAllValidity(t *testing.T) {
	tests := []struct {
		name                 string
		inputPageNumber      int
		inputPageSize        int
		expectedTotalCount   int
		expectedItemsCount   int
		expectedServiceError bool
		expectedDbError      bool
		mockData             mockData
	}{
		{
			name:               "Successful GetAll",
			inputPageNumber:    1,
			inputPageSize:      10,
			expectedTotalCount: 20,
			expectedItemsCount: 10,
			mockData:           newMockData(10, 20, nil),
		},
		{
			name:            "Repository Error",
			inputPageNumber: 1,
			inputPageSize:   10,
			expectedDbError: true,
			mockData:        newMockData(100, 100, fmt.Errorf("repository error")),
		},
		{
			name:                 "Invalid Input Page",
			inputPageNumber:      0,
			inputPageSize:        10,
			expectedServiceError: true,
			mockData:             newMockData(0, 0, nil),
		},
		{
			name:                 "Invalid Page Size",
			inputPageNumber:      1,
			inputPageSize:        0,
			expectedServiceError: true,
			mockData:             newMockData(0, 0, nil),
		},
		{
			name:               "Out of Range page",
			inputPageNumber:    11,
			inputPageSize:      10,
			expectedTotalCount: 100,
			expectedItemsCount: 0,
			mockData:           newMockData(0, 100, nil),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// setup
			mockRepo := new(mocks.MockDeviceRepository)
			if !test.expectedServiceError {
				mockRepo.On("GetAll", test.inputPageNumber, test.inputPageSize).
					Return(test.mockData.Devices, test.mockData.TotalCount, test.mockData.Error)
			}

			service := NewDeviceService(mockRepo)

			// execute
			devices, totalCount, err := service.GetAll(test.inputPageNumber, test.inputPageSize)

			if test.expectedServiceError || test.expectedDbError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expectedTotalCount, totalCount)
				assert.Equal(t, test.expectedItemsCount, len(devices))
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestGetById(t *testing.T) {
	tests := []struct {
		name          string
		inputDeviceId string
		mockDevice    *domain.Device
		mockError     error
		expectedError bool
	}{
		{
			name:          "Successful GetById",
			inputDeviceId: "1",
			mockDevice:    &domain.Device{ID: "1"},
		},
		{
			name:          "Db Error",
			inputDeviceId: "1",
			mockError:     fmt.Errorf("db error"),
			expectedError: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// setup
			mockRepo := new(mocks.MockDeviceRepository)

			mockRepo.On("FindByID", test.inputDeviceId).
				Return(test.mockDevice, test.mockError)

			service := NewDeviceService(mockRepo)

			// execute
			device, err := service.GetById(test.inputDeviceId)

			// assertions
			if test.expectedError {
				assert.Error(t, err)
				assert.Nil(t, device)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.inputDeviceId, device.ID)
			}

			mockRepo.AssertExpectations(t)

		})
	}
}

func TestSave(t *testing.T) {
	tests := []struct {
		name                 string
		inputDevice          *domain.Device
		mockError            error
		expectedServiceError bool
		expectedDbError      bool
	}{
		{
			name: "Valid ECC Type",
			inputDevice: &domain.Device{
				ID:            "1",
				AlgorithmType: domain.AlgorithmTypeECC,
			},
		},
		{
			name: "Valid RSA Type",
			inputDevice: &domain.Device{
				ID:            "1",
				AlgorithmType: domain.AlgorithmTypeRSA,
			},
		},
		{
			name: "Invalid Algorithm Type",
			inputDevice: &domain.Device{
				ID:            "1",
				AlgorithmType: domain.AlgorithmTypeUnknown,
			},
			expectedServiceError: true,
		},
		{
			name: "Invalid Id",
			inputDevice: &domain.Device{
				ID:            "",
				AlgorithmType: domain.AlgorithmTypeRSA,
			},
			expectedServiceError: true,
		},
		{
			name:                 "Invalid Input",
			inputDevice:          nil,
			expectedServiceError: true,
		},
		{
			name: "Db Error",
			inputDevice: &domain.Device{
				ID:            "1",
				AlgorithmType: domain.AlgorithmTypeRSA,
			},
			expectedDbError: true,
			mockError:       fmt.Errorf("db error"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mockRepo := new(mocks.MockDeviceRepository)
			service := NewDeviceService(mockRepo)
			if !test.expectedServiceError {
				mockRepo.On("Save", mock.Anything).Return(test.mockError)
			}

			err := service.Save(test.inputDevice)

			if test.expectedDbError || test.expectedServiceError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			if !test.expectedServiceError {
				mockRepo.AssertExpectations(t)
			}

		})
	}

}

type mockData struct {
	Devices    []*domain.Device
	TotalCount int
	Error      error
}

func newMockData(mockCount int, totalCount int, err error) mockData {
	var devices []*domain.Device
	for i := 0; i < mockCount; i++ {
		devices = append(devices, &domain.Device{
			ID: strconv.Itoa(i),
		})
	}
	return mockData{
		Devices:    devices,
		TotalCount: totalCount,
		Error:      err,
	}
}
