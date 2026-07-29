package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/khaivutri/bookmark-service/internal/model"
	"github.com/khaivutri/bookmark-service/internal/repository"
)

// HealthCheck defines the interface for the health check service.
type HealthCheck interface {
	// Check verifies the connectivity of Redis and database adapters.
	Check(ctx context.Context) (*model.HealthReport, error)
}
type redisPinger interface {
	Ping(ctx context.Context) error
}

type dbPinger interface {
	Ping(ctx context.Context) error
}

type healthCheck struct {
	serviceName 	string
	instanceID  	string
	
	redisAdapter 	redisPinger
	dbAdapter    	dbPinger
}

// NewHealthCheck constructs a new HealthCheck service instance.
func NewHealthCheck(serviceName, instanceID string, redisAdapter redisPinger, dbAdapter dbPinger) HealthCheck {
	return &healthCheck{	
							serviceName: 	serviceName,
							instanceID:		instanceID,
							redisAdapter: 	redisAdapter,
							dbAdapter:    	dbAdapter,
						}
}

// Check verifies the connectivity of Redis and database adapters.
func (hc *healthCheck) Check(ctx context.Context) (*model.HealthReport, error) {
	dpc := make(map[string]string)

	var msg string = "OK"
	var firstErr error

	pingCtx, cancel := context.WithTimeout(ctx, 5 * time.Second)
	defer cancel()

	if pingErr := hc.redisAdapter.Ping(pingCtx); pingErr != nil {
		dpc["redis"] = "DOWN"
		msg = "DEGRADED"
		firstErr = fmt.Errorf("%w: redis: %v", repository.ErrDependencyDown, pingErr)
	} else {
		dpc["redis"] = "UP"
	}

	if pingErr := hc.dbAdapter.Ping(pingCtx); pingErr != nil {
		dpc["postgres"] = "DOWN"
		msg = "DEGRADED"
		pgErr := fmt.Errorf("%w: postgres: %w", repository.ErrDependencyDown, pingErr)

		if firstErr == nil {
			firstErr = pgErr
		} else {
			firstErr = errors.Join(firstErr, pgErr)
		}
	} else {
		dpc["postgres"] = "UP"
	}

	return &model.HealthReport{	Message: msg, 
								ServiceName: hc.serviceName, 
								InstanceID: hc.instanceID, 
								Dependencies: dpc,
							}, firstErr

}