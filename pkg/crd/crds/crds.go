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

package crds

import (
	"fmt"

	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/util/yaml"

	"github.com/arangodb/go-driver"
	"github.com/arangodb/kube-arangodb/pkg/util"
)

type Definition struct {
	Version       driver.Version
	CRD           *apiextensions.CustomResourceDefinition
	CRDWithSchema *apiextensions.CustomResourceDefinition
}

func AllDefinitions() []Definition {
	return []Definition{
		// Deployment
		DatabaseDeploymentDefinition(),
		DatabaseMemberDefinition(),

		// ACS
		DatabaseClusterSynchronizationDefinition(),

		// ArangoSync
		ReplicationDeploymentReplicationDefinition(),

		// Storage
		StorageLocalStorageDefinition(),

		// Apps
		AppsJobDefinition(),
		DatabaseTaskDefinition(),

		// Backups
		BackupsBackupDefinition(),
		BackupsBackupPolicyDefinition(),

		// ML
		MLExtensionDefinition(),
		MLStorageDefinition(),

		MLCronJobDefinition(),
		MLBatchJobDefinition(),
	}
}

func mustLoadCRD(crdRaw, crdSchemasRaw []byte, crd, crdWithSchema *apiextensions.CustomResourceDefinition) {
	if err := yaml.Unmarshal(crdRaw, crd); err != nil {
		panic(err)
	}

	var crdSchemas map[string]apiextensions.CustomResourceValidation
	if err := yaml.Unmarshal(crdSchemasRaw, &crdSchemas); err != nil {
		panic(err)
	}

	aCopy := crd.DeepCopy()
	*crdWithSchema = *aCopy
	for i, v := range crdWithSchema.Spec.Versions {
		schema, ok := crdSchemas[v.Name]
		if !ok {
			panic(fmt.Sprintf("Validation schema is not defined for version %s of %s", v.Name, crd.Name))
		}
		crdWithSchema.Spec.Versions[i].Schema = schema.DeepCopy()
	}
}

type GetCRDOptions struct {
	WithSchema *bool
}

func WithSchema() GetCRDOptions {
	return GetCRDOptions{
		WithSchema: util.NewType(true),
	}
}

func getCRD(crd, crdWithSchema apiextensions.CustomResourceDefinition, opts ...GetCRDOptions) *apiextensions.CustomResourceDefinition {
	var o GetCRDOptions
	if len(opts) > 0 {
		o = opts[0]
	}
	if o.WithSchema != nil && *o.WithSchema {
		return crdWithSchema.DeepCopy()
	}
	return crd.DeepCopy()
}
