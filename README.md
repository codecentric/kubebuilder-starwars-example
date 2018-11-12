# kubebuilder-starwars-example

That is a demo repository to add a custom CRD to a kubernetes clusters, for example via docker 4 mac or minikube.

## Usage
### Docker 4 Mac or Minikube

* `make install` to install the CRD
* `make run` to let the manager run
* `kubectl apply -f config/sample/ships_v1beta1_starship.yaml` to spawn the first starhship

More Deployment options [here](https://book.kubebuilder.io/basics/project_creation_and_structure.html)