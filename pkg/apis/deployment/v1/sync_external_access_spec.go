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
	"net"
	"net/url"
	"strconv"

	"github.com/arangodb/kube-arangodb/pkg/apis/shared"
	"github.com/arangodb/kube-arangodb/pkg/util/errors"
)

// SyncExternalAccessSpec holds configuration for the external access provided for the sync deployment.
type SyncExternalAccessSpec struct {
	ExternalAccessSpec
	// MasterEndpoint setting specifies the master endpoint(s) advertised by the ArangoSync SyncMasters.
	// If not set, this setting defaults to:
	// - If `spec.sync.externalAccess.loadBalancerIP` is set, it defaults to `https://<load-balancer-ip>:<8629>`.
	// - Otherwise it defaults to `https://<sync-service-dns-name>:<8629>`.
	// +doc/type: []string
	MasterEndpoint []string `json:"masterEndpoint,omitempty"`
	// AccessPackageSecretNames setting specifies the names of zero of more `Secrets` that will be created by the deployment
	// operator containing "access packages". An access package contains those `Secrets` that are needed
	// to access the SyncMasters of this `ArangoDeployment`.
	// By removing a name from this setting, the corresponding `Secret` is also deleted.
	// Note that to remove all access packages, leave an empty array in place (`[]`).
	// Completely removing the setting results in not modifying the list.
	// +doc/type: []string
	// +doc/link: See the ArangoDeploymentReplication specification|deployment-replication-resource-reference.md
	AccessPackageSecretNames []string `json:"accessPackageSecretNames,omitempty"`
}

// GetMasterEndpoint returns the value of masterEndpoint.
func (s SyncExternalAccessSpec) GetMasterEndpoint() []string {
	return s.MasterEndpoint
}

// GetAccessPackageSecretNames returns the value of accessPackageSecretNames.
func (s SyncExternalAccessSpec) GetAccessPackageSecretNames() []string {
	return s.AccessPackageSecretNames
}

// ResolveMasterEndpoint returns the value of `--master.endpoint` option passed to arangosync.
func (s SyncExternalAccessSpec) ResolveMasterEndpoint(syncServiceHostName string, syncServicePort int) []string {
	if len(s.MasterEndpoint) > 0 {
		return s.MasterEndpoint
	}
	if ip := s.GetLoadBalancerIP(); ip != "" {
		syncServiceHostName = ip
	}
	return []string{"https://" + net.JoinHostPort(syncServiceHostName, strconv.Itoa(syncServicePort))}
}

// Validate the given spec
func (s SyncExternalAccessSpec) Validate() error {
	if err := s.ExternalAccessSpec.Validate(); err != nil {
		return errors.WithStack(err)
	}
	for _, ep := range s.MasterEndpoint {
		if u, err := url.Parse(ep); err != nil {
			return errors.WithStack(errors.Newf("Failed to parse master endpoint '%s': %s", ep, err))
		} else {
			if u.Scheme != "http" && u.Scheme != "https" {
				return errors.WithStack(errors.Newf("Invalid scheme '%s' in master endpoint '%s'", u.Scheme, ep))
			}
			if u.Host == "" {
				return errors.WithStack(errors.Newf("Missing host in master endpoint '%s'", ep))
			}
		}
	}
	for _, name := range s.AccessPackageSecretNames {
		if err := shared.ValidateResourceName(name); err != nil {
			return errors.WithStack(errors.Newf("Invalid name '%s' in accessPackageSecretNames: %s", name, err))
		}
	}
	return nil
}

// SetDefaults fills in missing defaults
func (s *SyncExternalAccessSpec) SetDefaults() {
	s.ExternalAccessSpec.SetDefaults()
}

// SetDefaultsFrom fills unspecified fields with a value from given source spec.
func (s *SyncExternalAccessSpec) SetDefaultsFrom(source SyncExternalAccessSpec) {
	s.ExternalAccessSpec.SetDefaultsFrom(source.ExternalAccessSpec)
	if s.MasterEndpoint == nil && source.MasterEndpoint != nil {
		s.MasterEndpoint = append([]string{}, source.MasterEndpoint...)
	}
	if s.AccessPackageSecretNames == nil && source.AccessPackageSecretNames != nil {
		s.AccessPackageSecretNames = append([]string{}, source.AccessPackageSecretNames...)
	}
}

// ResetImmutableFields replaces all immutable fields in the given target with values from the source spec.
// It returns a list of fields that have been reset.
// Field names are relative to given field prefix.
func (s SyncExternalAccessSpec) ResetImmutableFields(fieldPrefix string, target *SyncExternalAccessSpec) []string {
	result := s.ExternalAccessSpec.ResetImmutableFields(fieldPrefix, &s.ExternalAccessSpec)
	return result
}
