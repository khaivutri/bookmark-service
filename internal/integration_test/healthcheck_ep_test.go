package integrationtest

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/khaivutri/bookmark-service/internal/api"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheckEndpoint(t *testing.T){
	testCases := []struct{
		name 					string

		reqMethod 				string
		reqPath					string 

		setupEnv func()

		expectedConfigError 	bool 
		expectedStatusCode 		int
		expectedResponseBody 	string 
	}{
		{
			name:                 "valid health report with custom UUID",
			reqMethod:            http.MethodGet,
			reqPath:              "/health-check",
			setupEnv: func() {
				t.Setenv("INSTANCE_ID", "cbe1a562-596b-45d0-bf8b-a999b23b184a")
				t.Setenv("SERVICE_NAME", "my_service")
			},
			expectedConfigError:  false,
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"message":"OK",
									"service_name":"my_service",
									"instance_id":"cbe1a562-596b-45d0-bf8b-a999b23b184a"}`,
		},
		{
			name:                 "empty instance id should auto generate valid UUID",
			reqMethod:            http.MethodGet,
			reqPath:              "/health-check",
			setupEnv: func() {
				t.Setenv("INSTANCE_ID", "")
			},
			expectedConfigError:  false,
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: "", 
		},
		{
			name:                 "invalid UUID format should block config initialization",
			reqMethod:            http.MethodGet,
			reqPath:              "/health-check",
			setupEnv: func() {
				t.Setenv("INSTANCE_ID", "invalid-uuid-12345")
			},
			expectedConfigError:  true,
			expectedStatusCode:   0,
			expectedResponseBody: "",
		},
		{
			name:                 "invalid endpoint with POST method",
			reqMethod:            http.MethodPost, 
			reqPath:              "/wrong-endpoint", 
			setupEnv: func() {
				t.Setenv("INSTANCE_ID", "cbe1a562-596b-45d0-bf8b-a999b23b184a") 
			},
			expectedConfigError:  false,
			expectedStatusCode:   http.StatusNotFound, 
			expectedResponseBody: "",            
		},

		{
			name:                 "invalid POST method on health-check endpoint",
			reqMethod:            http.MethodPost,
			reqPath:              "/health-check", 
			setupEnv: func() {
				t.Setenv("INSTANCE_ID", "cbe1a562-596b-45d0-bf8b-a999b23b184a")
			},
			expectedConfigError:  false,
			expectedStatusCode:   http.StatusMethodNotAllowed, 
			expectedResponseBody: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setupEnv != nil {
				tc.setupEnv()
			}

			cfg, err := api.NewConfig()
			if tc.expectedConfigError {
				// Assert that an error occurred during config initialization
				assert.Error(t, err)

				// Assert that no config object was created on failure
				assert.Nil(t, cfg)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, cfg)

			testAPI := api.NewEngine(cfg)
			
			// Create a new HTTP request based on the test case's method and path
			req, _ := http.NewRequest(tc.reqMethod, tc.reqPath, nil)
			recorder := httptest.NewRecorder()
			
			testAPI.ServeHTTP(recorder, req)

			// Assert Status Code (404 for wroong endpoint, 200 for valid, 405 for wrong method)
			assert.Equal(t, tc.expectedStatusCode, recorder.Code)

			// If the expected status code is 200 OK, we can check the response body. Otherwise, we may skip this check or assert for an empty body.
			if tc.expectedStatusCode == http.StatusOK && tc.expectedResponseBody != "" {
				assert.JSONEq(t, tc.expectedResponseBody, recorder.Body.String())
			}
		})
	}
}