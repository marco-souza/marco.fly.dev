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

func Injectables(entities ...interface{}) {
	for _, entity := range entities {
		Injectable(entity)
	}
}

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

func Invoke[F any](cb F) error {
	// get cb arguments into args
	fn := reflect.ValueOf(cb)
	args := []reflect.Type{}
	for i := 0; i < fn.Type().NumIn(); i++ {
		elem := fn.Type().In(i)
		args = append(args, elem)
	}

	// get real values for cb deps
	deps := make([]reflect.Value, len(args))
	for i, arg := range args {
		dep := ctx.Value(arg)
		if dep == nil {
			return fmt.Errorf("dependency not found: %s", arg)
		}

		deps[i] = reflect.ValueOf(dep)
	}

	// call db with deps
	result := reflect.ValueOf(cb).Call(deps)[0]
	if result.IsNil() {
		return nil
	}

	return result.Interface().(error)
}

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

func MustInject[T any](entity T) *T {
	t, err := Inject(entity)
	if err != nil {
		panic(err)
	}

	return t
}

// clean the Container
func Clean() {
	for _, svc := range teardownServices {
		logger.Info("tearing down", "service", reflect.TypeOf(svc))
		if err := svc.Stop(); err != nil {
			logger.Warn("failed to teardown dependency", "err", err)
		}
	}

	ctx = context.Background()
	teardownServices = []Service{}
}
