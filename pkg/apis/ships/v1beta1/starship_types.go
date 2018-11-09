package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// StarshipSpec defines the desired state of Starship
type StarshipSpec struct {
	Name string `json:"name,omitempty"`
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// StarshipStatus defines the observed state of Starship
type StarshipStatus struct {
	Name       string `json:"name,omitempty"`
	Model      string `json:"model,omitempty"`
	Crew       string `json:"crew,omitempty"`
	Passengers string `json:"passengers,omitempty"`
	Costs      string `json:"costs,omitempty"`
	Capacity   string `json:"capacity,omitempty"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Starship is the Schema for the starships API
// +k8s:openapi-gen=true
type Starship struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   StarshipSpec   `json:"spec,omitempty"`
	Status StarshipStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// StarshipList contains a list of Starship
type StarshipList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Starship `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Starship{}, &StarshipList{})
}
