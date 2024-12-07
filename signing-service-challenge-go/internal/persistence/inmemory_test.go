package persistence

import (
	"github.com/fiskaly/coding-challenges/signing-service-challenge/internal/domain"
	"strconv"
	"testing"
)

func TestGetAll(t *testing.T) {

	tests := []struct {
		name                 string
		totalElements        int
		page                 int
		pageSize             int
		expectedResultLength int
		expectedTotal        int
	}{
		{
			name:                 "test with 100 element page 1",
			totalElements:        100,
			page:                 1,
			pageSize:             10,
			expectedResultLength: 10,
			expectedTotal:        100,
		},
		{
			name:                 "test with 100 element page 2",
			totalElements:        100,
			page:                 2,
			pageSize:             10,
			expectedResultLength: 10,
			expectedTotal:        100,
		},
		{
			name:                 "test with 100 element page 11",
			totalElements:        100,
			page:                 11,
			pageSize:             10,
			expectedResultLength: 0,
			expectedTotal:        100,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// prepare
			store := NewInMemoryStorage()
			for i := 0; i < test.totalElements; i++ {
				if err := store.Save(domain.SignatureDevice{ID: strconv.Itoa(i), Counter: int64(i)}); err != nil {
					t.Error(err)
				}
			}

			list, total, err := store.GetAll(test.page, test.pageSize)
			if err != nil {
				t.Error(err)
			}

			if len(list) != test.expectedResultLength {
				t.Errorf("Expected length %d, got %d", test.expectedResultLength, len(list))
			}
			if total != test.expectedTotal {
				t.Errorf("Expected total %d, got %d", test.expectedTotal, total)
			}

		})
	}
}
