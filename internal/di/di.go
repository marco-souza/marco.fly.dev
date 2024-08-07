package di

import (
	"context"
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
		logger.Warn("entity is a pointer", "entity", entityPointer, "value", entity)
	}

	t := reflect.TypeOf(entity)
	if val := ctx.Value(t); val != nil {
		logger.Warn("entity already exists", "entity", t)
		return
	}

	// check if entity is a Dependency
	if d, ok := entityPointer.(Service); ok {
		logger.Info("starting service", "entity", t)

		if err := d.Start(); err != nil {
			logger.Warn("failed to start dependency", "entity", t, "err", err)
			return
		}

		// register for teardown
		teardownServices = append(teardownServices, d)
	}

	// register the dependency
	ctx = context.WithValue(ctx, t, entity)
	logger.Info("dependency injected", "dependency", t)
}

func Inject[T any](entity T) *T {
	t := reflect.TypeOf(entity)
	instance := ctx.Value(t)
	if instance == nil {
		return nil
	}

	if t, ok := instance.(T); ok {
		return &t
	}

	return nil
}

// clean the Container
func Clean() {
	// run all teardown functions
	logger.Info("cleaning dependencies", "count", len(teardownServices))
	for _, svc := range teardownServices {
		logger.Info("tearing down", "svc", svc)
		if err := svc.Stop(); err != nil {
			logger.Warn("failed to teardown dependency", "err", err)
		}
	}

	ctx = context.Background()
	teardownServices = []Service{}
}
