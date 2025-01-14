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
	DatabaseDeploymentVersion = driver.Version("1.0.1")
)

func init() {
	mustLoadCRD(databaseDeployment, databaseDeploymentSchemaRaw, &databaseDeploymentCRD, &databaseDeploymentCRDSchemas)
}

// Deprecated: use DatabaseDeploymentWithOptions instead
func DatabaseDeployment() *apiextensions.CustomResourceDefinition {
	return DatabaseDeploymentWithOptions()
}

func DatabaseDeploymentWithOptions(opts ...func(*CRDOptions)) *apiextensions.CustomResourceDefinition {
	return getCRD(databaseDeploymentCRD, databaseDeploymentCRDSchemas, opts...)
}

// Deprecated: use DatabaseDeploymentDefinitionWithOptions instead
func DatabaseDeploymentDefinition() Definition {
	return DatabaseDeploymentDefinitionWithOptions()
}

func DatabaseDeploymentDefinitionWithOptions(opts ...func(*CRDOptions)) Definition {
	return Definition{
		Version: DatabaseDeploymentVersion,
		CRD:     DatabaseDeploymentWithOptions(opts...),
	}
}

var databaseDeploymentCRD apiextensions.CustomResourceDefinition
var databaseDeploymentCRDSchemas crdSchemas

//go:embed database-deployment.yaml
var databaseDeployment []byte

//go:embed database-deployment.schema.generated.yaml
var databaseDeploymentSchemaRaw []byte
