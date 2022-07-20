package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
)

type MyResourceSpec struct {
	FieldA string `json:"field_a,omitempty"`
	FieldB string `json:"field_b,omitempty"`
}

type MyResourceStatus struct {
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

type MyResource struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MyResourceSpec   `json:"spec,omitempty"`
	Status MyResourceStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

type MyResourceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MyResource `json:"items"`
}

func init() {
	// TODO clean up

	// GroupVersion is group version used to register these objects
	groupVersion := schema.GroupVersion{Group: "application.sample.ibm.com", Version: "v1alpha1"}

	// SchemeBuilder is used to add go types to the GroupVersionKind scheme
	schemeBuilder := &scheme.Builder{GroupVersion: groupVersion}

	// AddToScheme adds the types in this group-version to the given scheme.
	addToScheme := schemeBuilder.AddToScheme
	var _ = addToScheme

	schemeBuilder.Register(&MyResource{}, &MyResourceList{})
}
