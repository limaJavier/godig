package godig

import (
	"fmt"
	"reflect"
	"sync"

	"go.uber.org/dig"
)

type Resolver struct {
	mutex           sync.Mutex
	container       *dig.Container
	registeredTypes set[reflect.Type]
	resolved        bool
}

func New() *Resolver {
	return &Resolver{
		mutex:           sync.Mutex{},
		container:       dig.New(),
		registeredTypes: make(set[reflect.Type]),
	}
}

func Register(resolver *Resolver, constructors ...any) error {
	resolver.mutex.Lock()
	defer resolver.mutex.Unlock()
	if resolver.resolved {
		return fmt.Errorf("cannot register anymore constructors since a type has already been resolved")
	}

	for _, constructor := range constructors {
		constructorType := reflect.TypeOf(constructor)
		constructorOutCount := constructorType.NumOut()
		if constructorType.Kind() != reflect.Func {
			return fmt.Errorf("provided constructor is not a function: constructor of type \"%v\"", constructorType.Kind())
		}
		if constructorOutCount < 1 || constructorOutCount > 2 {
			return fmt.Errorf("constructor must return either 1 output or 2 outputs where the second one is of type \"error\": constructor has %v outputs", constructorOutCount)
		}
		if constructorOutCount == 2 && constructorType.Out(1) != reflect.TypeOf((*error)(nil)).Elem() {
			secondOutputTypeName := constructorType.Out(1).Name()
			if secondOutputTypeName != "" {
				secondOutputTypeName = fmt.Sprintf("(%v)", secondOutputTypeName)
			}
			return fmt.Errorf("second return value must be of type \"error\": second return value is of type \"%v%v\"", constructorType.Out(1).Kind(), secondOutputTypeName)
		}

		outputType := constructorType.Out(0)
		if resolver.registeredTypes.Contains(outputType) {
			outputTypeName := outputType.Name()
			if outputTypeName != "" {
				outputTypeName = fmt.Sprintf("(%v)", outputTypeName)
			}
			return fmt.Errorf("cannot register more than one constructor per type: type \"%v%v\" has two registered constructors", outputType.Kind(), outputTypeName)
		} else {
			resolver.registeredTypes.Add(outputType)
		}

		if err := resolver.container.Provide(constructor); err != nil {
			return err
		}
	}
	return nil
}

func Resolve[T any](resolver *Resolver) (T, error) {
	var instance T
	if err := resolver.container.Invoke(func(obj T) { instance = obj }); err != nil {
		return instance, err
	}

	resolver.mutex.Lock()
	resolver.resolved = true
	resolver.mutex.Unlock()

	return instance, nil
}
