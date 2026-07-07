package service

import "github.com/khaivutri/bookmark-service/internal/model"

// HealthCheck defines the interface for health check service.
type HealthCheck interface {
	Check() *model.HealthReport
}

type healthCheck struct {
	serviceName 	string
	instanceID  	string
}

// NewHealthCheck creates and returns a new HealthCheck service instance.
func NewHealthCheck(serviceName, instanceID string) HealthCheck {
	return &healthCheck{	
							serviceName: serviceName,
							instanceID: instanceID,
						}
}

// Check returns the health status report for the service.
func (hc *healthCheck) Check() *model.HealthReport {
	return &model.HealthReport{
		Message:     	"OK",
		ServiceName:	hc.serviceName,
		InstanceID:  	hc.instanceID,
	}
}