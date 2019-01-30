package v1alpha1

import (
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
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

	Status *DeploydaemonStatus `json:"status"`
}

func (d DeployDaemon) String() string {
	var _output string

	_output+=fmt.Sprintln("Name: ", d.Name)
	_output+=fmt.Sprintln("Component: ", d.Spec.Component)
	_output+=fmt.Sprintln("Image: ", d.Spec.Image)
	_output+=fmt.Sprintln("Version: ", d.Spec.Version)
	_output+=fmt.Sprintln("Expose: ", d.Spec.Expose)
	_output+=fmt.Sprintln("ConfigRef: ", d.Spec.Config)

	return _output
}

func (d *DeployDaemon) GetDeploymentName() string {

	return d.Spec.Tenant+d.Spec.Environment+d.Spec.EnvType+"-"+d.Spec.Component
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
	Tenant string `json:"tenant"`
	Environment string `json:"environment"`
	EnvType   string   `json:"envtype"`
	Scheduler string `json:"scheduler,omitempty"`
    Image     string  `json:"image"`
	Version   string `json:"version"`
	Config    string `json:"configRef"`
	Secrets   []SecretsRef  `json:"secretRefs"'`
	Expose    string `json:"expose"`
	Replica   *int32 `json:"instance"`
}



// Define Secret Reference, Support associate with multiple secret as environment parameter
type SecretsRef struct{
	Name string `json:"paramName"`
	Secret string `json:"secretName"`
}

type DeploydaemonStatus struct {

	Cluster *ClusterSpec `json:"cluster,omitempty"`

	// StartTime is the time the build is actually started.
	// +optional
	StartTime *metav1.Time `json:"startTime,omitempty"`

	// CompletionTime is the time the build completed.
	// +optional
	CompletionTime *metav1.Time `json:"completionTime,omitempty"`

	// Define Current Deploy Daemon Status
	Conditions ConditionsSpec `json:"conditions,omitempty"`

	// Define Deployment Status
	Deployment appsv1.DeploymentStatus `json:"deploymentStatus,omitempty"`
}

type ClusterSpec struct{
	Name      string   `json:"name,omitempty"`
	NameSpace string `json:"namespace,omitempty"`
	DeploymentName string `json:"deployment,omitempty"`
}

type ConditionsSpec struct{
	LastUpdateTime metav1.Time `json:"lastupdatetime,omitempty"`
	Type           string      `json:"type,omitempty"`
	Status         bool      `json:"status,omitempty"`
	Reason         string    `json:"reason,omitempty"`
	Message        string   `json:"message,omitempty"`
}


type ConditionManager struct {
   status DeploydaemonStatus
}


func (condM ConditionManager) UpdateStatus(dd *DeployDaemon, dp *appsv1.Deployment) {
	//here we update or create new status for DeployDaemon
}

var condManager ConditionManager

func init(){
	condManager = ConditionManager{}
}





