package model

// HealthConfig defines the interface for configuration access.
type HealthConfig interface {
	GetServiceName() 	string
	GetInstanceId() 	string
}

// HealthReport represents the health status response.
type HealthReport struct {
	Message 		string 		`json:"message"`
	ServiceName 	string 		`json:"service_name"`
	InstanceID 		string		`json:"instance_id"`
}