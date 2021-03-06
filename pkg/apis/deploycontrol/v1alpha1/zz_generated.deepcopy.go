// +build !ignore_autogenerated

/*
Copyright The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by deepcopy-gen. DO NOT EDIT.

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterSpec) DeepCopyInto(out *ClusterSpec) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterSpec.
func (in *ClusterSpec) DeepCopy() *ClusterSpec {
	if in == nil {
		return nil
	}
	out := new(ClusterSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ConditionManager) DeepCopyInto(out *ConditionManager) {
	*out = *in
	in.status.DeepCopyInto(&out.status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ConditionManager.
func (in *ConditionManager) DeepCopy() *ConditionManager {
	if in == nil {
		return nil
	}
	out := new(ConditionManager)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ConditionsSpec) DeepCopyInto(out *ConditionsSpec) {
	*out = *in
	in.LastUpdateTime.DeepCopyInto(&out.LastUpdateTime)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ConditionsSpec.
func (in *ConditionsSpec) DeepCopy() *ConditionsSpec {
	if in == nil {
		return nil
	}
	out := new(ConditionsSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DeployDaemon) DeepCopyInto(out *DeployDaemon) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	if in.Status != nil {
		in, out := &in.Status, &out.Status
		*out = new(DeploydaemonStatus)
		(*in).DeepCopyInto(*out)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DeployDaemon.
func (in *DeployDaemon) DeepCopy() *DeployDaemon {
	if in == nil {
		return nil
	}
	out := new(DeployDaemon)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *DeployDaemon) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DeployDaemonList) DeepCopyInto(out *DeployDaemonList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]DeployDaemon, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DeployDaemonList.
func (in *DeployDaemonList) DeepCopy() *DeployDaemonList {
	if in == nil {
		return nil
	}
	out := new(DeployDaemonList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *DeployDaemonList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DeploydaemonSpec) DeepCopyInto(out *DeploydaemonSpec) {
	*out = *in
	if in.Secrets != nil {
		in, out := &in.Secrets, &out.Secrets
		*out = make([]SecretsRef, len(*in))
		copy(*out, *in)
	}
	if in.Replica != nil {
		in, out := &in.Replica, &out.Replica
		*out = new(int32)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DeploydaemonSpec.
func (in *DeploydaemonSpec) DeepCopy() *DeploydaemonSpec {
	if in == nil {
		return nil
	}
	out := new(DeploydaemonSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DeploydaemonStatus) DeepCopyInto(out *DeploydaemonStatus) {
	*out = *in
	if in.Cluster != nil {
		in, out := &in.Cluster, &out.Cluster
		*out = new(ClusterSpec)
		**out = **in
	}
	if in.StartTime != nil {
		in, out := &in.StartTime, &out.StartTime
		*out = (*in).DeepCopy()
	}
	if in.CompletionTime != nil {
		in, out := &in.CompletionTime, &out.CompletionTime
		*out = (*in).DeepCopy()
	}
	in.Conditions.DeepCopyInto(&out.Conditions)
	in.Deployment.DeepCopyInto(&out.Deployment)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DeploydaemonStatus.
func (in *DeploydaemonStatus) DeepCopy() *DeploydaemonStatus {
	if in == nil {
		return nil
	}
	out := new(DeploydaemonStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SecretsRef) DeepCopyInto(out *SecretsRef) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SecretsRef.
func (in *SecretsRef) DeepCopy() *SecretsRef {
	if in == nil {
		return nil
	}
	out := new(SecretsRef)
	in.DeepCopyInto(out)
	return out
}
