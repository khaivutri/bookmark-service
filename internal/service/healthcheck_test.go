package service

import (
	"testing"

	"github.com/khaivutri/bookmark-service/internal/model"
	"github.com/stretchr/testify/assert"
)

type fakeHealthConfig struct {
	serviceName string
	instanceID  string
}

func (f *fakeHealthConfig) GetServiceName() string {
	return f.serviceName
}

func (f *fakeHealthConfig) GetInstanceId() string {
	return f.instanceID
}


func TestHealthCheck_GetHealthReport(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name 			string
		cfg  			*fakeHealthConfig

		expectedReport 	*model.HealthReport

	}{
		{
			name: 				"valid health report - 1",

			cfg : &fakeHealthConfig{
				serviceName: 	"bookmark_service",
				instanceID:  	"2947cc38-7c27-4b15-9d2f-50e52e638935",
			},
			
			expectedReport : &model.HealthReport{
				Message:     	"OK",
				ServiceName: 	"bookmark_service",
				InstanceID:  	"2947cc38-7c27-4b15-9d2f-50e52e638935",
			},
		},
		{
			name: 				"valid health report -2",

			cfg : &fakeHealthConfig{
				serviceName: 	"bookmark_service",
				instanceID:  	"3c983364-29f3-4501-a7cc-5603e16f6827",
			},
			
			expectedReport : &model.HealthReport{
				Message:     	"OK",
				ServiceName: 	"bookmark_service",
				InstanceID:  	"3c983364-29f3-4501-a7cc-5603e16f6827",
			},
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			hc := NewHealthCheck(tc.cfg)
			report := hc.Check()
			assert.Equal(t, tc.expectedReport, report)
		})
	}

}