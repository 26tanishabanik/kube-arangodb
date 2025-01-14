//
// DISCLAIMER
//
// Copyright 2016-2023 ArangoDB GmbH, Cologne, Germany
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

package v1

import (
	core "k8s.io/api/core/v1"
)

type LifecycleSpec struct {
	// Resources holds resource requests & limits
	// +doc/type: core.ResourceRequirements
	// +doc/link: Documentation of core.ResourceRequirements|https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.26/#resourcerequirements-v1-core
	Resources core.ResourceRequirements `json:"resources,omitempty"`
}

// SetDefaultsFrom fills unspecified fields with a value from given source spec.
func (s *LifecycleSpec) SetDefaultsFrom(source LifecycleSpec) {
	setStorageDefaultsFromResourceList(&s.Resources.Limits, source.Resources.Limits)
	setStorageDefaultsFromResourceList(&s.Resources.Requests, source.Resources.Requests)
}
