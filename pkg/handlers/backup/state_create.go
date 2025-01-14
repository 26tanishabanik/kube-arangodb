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

package backup

import (
	"github.com/arangodb/go-driver"
	"github.com/arangodb/go-driver/util/connection/wrappers/async"

	backupApi "github.com/arangodb/kube-arangodb/pkg/apis/backup/v1"
	"github.com/arangodb/kube-arangodb/pkg/deployment/features"
)

func stateCreateHandler(h *handler, backup *backupApi.ArangoBackup) (*backupApi.ArangoBackupStatus, error) {
	deployment, err := h.getArangoDeploymentObject(backup)
	if err != nil {
		return nil, err
	}

	client, err := h.arangoClientFactory(deployment, backup)
	if err != nil {
		return nil, newTemporaryError(err)
	}

	if features.AsyncBackupCreation().Enabled() {
		return asyncBackup(client, backup)
	} else {
		return syncBackup(client, backup)
	}
}

func syncBackup(client ArangoBackupClient, backup *backupApi.ArangoBackup) (*backupApi.ArangoBackupStatus, error) {
	response, err := client.Create()
	if err != nil {
		return wrapUpdateStatus(backup,
			updateStatusState(backupApi.ArangoBackupStateCreateError, "Create failed with error: %s", err.Error()),
			cleanStatusJob(),
			updateStatusAvailable(false),
			addBackOff(backup.Spec),
		)
	}

	backupMeta, err := client.Get(response.ID)
	if err != nil {
		if driver.IsNotFoundGeneral(err) {
			return wrapUpdateStatus(backup,
				updateStatusState(backupApi.ArangoBackupStateFailed, "backup is not present after creation"),
				cleanStatusJob(),
			)
		}

		return nil, newFatalError(err)
	}

	return wrapUpdateStatus(backup,
		updateStatusState(backupApi.ArangoBackupStateReady, ""),
		updateStatusAvailable(true),
		updateStatusBackup(backupMeta),
		cleanBackOff(),
	)
}

func asyncBackup(client ArangoBackupClient, backup *backupApi.ArangoBackup) (*backupApi.ArangoBackupStatus, error) {
	_, err := client.CreateAsync("")

	if err == nil {
		return wrapUpdateStatus(backup,
			updateStatusState(backupApi.ArangoBackupStateCreateError, "Start Async backup failed - no jobID response"),
			cleanStatusJob(),
			updateStatusAvailable(false),
			addBackOff(backup.Spec),
		)
	}

	jobID, isAsyncId := async.IsAsyncJobInProgress(err)
	if !isAsyncId {
		return wrapUpdateStatus(backup,
			updateStatusState(backupApi.ArangoBackupStateCreateError, "Create backup failed with error: %s", err.Error()),
			cleanStatusJob(),
			updateStatusAvailable(false),
			addBackOff(backup.Spec),
		)
	}

	return wrapUpdateStatus(backup,
		updateStatusState(backupApi.ArangoBackupStateCreating, ""),
		updateStatusJob(jobID, "0%"),
		updateStatusAvailable(false),
	)
}
