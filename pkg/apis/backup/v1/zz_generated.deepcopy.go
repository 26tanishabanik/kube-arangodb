//go:build !ignore_autogenerated
// +build !ignore_autogenerated

//
// DISCLAIMER
//
// Copyright 2023 ArangoDB GmbH, Cologne, Germany
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Copyright holder is ArangoDB GmbH, Cologne, Germany
//

// Code generated by deepcopy-gen. DO NOT EDIT.

package v1

import (
	sharedv1 "github.com/arangodb/kube-arangodb/pkg/apis/shared/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ArangoBackup) DeepCopyInto(out *ArangoBackup) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ArangoBackup.
func (in *ArangoBackup) DeepCopy() *ArangoBackup {
	if in == nil {
		return nil
	}
	out := new(ArangoBackup)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ArangoBackup) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ArangoBackupDetails) DeepCopyInto(out *ArangoBackupDetails) {
	*out = *in
	if in.PotentiallyInconsistent != nil {
		in, out := &in.PotentiallyInconsistent, &out.PotentiallyInconsistent
		*out = new(bool)
		**out = **in
	}
	if in.Uploaded != nil {
		in, out := &in.Uploaded, &out.Uploaded
		*out = new(bool)
		**out = **in
	}
	if in.Downloaded != nil {
		in, out := &in.Downloaded, &out.Downloaded
		*out = new(bool)
		**out = **in
	}
	if in.Imported != nil {
		in, out := &in.Imported, &out.Imported
		*out = new(bool)
		**out = **in
	}
	in.CreationTimestamp.DeepCopyInto(&out.CreationTimestamp)
	if in.Keys != nil {
		in, out := &in.Keys, &out.Keys
		*out = make(sharedv1.HashList, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ArangoBackupDetails.
func (in *ArangoBackupDetails) DeepCopy() *ArangoBackupDetails {
	if in == nil {
		return nil
	}
	out := new(ArangoBackupDetails)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ArangoBackupList) DeepCopyInto(out *ArangoBackupList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ArangoBackup, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ArangoBackupList.
func (in *ArangoBackupList) DeepCopy() *ArangoBackupList {
	if in == nil {
		return nil
	}
	out := new(ArangoBackupList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ArangoBackupList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ArangoBackupPolicy) DeepCopyInto(out *ArangoBackupPolicy) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ArangoBackupPolicy.
func (in *ArangoBackupPolicy) DeepCopy() *ArangoBackupPolicy {
	if in == nil {
		return nil
	}
	out := new(ArangoBackupPolicy)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ArangoBackupPolicy) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ArangoBackupPolicyList) DeepCopyInto(out *ArangoBackupPolicyList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ArangoBackupPolicy, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ArangoBackupPolicyList.
func (in *ArangoBackupPolicyList) DeepCopy() *ArangoBackupPolicyList {
	if in == nil {
		return nil
	}
	out := new(ArangoBackupPolicyList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ArangoBackupPolicyList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ArangoBackupPolicySpec) DeepCopyInto(out *ArangoBackupPolicySpec) {
	*out = *in
	if in.DeploymentSelector != nil {
		in, out := &in.DeploymentSelector, &out.DeploymentSelector
		*out = new(metav1.LabelSelector)
		(*in).DeepCopyInto(*out)
	}
	in.BackupTemplate.DeepCopyInto(&out.BackupTemplate)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ArangoBackupPolicySpec.
func (in *ArangoBackupPolicySpec) DeepCopy() *ArangoBackupPolicySpec {
	if in == nil {
		return nil
	}
	out := new(ArangoBackupPolicySpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ArangoBackupPolicyStatus) DeepCopyInto(out *ArangoBackupPolicyStatus) {
	*out = *in
	in.Scheduled.DeepCopyInto(&out.Scheduled)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ArangoBackupPolicyStatus.
func (in *ArangoBackupPolicyStatus) DeepCopy() *ArangoBackupPolicyStatus {
	if in == nil {
		return nil
	}
	out := new(ArangoBackupPolicyStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ArangoBackupProgress) DeepCopyInto(out *ArangoBackupProgress) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ArangoBackupProgress.
func (in *ArangoBackupProgress) DeepCopy() *ArangoBackupProgress {
	if in == nil {
		return nil
	}
	out := new(ArangoBackupProgress)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ArangoBackupSpec) DeepCopyInto(out *ArangoBackupSpec) {
	*out = *in
	out.Deployment = in.Deployment
	if in.Options != nil {
		in, out := &in.Options, &out.Options
		*out = new(ArangoBackupSpecOptions)
		(*in).DeepCopyInto(*out)
	}
	if in.Download != nil {
		in, out := &in.Download, &out.Download
		*out = new(ArangoBackupSpecDownload)
		**out = **in
	}
	if in.Upload != nil {
		in, out := &in.Upload, &out.Upload
		*out = new(ArangoBackupSpecOperation)
		**out = **in
	}
	if in.PolicyName != nil {
		in, out := &in.PolicyName, &out.PolicyName
		*out = new(string)
		**out = **in
	}
	if in.Backoff != nil {
		in, out := &in.Backoff, &out.Backoff
		*out = new(ArangoBackupSpecBackOff)
		(*in).DeepCopyInto(*out)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ArangoBackupSpec.
func (in *ArangoBackupSpec) DeepCopy() *ArangoBackupSpec {
	if in == nil {
		return nil
	}
	out := new(ArangoBackupSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ArangoBackupSpecBackOff) DeepCopyInto(out *ArangoBackupSpecBackOff) {
	*out = *in
	if in.MinDelay != nil {
		in, out := &in.MinDelay, &out.MinDelay
		*out = new(int)
		**out = **in
	}
	if in.MaxDelay != nil {
		in, out := &in.MaxDelay, &out.MaxDelay
		*out = new(int)
		**out = **in
	}
	if in.Iterations != nil {
		in, out := &in.Iterations, &out.Iterations
		*out = new(int)
		**out = **in
	}
	if in.MaxIterations != nil {
		in, out := &in.MaxIterations, &out.MaxIterations
		*out = new(int)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ArangoBackupSpecBackOff.
func (in *ArangoBackupSpecBackOff) DeepCopy() *ArangoBackupSpecBackOff {
	if in == nil {
		return nil
	}
	out := new(ArangoBackupSpecBackOff)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ArangoBackupSpecDeployment) DeepCopyInto(out *ArangoBackupSpecDeployment) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ArangoBackupSpecDeployment.
func (in *ArangoBackupSpecDeployment) DeepCopy() *ArangoBackupSpecDeployment {
	if in == nil {
		return nil
	}
	out := new(ArangoBackupSpecDeployment)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ArangoBackupSpecDownload) DeepCopyInto(out *ArangoBackupSpecDownload) {
	*out = *in
	out.ArangoBackupSpecOperation = in.ArangoBackupSpecOperation
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ArangoBackupSpecDownload.
func (in *ArangoBackupSpecDownload) DeepCopy() *ArangoBackupSpecDownload {
	if in == nil {
		return nil
	}
	out := new(ArangoBackupSpecDownload)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ArangoBackupSpecOperation) DeepCopyInto(out *ArangoBackupSpecOperation) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ArangoBackupSpecOperation.
func (in *ArangoBackupSpecOperation) DeepCopy() *ArangoBackupSpecOperation {
	if in == nil {
		return nil
	}
	out := new(ArangoBackupSpecOperation)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ArangoBackupSpecOptions) DeepCopyInto(out *ArangoBackupSpecOptions) {
	*out = *in
	if in.Timeout != nil {
		in, out := &in.Timeout, &out.Timeout
		*out = new(float32)
		**out = **in
	}
	if in.AllowInconsistent != nil {
		in, out := &in.AllowInconsistent, &out.AllowInconsistent
		*out = new(bool)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ArangoBackupSpecOptions.
func (in *ArangoBackupSpecOptions) DeepCopy() *ArangoBackupSpecOptions {
	if in == nil {
		return nil
	}
	out := new(ArangoBackupSpecOptions)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ArangoBackupState) DeepCopyInto(out *ArangoBackupState) {
	*out = *in
	in.Time.DeepCopyInto(&out.Time)
	if in.Progress != nil {
		in, out := &in.Progress, &out.Progress
		*out = new(ArangoBackupProgress)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ArangoBackupState.
func (in *ArangoBackupState) DeepCopy() *ArangoBackupState {
	if in == nil {
		return nil
	}
	out := new(ArangoBackupState)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ArangoBackupStatus) DeepCopyInto(out *ArangoBackupStatus) {
	*out = *in
	in.ArangoBackupState.DeepCopyInto(&out.ArangoBackupState)
	if in.Backup != nil {
		in, out := &in.Backup, &out.Backup
		*out = new(ArangoBackupDetails)
		(*in).DeepCopyInto(*out)
	}
	if in.Backoff != nil {
		in, out := &in.Backoff, &out.Backoff
		*out = new(ArangoBackupStatusBackOff)
		(*in).DeepCopyInto(*out)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ArangoBackupStatus.
func (in *ArangoBackupStatus) DeepCopy() *ArangoBackupStatus {
	if in == nil {
		return nil
	}
	out := new(ArangoBackupStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ArangoBackupStatusBackOff) DeepCopyInto(out *ArangoBackupStatusBackOff) {
	*out = *in
	in.Next.DeepCopyInto(&out.Next)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ArangoBackupStatusBackOff.
func (in *ArangoBackupStatusBackOff) DeepCopy() *ArangoBackupStatusBackOff {
	if in == nil {
		return nil
	}
	out := new(ArangoBackupStatusBackOff)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ArangoBackupTemplate) DeepCopyInto(out *ArangoBackupTemplate) {
	*out = *in
	if in.Options != nil {
		in, out := &in.Options, &out.Options
		*out = new(ArangoBackupSpecOptions)
		(*in).DeepCopyInto(*out)
	}
	if in.Upload != nil {
		in, out := &in.Upload, &out.Upload
		*out = new(ArangoBackupSpecOperation)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ArangoBackupTemplate.
func (in *ArangoBackupTemplate) DeepCopy() *ArangoBackupTemplate {
	if in == nil {
		return nil
	}
	out := new(ArangoBackupTemplate)
	in.DeepCopyInto(out)
	return out
}
