package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type DeployDaemon struct {

	// TypeMeta is the metadata for the resource, like kind and apiversion
	metav1.TypeMeta `json:",inline"`

	// ObjectMeta contains the metadata for the particular object, including
	// things like...
	//  - name
	//  - namespace
	//  - self link
	//  - labels
	//  - ... etc ...
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec DeploydaemonSpec `json:"spec"`

	Status DeploydaemonStatus `json:"Status"`
}


// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// DeployDaemonList is a list of DeployDaemon resources
type DeployDaemonList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`
	Items []DeployDaemon `json:"items"`
}


type DeploydaemonSpec struct {
	Component string `json:"component"`
    Image     string  `json:"image"`
	Version   string `json:"version"`
	Config    string `json:"configRef"`
	Secrets   []SecretsRef  `json:"secretRefs"'`
	Replica   int `json:"instanceNum"`
}


// Define Secret Reference, Support associate with multiple secret as environment parameter
type SecretsRef struct{
	Name string `json:"paramName"`
	Secret string `json:"secretName"`
}

type DeploydaemonStatus struct {

	Cluster ClusterSpec `json:"cluster,omitempty"`

	NameSpace string `json:"namespace,omitempty"`

	// StartTime is the time the build is actually started.
	// +optional
	StartTime *metav1.Time `json:"startTime,omitempty"`

	// CompletionTime is the time the build completed.
	// +optional
	CompletionTime *metav1.Time `json:"completionTime,omitempty"`

	// Define Status
	Status StatusSpec `json:"status"`
}

type ClusterSpec struct{
	Name  string   `json:"name,omitempty"`
	NameSpace string `json:"namespace,omitempty"`
}

type StatusSpec struct{
	Success bool     `json:"success"`
	Message string   `json:"message"`
}


