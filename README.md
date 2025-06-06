# godig

**godig** is a lightweight, thread-safe wrapper around [Uber's dig](https://pkg.go.dev/go.uber.org/dig) dependency injection framework. It enforces best practices by validating constructor functions, preventing duplicate registrations, and ensuring safe concurrent access.

## 🚀 Features

* **Constructor Validation**: Accepts functions returning either a single value or a value and an error.
* **Duplicate Prevention**: Prevents multiple constructors for the same type.
* **Post-Resolution Lock**: Disallows further registrations after the first successful resolution.
* **Generic Resolution**: Provides a type-safe `Resolve[T]` function.
* **Thread-Safe**: Utilizes `sync.Mutex` to ensure safe concurrent operations.

## 📦 Installation

To install `godig`, use `go get`:

```bash
go get github.com/limaJavier/godig
```

Then, import it in your Go code:

```go
import "github.com/limaJavier/godig"
```

## 🛠️ Usage

Here's a basic example of how to use `godig`:

```go
package main

import (
	"fmt"

	"github.com/limaJavier/godig"
)

// Vehicle is an interface with two methods.
type Vehicle interface {
	Start()
	Stop()
}

// Car is a concrete implementation of Vehicle.
type Car struct{}

func (c *Car) Start() {
	fmt.Println("Car started")
}

func (c *Car) Stop() {
	fmt.Println("Car stopped")
}

// NewCar is a constructor that returns a Car and possibly an error.
func NewCar() (*Car, error) {
	return &Car{}, nil
}

// NewVehicle returns a Vehicle interface, built from a *Car.
func NewVehicle(car *Car) Vehicle {
	return car
}

func main() {
	resolver := godig.New()

	// Register constructors
	if err := godig.Register(resolver, NewCar, NewVehicle); err != nil {
		panic(err)
	}

	// Resolve the interface
	vehicle, err := godig.Resolve[Vehicle](resolver)
	if err != nil {
		panic(err)
	}
	vehicle.Start()

	// Resolve the concrete implementation
	car, err := godig.Resolve[*Car](resolver)
	if err != nil {
		panic(err)
	}
	car.Stop()
}
```

## 🧪 Testing

To run tests:

```bash
go test ./...
```

Ensure you have Go installed and set up properly.

## 📄 License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
