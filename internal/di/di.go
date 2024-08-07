package di

import (
	"context"
	"fmt"
	"log/slog"
	"reflect"
)

type Service interface {
	Start() error
	Stop() error
}

var logger = slog.With("service", "di")
var ctx = context.Background()
var teardownServices = []Service{}

func Injectable(entityPointer interface{}) {
	if entityPointer == nil {
		logger.Warn("entity is nil")
		return
	}

	// if poniter, get value
	entity := entityPointer
	if reflect.TypeOf(entityPointer).Kind() == reflect.Ptr {
		entity = reflect.ValueOf(entityPointer).Elem().Interface()
	}

	t := reflect.TypeOf(entity)
	if val := ctx.Value(t); val != nil {
		logger.Warn("entity already exists", "entity", t)
		return
	}

	if d, ok := entityPointer.(Service); ok {
		logger.Info("injecting", "service", t)

		if err := d.Start(); err != nil {
			logger.Warn("failed to start service", "service", t, "err", err)
		}

		// register for teardown
		teardownServices = append(teardownServices, d)
	}

	// register the dependency
	ctx = context.WithValue(ctx, t, entity)
}

// TODO: Implement `Invoke` method to inject dependencies automatically

func Inject[T any](entity T) (*T, error) {
	t := reflect.TypeOf(entity)
	instance := ctx.Value(t)
	if instance == nil {
		return nil, fmt.Errorf("dependency not found: %s", t)
	}

	if t, ok := instance.(T); ok {
		return &t, nil
	}

	return nil, fmt.Errorf("failed to cast dependency: %s", t)
}

// clean the Container
func Clean() {
	for _, svc := range teardownServices {
		logger.Info("tearing down", "service", svc)
		if err := svc.Stop(); err != nil {
			logger.Warn("failed to teardown dependency", "err", err)
		}
	}

	ctx = context.Background()
	teardownServices = []Service{}
}
