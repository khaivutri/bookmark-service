package service

import "github.com/khaivutri/bookmark-service/internal/model"

// HealthCheck defines the interface for health check service.
type HealthCheck interface {
	Check() *model.HealthReport
}

type healthCheck struct {
	config model.HealthConfig
}


// NewHealthCheck creates and returns a new HealthCheck service instance.
func NewHealthCheck(cfg model.HealthConfig) HealthCheck {
	return &healthCheck{config: cfg}
}

// Check returns the health status report for the service.
func (hc *healthCheck) Check()*model.HealthReport{
	serviceName := hc.config.GetServiceName()
	instanceId := hc.config.GetInstanceId()
	return &model.HealthReport{
					Message:		 "OK", 
					ServiceName:	 serviceName, 
					InstanceId: 	 instanceId,
				}
}