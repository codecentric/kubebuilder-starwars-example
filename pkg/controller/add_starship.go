package controller

import (
	"github.com/codecentric/kubebuilder-starwars-example/pkg/controller/starship"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, starship.Add)
}
