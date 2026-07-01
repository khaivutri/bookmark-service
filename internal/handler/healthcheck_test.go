package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/khaivutri/bookmark-service/internal/model"
	"github.com/khaivutri/bookmark-service/internal/service/mocks"
	"github.com/stretchr/testify/assert"
)

func TestHealtCheck_HealthCheck(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupRequest func(ctx *gin.Context)
		setupMockService func(context context.Context) *mocks.HealthCheck

		expectedStatusCode int
		expectedResponseBody string 
	}{
		{
			name: 			"valid health report - 1",

			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/health-check", nil)
			},

			setupMockService: func(ctx context.Context) *mocks.HealthCheck {
				mockSvc := mocks.NewHealthCheck(t)
				mockSvc.On("Check").Return(&model.HealthReport{	Message: 		"OK",
																			ServiceName: 	"bookmark_service", 
																			InstanceID: 	"cbe1a562-596b-45d0-bf8b-a999b23b184a"}).Once(	)
				return mockSvc
			},

			expectedStatusCode: http.StatusOK,
			expectedResponseBody: `{"message":"OK",
									"service_name":"bookmark_service",
									"instance_id":"cbe1a562-596b-45d0-bf8b-a999b23b184a"}`,
		},
		{
			name: 			"valid health report  - 2",

			setupRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/health-check", nil)
			},

			setupMockService: func(ctx context.Context) *mocks.HealthCheck {
				mockSvc := mocks.NewHealthCheck(t)
				mockSvc.On("Check").Return(&model.HealthReport{	Message: 		"OK",
																			ServiceName: 	"my_service", 
																			InstanceID: 	"cbe1a562-596b-45d0-bf8b-a999b23b184a"}).Once(	)
				return mockSvc
			},

			expectedStatusCode: http.StatusOK,
			expectedResponseBody: `{"message":"OK",
									"service_name":"my_service",
									"instance_id":"cbe1a562-596b-45d0-bf8b-a999b23b184a"}`,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()	
			ctx, _ := gin.CreateTestContext(rec)

			testCase.setupRequest(ctx)

			mockSvc := testCase.setupMockService(ctx)

			testHandler := NewHealthCheck(mockSvc)
			testHandler.HealthCheck(ctx)

			assert.Equal(t, testCase.expectedStatusCode, rec.Code)
			assert.JSONEq(t, testCase.expectedResponseBody, rec.Body.String())
		})
	}
}