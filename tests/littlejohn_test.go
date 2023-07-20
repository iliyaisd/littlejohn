package tests

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	tickersPath    = "http://localhost:8080/tickers"
	historyPathTpl = "http://localhost:8080/tickers/%s/history"
)

func TestPortfolio(t *testing.T) {
	testCases := map[string]struct {
		login        string
		expectedCode int
	}{
		"it should return http status 403 if no basic auth was provided": {
			expectedCode: http.StatusForbidden,
		},
		"it should return http status 403 on wrong user login": {
			login:        "non-existent",
			expectedCode: http.StatusForbidden,
		},
		"it should return idempotent results for the same valid user": {
			login:        "johndoe",
			expectedCode: http.StatusOK,
		},
	}

	for testName, testCase := range testCases {
		t.Run(testName, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, tickersPath, nil)
			require.NoError(t, err)

			if len(testCase.login) > 0 {
				req.Header.Add("Authorization", "Basic "+
					base64.StdEncoding.EncodeToString([]byte(testCase.login+":")))
			}

			client := http.Client{}
			resp, err := client.Do(req)
			require.NoError(t, err)

			assert.Equal(t, testCase.expectedCode, resp.StatusCode)
			if resp.StatusCode >= 400 {
				return
			}

			//on success, doing second request to check idempotency
			otherResp, err := client.Do(req)
			require.NoError(t, err)
			assert.Equal(t, resp.StatusCode, otherResp.StatusCode)

			var tickers, otherTickers []interface{}
			require.NoError(t, json.NewDecoder(resp.Body).Decode(&tickers))
			require.NoError(t, json.NewDecoder(otherResp.Body).Decode(&otherTickers))

			assert.ElementsMatch(t, tickers, otherTickers)
		})
	}
}

func TestHistoricalPrices(t *testing.T) {
	testCases := map[string]struct {
		login               string
		ticker              string
		expectedCode        int
		expectedResultCount int
	}{
		"it should return http status 403 if no basic auth was provided": {
			expectedCode: http.StatusForbidden,
			ticker:       "AAPL",
		},
		"it should return http status 403 on wrong user login": {
			login:        "non-existent",
			ticker:       "AAPL",
			expectedCode: http.StatusForbidden,
		},
		"it should return http status 404 if ticker does not exist for user": {
			login:        "johndoe",
			ticker:       "wrong_name",
			expectedCode: http.StatusNotFound,
		},
		"it should return correct number of days of results on success": {
			login:               "johndoe",
			ticker:              "GOOG",
			expectedCode:        http.StatusOK,
			expectedResultCount: 90,
		},
	}

	for testName, testCase := range testCases {
		t.Run(testName, func(t *testing.T) {
			req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(historyPathTpl, testCase.ticker), nil)
			require.NoError(t, err)

			if len(testCase.login) > 0 {
				req.Header.Add("Authorization", "Basic "+
					base64.StdEncoding.EncodeToString([]byte(testCase.login+":")))
			}

			client := http.Client{}
			resp, err := client.Do(req)
			require.NoError(t, err)

			assert.Equal(t, testCase.expectedCode, resp.StatusCode)
			if resp.StatusCode >= 400 {
				return
			}

			var prices []interface{}
			require.NoError(t, json.NewDecoder(resp.Body).Decode(&prices))

			assert.Equal(t, testCase.expectedResultCount, len(prices))
		})
	}
}
