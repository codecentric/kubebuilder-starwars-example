package starship

import (
	"context"
	"encoding/json"
	"fmt"
	shipsv1beta1 "github.com/codecentric/kubebuilder-starwars-example/pkg/apis/ships/v1beta1"
	"io/ioutil"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

type apiResponse struct {
	Count    int            `json:"count"`
	Next     string         `json:"next"`
	Previous string         `json:"previous,"`
	Results  []shipResponse `json:"results"`
}

type shipResponse struct {
	Name         string `json:"name"`
	Model        string `json:"count"`
	Manufacturer string `json:"manufacturer"`
	Costs        string `json:"cost_in_credits"`
	Passengers   string `json:"passengers"`
	Crew         string `json:"crew"`
	Capacity     string `json:"cargo_capacity"`
}

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new Starship Controller and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
// USER ACTION REQUIRED: update cmd/manager/main.go to call this ships.Add(mgr) to install this Controller
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileStarship{Client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("starship-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to Starship
	err = c.Watch(&source.Kind{Type: &shipsv1beta1.Starship{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create
	// Uncomment watch a Deployment created by Starship - change this for objects you create
	err = c.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &shipsv1beta1.Starship{},
	})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileStarship{}

// ReconcileStarship reconciles a Starship object
type ReconcileStarship struct {
	client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a Starship object and makes changes based on the state read
// and what is in the Starship.Spec
// a Deployment as an example
// Automatically generate RBAC rules to allow the Controller to read and write Deployments
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=ships.codecentric.de,resources=starships,verbs=get;list;watch;create;update;patch;delete
func (r *ReconcileStarship) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	// Fetch the Starship instance
	instance := &shipsv1beta1.Starship{}
	err := r.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Object not found, return.  Created objects are automatically garbage collected.
			// For additional cleanup logic use finalizers.
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	var foundShip shipResponse
	var found bool

	url := "https://swapi.co/api/starships/"

	for !found {
		data := getData(url)

		ship := findShip(data.Results, instance.Spec.Name)

		if ship != nil {
			foundShip = *ship
			found = true
		}

		if ship == nil && data.Next != "" {
			url = data.Next
		}

		if !found && data.Next == "" {
			return reconcile.Result{}, errors.NewBadRequest("Invalid Ship Name given")
		}
	}

	instance.Spec.Name = foundShip.Name
	instance.Status.Capacity = foundShip.Capacity
	instance.Status.Costs = foundShip.Costs
	instance.Status.Crew = foundShip.Crew
	instance.Status.Model = foundShip.Model
	instance.Status.Passengers = foundShip.Passengers

	err = r.Update(context.TODO(), instance)

	if err != nil {
		return reconcile.Result{}, err
	}

	return reconcile.Result{}, nil
}

func getData(url string) apiResponse {
	response, err := http.Get(url)
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	}

	data, _ := ioutil.ReadAll(response.Body)

	var result apiResponse
	err = json.Unmarshal(data, &result)

	return result
}

func findShip(ships []shipResponse, shipname string) *shipResponse {
	for _, ship := range ships {
		if ship.Name == shipname {
			return &ship
		}
	}
	return nil
}
