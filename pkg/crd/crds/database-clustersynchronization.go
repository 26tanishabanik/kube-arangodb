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
	_ "embed"

	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"

	"github.com/arangodb/go-driver"
)

const (
	DatabaseClusterSynchronizationVersion = driver.Version("1.0.1")
)

func init() {
	mustLoadCRD(databaseClusterSynchronization, databaseClusterSynchronizationSchemaRaw, &databaseClusterSynchronizationCRD, &databaseClusterSynchronizationCRDSchemas)
}

// Deprecated: use DatabaseClusterSynchronizationWithOptions instead
func DatabaseClusterSynchronization() *apiextensions.CustomResourceDefinition {
	return DatabaseClusterSynchronizationWithOptions()
}

func DatabaseClusterSynchronizationWithOptions(opts ...func(*CRDOptions)) *apiextensions.CustomResourceDefinition {
	return getCRD(databaseClusterSynchronizationCRD, databaseClusterSynchronizationCRDSchemas, opts...)
}

// Deprecated: use DatabaseClusterSynchronizationDefinitionWithOptions instead
func DatabaseClusterSynchronizationDefinition() Definition {
	return DatabaseClusterSynchronizationDefinitionWithOptions()
}

func DatabaseClusterSynchronizationDefinitionWithOptions(opts ...func(*CRDOptions)) Definition {
	return Definition{
		Version: DatabaseClusterSynchronizationVersion,
		CRD:     DatabaseClusterSynchronizationWithOptions(opts...),
	}
}

var databaseClusterSynchronizationCRD apiextensions.CustomResourceDefinition
var databaseClusterSynchronizationCRDSchemas crdSchemas

//go:embed database-clustersynchronization.yaml
var databaseClusterSynchronization []byte

//go:embed database-clustersynchronization.schema.generated.yaml
var databaseClusterSynchronizationSchemaRaw []byte
