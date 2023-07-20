package api_test

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/iliyaisd/littlejohn/internal/api"
	"github.com/iliyaisd/littlejohn/ljlib"
	"github.com/stretchr/testify/require"
)

const testUsername = "johndoe"

func TestAPIKeyAuthorizer_Authorize(t *testing.T) {
	testCases := map[string]struct {
		authHeaderExists bool
		username         string
		password         string
		expectedError    bool
	}{
		"it should not authorize when there's no Authorization header": {
			expectedError: true,
		},
		"it should not authorize when username does not exist": {
			authHeaderExists: true,
			username:         "non-existent",
			expectedError:    true,
		},
		"it should not authorize if the password is different from default": {
			authHeaderExists: true,
			username:         "johndoe",
			password:         "password",
			expectedError:    true,
		},
		"it should authorize the user when credentials are correct, and set the user to context": {
			authHeaderExists: true,
			username:         "johndoe",
			password:         "",
		},
	}

	for testName, testCase := range testCases {
		t.Run(testName, func(t *testing.T) {
			userRepository := mockUserRepository{}
			authorizer := api.NewAPIKeyAuthorizer(userRepository)
			r, err := http.NewRequest(http.MethodGet, "", nil)
			require.NoError(t, err)
			if testCase.authHeaderExists {
				r.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString(
					[]byte(fmt.Sprintf("%s:%s", testCase.username, testCase.password))))
			}
			r, err = authorizer.Authorize(r)
			if testCase.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				user, ok := r.Context().Value("user").(*ljlib.User)
				require.True(t, ok)
				require.NotNil(t, user)
				require.Equal(t, testUsername, user.Username)
			}
		})
	}
}

type mockUserRepository struct{}

func (m mockUserRepository) GetUserByUsername(username string) (*ljlib.User, error) {
	if username == testUsername {
		return &ljlib.User{ID: uuid.MustParse("8a8d28aa-6c15-43be-8363-eb9862466063"), Username: testUsername}, nil
	}
	return nil, fmt.Errorf("user not found")
}
