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
	DatabaseTaskVersion = driver.Version("1.0.1")
)

func init() {
	mustLoadCRD(databaseTask, databaseTaskSchemaRaw, &databaseTaskCRD, &databaseTaskCRDWithSchema)
}

func DatabaseTask(opts ...func(*CRDOptions)) *apiextensions.CustomResourceDefinition {
	return getCRD(databaseTaskCRD, databaseTaskCRDWithSchema, opts...)
}

func DatabaseTaskDefinition() Definition {
	return Definition{
		Version:       DatabaseTaskVersion,
		CRD:           databaseTaskCRD.DeepCopy(),
		CRDWithSchema: databaseTaskCRDWithSchema.DeepCopy(),
	}
}

var databaseTaskCRD apiextensions.CustomResourceDefinition
var databaseTaskCRDWithSchema apiextensions.CustomResourceDefinition

//go:embed database-task.yaml
var databaseTask []byte

//go:embed database-task.schema.generated.yaml
var databaseTaskSchemaRaw []byte
