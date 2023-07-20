package datasource_test

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/iliyaisd/littlejohn/internal/datasource"
	"github.com/iliyaisd/littlejohn/ljlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLocalDatasource_GetHistoricalPrices(t *testing.T) {
	testCases := map[string]struct {
		ticker         string
		expectedError  bool
		expectedLength int
	}{
		"it should return IllegalArgumentError for non-existent ticker": {
			ticker:        "non-existent",
			expectedError: true,
		},
		"it should return correct number of prices for an existent ticker": {
			ticker:         "AAPL",
			expectedLength: 11,
		},
	}
	for testName, testCase := range testCases {
		t.Run(testName, func(t *testing.T) {
			localDS := datasource.NewLocalDatasource()
			prices, err := localDS.GetHistoricalPrices(testCase.ticker, mustParseDate(t, "2023-02-10"), mustParseDate(t, "2023-02-20"))
			if testCase.expectedError {
				require.Error(t, err)
				assert.True(t, errors.As(err, &ljlib.IllegalArgumentError{}))
			} else {
				require.NoError(t, err)
				assert.Equal(t, testCase.expectedLength, len(prices))
			}
		})
	}
}

func TestLocalDatasource_GetUserPortfolio(t *testing.T) {
	testCases := map[string]struct {
		userUUID       uuid.UUID
		expectedError  bool
		expectedLength int
	}{
		"it should return not found error for non-existent users": {
			userUUID:      uuid.MustParse("1a97a523-e1dd-4a9b-878e-cc4d391e9860"),
			expectedError: true,
		},
		"it should return correct number of items in portfolio for an existent ticker": {
			userUUID:       uuid.MustParse("f2f208c8-16a4-4ef6-80e3-88103f6471a2"),
			expectedLength: 8,
		},
	}
	for testName, testCase := range testCases {
		t.Run(testName, func(t *testing.T) {
			localDS := datasource.NewLocalDatasource()
			prices, err := localDS.GetUserPortfolio(testCase.userUUID)
			if testCase.expectedError {
				require.Error(t, err)
				assert.True(t, errors.As(err, &ljlib.NotFoundError{}))
			} else {
				require.NoError(t, err)
				assert.Equal(t, testCase.expectedLength, len(prices))
			}
		})
	}
}

func mustParseDate(t *testing.T, dt string) time.Time {
	tm, err := time.Parse(time.DateOnly, dt)
	require.NoError(t, err)
	return tm
}
