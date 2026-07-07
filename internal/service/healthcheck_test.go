package service

import (
	"testing"

	"github.com/khaivutri/bookmark-service/internal/model"
	"github.com/stretchr/testify/assert"
)

func TestHealthCheck_GetHealthReport(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name         		string
		serviceName  		string
		instanceID   		string
		expectedReport 		*model.HealthReport
	}{
		{
			name:        		"valid health report - 1",
			serviceName: 		"bookmark_service",
			instanceID:  		"2947cc38-7c27-4b15-9d2f-50e52e638935",
			expectedReport: 	&model.HealthReport{
									Message:     "OK",
									ServiceName: "bookmark_service",
									InstanceID:  "2947cc38-7c27-4b15-9d2f-50e52e638935",
								},
		},
		{
			name:        		"valid health report -2",
			serviceName: 		"bookmark_service",
			instanceID:  		"3c983364-29f3-4501-a7cc-5603e16f6827",
			expectedReport: 	&model.HealthReport{
									Message:     "OK",
									ServiceName: "bookmark_service",
									InstanceID:  "3c983364-29f3-4501-a7cc-5603e16f6827",
								},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			hc := NewHealthCheck(tc.serviceName, tc.instanceID)
			report := hc.Check()
			assert.Equal(t, tc.expectedReport, report)
		})
	}
}