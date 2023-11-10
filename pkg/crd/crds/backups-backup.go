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
	BackupsBackupVersion = driver.Version("1.0.1")
)

func init() {
	mustLoadCRD(backupsBackup, backupsBackupSchemaRaw, &backupsBackupCRD, &backupsBackupCRDWithSchema)
}

func BackupsBackup(opts ...func(*CRDOptions)) *apiextensions.CustomResourceDefinition {
	return getCRD(backupsBackupCRD, backupsBackupCRDWithSchema, opts...)
}

func BackupsBackupDefinition() Definition {
	return Definition{
		Version:       BackupsBackupVersion,
		CRD:           backupsBackupCRD.DeepCopy(),
		CRDWithSchema: backupsBackupCRDWithSchema.DeepCopy(),
	}
}

var backupsBackupCRD apiextensions.CustomResourceDefinition
var backupsBackupCRDWithSchema apiextensions.CustomResourceDefinition

//go:embed backups-backup.yaml
var backupsBackup []byte

//go:embed backups-backup.schema.generated.yaml
var backupsBackupSchemaRaw []byte
