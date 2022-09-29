//
// DISCLAIMER
//
// Copyright 2016-2022 ArangoDB GmbH, Cologne, Germany
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

package replication

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"time"

	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/arangodb/arangosync-client/client"
	"github.com/arangodb/arangosync-client/client/synccheck"
	"github.com/arangodb/go-driver"

	api "github.com/arangodb/kube-arangodb/pkg/apis/replication/v1"
	"github.com/arangodb/kube-arangodb/pkg/deployment/features"
	"github.com/arangodb/kube-arangodb/pkg/util/errors"
)

// IsBeingDeleted returns true is object is marked for deletion.
// It can be removed in two modes:
// - set condition on the resource (recommended).
// - remove resource (not recommended).
func (dr *DeploymentReplication) IsBeingDeleted() (*meta.Time, bool) {
	if timestamp := dr.apiObject.GetDeletionTimestamp(); timestamp != nil {
		// In this mode of deletion we can not change resource anymore.
		return timestamp, true
	}

	if cond, ok := dr.status.Conditions.Get(api.ConditionTypeConfigured); ok {
		if cond.Reason == api.ConditionTypeConfiguredInvalid {
			// In this mode of deletion we can change resource, so we will see actual progress in the status.
			return &cond.LastTransitionTime, true
		}
	}

	return nil, false
}

// inspectDeploymentReplication inspects the entire deployment replication
// and configures the replication when needed.
// This function should be called when:
// - the deployment replication has changed
// - any of the underlying resources has changed
// - once in a while
// Returns the delay until this function should be called again.
func (dr *DeploymentReplication) inspectDeploymentReplication(ctx context.Context, lastInterval time.Duration) (time.Duration, error) {
	spec := dr.apiObject.Spec

	// Add finalizers
	if err := dr.addFinalizers(); err != nil {
		return 0, errors.WithMessage(err, "Failed to add finalizers")
	}

	// Is the deployment in failed state, if so, give up.
	if dr.status.Phase.IsFailed() {
		dr.log.Debug("Deployment replication is in Failed state.")
		return lastInterval, nil
	}

	if timestamp, ok := dr.IsBeingDeleted(); ok {
		// Resource is being deleted.
		retrySoon, err := dr.runFinalizers(ctx, dr.apiObject)
		if err != nil || retrySoon {
			timeout := CancellationTimeout + AbortTimeout
			if isTimeExceeded(timestamp, timeout) {
				// Cancellation and abort timeout exceeded, so it must go into failed state.
				dr.reportDeploymentReplicationErr(err, fmt.Sprintf("Failed to cancel synchronization in %s", timeout.String()))
			}
		}

		if err != nil {
			return cancellationInterval, errors.WithMessage(err, "Failed to run finalizers")
		}

		return cancellationInterval, nil
	}

	// Inspect configuration status
	destClient, err := dr.createSyncMasterClient(spec.Destination)
	if err != nil {
		return 0, errors.WithMessage(err, "Failed to create destination sync master client")
		//dr.status.Conditions.Update(api.ConditionTypeConfigured, true, api.ConditionTypeConfiguredActive,
		//	"Destination syncmaster is configured correctly and active")
		//dr.status.IncomingSynchronization = dr.inspectIncomingSynchronizationStatus(ctx, destClient,
		//	driver.Version(destArangosyncVersion.Version), destStatus.Shards)
		//if dr.status.Conditions.Update(api.ConditionTypeConfigured, false, api.ConditionTypeConfiguredInvalid,
		//	"Destination syncmaster is configured for different source") {
	}

	configureNeeded, updateStatusNeeded := false, false
	// Inspect destination DC.
	if ok, err := dr.checkDestinationDC(ctx, destClient); err != nil {
		dr.log.Err(err).Warn("Failed to reconcile destination data center")
		// Don't return here.
	} else if ok {
		if dr.isConfigurationNeeded() {
			configureNeeded = true
		}
		//if dr.status.Conditions.Update(api.ConditionTypeConfigured, false, api.ConditionTypeConfiguredInactive,
		//	"Destination syncmaster is configured correctly but in-active") {
		updateStatusNeeded = true
	}

	// Inspect source DC.
	if ok, err := dr.setSourceEndpoints(ctx, destClient); err != nil {
		dr.log.Err(err).Warn("Failed to set source endpoints")
		// Don't return here.
	} else if ok {
		updateStatusNeeded = true
	}

	if updateStatusNeeded {
		if err := dr.updateCRStatus(); err != nil {
			return time.Second * 3, errors.WithMessage(err, "Failed to update status")
		}
	}

	if configureNeeded {
		if err := dr.configureSynchronization(ctx, destClient); err != nil {
			return 0, errors.WithMessagef(err, "")
		}

		return time.Second * 10, nil
	}

	return nextInterval, nil
}

// isIncomingEndpoint returns true when given sync status's endpoint
// intersects with the given endpoint spec.
func (dr *DeploymentReplication) isIncomingEndpoint(status client.SyncInfo, epSpec api.EndpointSpec) (bool, error) {
	ep, err := dr.createArangoSyncEndpoint(epSpec)
	if err != nil {
		return false, errors.WithStack(err)
	}
	return !status.Source.Intersection(ep).IsEmpty(), nil
}

// hasOutgoingEndpoint returns true when given sync status has an outgoing
// item that intersects with the given endpoint spec.
// Returns: outgoing-ID, outgoing-found, error
func (dr *DeploymentReplication) hasOutgoingEndpoint(status client.SyncInfo, epSpec api.EndpointSpec, reportedEndpoint client.Endpoint) (string, bool, error) {
	ep, err := dr.createArangoSyncEndpoint(epSpec)
	if err != nil {
		return "", false, errors.WithStack(err)
	}
	ep = ep.Merge(reportedEndpoint...)
	for _, o := range status.Outgoing {
		if !o.Endpoint.Intersection(ep).IsEmpty() {
			return o.ID, true, nil
		}
	}
	return "", false, nil
}

// inspectIncomingSynchronizationStatus returns the synchronization status for the incoming sync
func (dr *DeploymentReplication) inspectIncomingSynchronizationStatus(ctx context.Context, syncClient client.API, arangosyncVersion driver.Version, localShards []client.ShardSyncInfo) api.SynchronizationStatus {
	dataCentersResp, err := syncClient.Master().GetDataCentersInfo(ctx)
	if err != nil {
		errMsg := "Failed to fetch data-centers info"
		dr.log.Err(err).Warn(errMsg)
		return api.SynchronizationStatus{
			Error: fmt.Sprintf("%s: %s", errMsg, err.Error()),
		}
	}

	ch := synccheck.NewSynchronizationChecker(syncClient, time.Minute)
	incomingSyncStatus, err := ch.CheckSync(ctx, &dataCentersResp, localShards)
	if err != nil {
		errMsg := "Failed to check synchronization status"
		dr.log.Err(err).Warn(errMsg)
		return api.SynchronizationStatus{
			Error: fmt.Sprintf("%s: %s", errMsg, err.Error()),
		}
	}
	return dr.createSynchronizationStatus(arangosyncVersion, incomingSyncStatus)
}

// createSynchronizationStatus returns aggregated info about DCSyncStatus
func (dr *DeploymentReplication) createSynchronizationStatus(arangosyncVersion driver.Version, dcSyncStatus *synccheck.DCSyncStatus) api.SynchronizationStatus {
	dbs := make(map[string]api.DatabaseSynchronizationStatus, len(dcSyncStatus.Databases))
	i := 0
	for dbName, dbSyncStatus := range dcSyncStatus.Databases {
		i++
		db := dbName
		if features.SensitiveInformationProtection().Enabled() {
			// internal IDs are not available in older versions
			if arangosyncVersion.CompareTo("2.12.0") >= 0 {
				db = dbSyncStatus.ID
			} else {
				db = fmt.Sprintf("<PROTECTED_INFO_%d>", i)
			}
		}
		dbs[db] = dr.createDatabaseSynchronizationStatus(dbSyncStatus)
	}
	return api.SynchronizationStatus{
		AllInSync: dcSyncStatus.AllInSync(),
		Databases: dbs,
		Error:     "",
	}
}

// createDatabaseSynchronizationStatus returns sync status for DB
func (dr *DeploymentReplication) createDatabaseSynchronizationStatus(dbSyncStatus synccheck.DatabaseSyncStatus) api.DatabaseSynchronizationStatus {
	// use limit for errors because the resulting status object should not be too big
	const maxReportedIncomingSyncErrors = 20

	var errs []api.DatabaseSynchronizationError
	var shardsTotal, shardsInSync int
	for colName, colSyncStatus := range dbSyncStatus.Collections {
		col := colName
		if features.SensitiveInformationProtection().Enabled() {
			col = colSyncStatus.ID
		}
		if colSyncStatus.Error != "" && len(errs) < maxReportedIncomingSyncErrors {
			errs = append(errs, api.DatabaseSynchronizationError{
				Collection: col,
				Shard:      "",
				Message:    colSyncStatus.Error,
			})
		}

		shardsTotal += len(colSyncStatus.Shards)
		for shardIndex, shardSyncStatus := range colSyncStatus.Shards {
			if shardSyncStatus.InSync {
				shardsInSync++
			} else if len(errs) < maxReportedIncomingSyncErrors {
				errs = append(errs, api.DatabaseSynchronizationError{
					Collection: col,
					Shard:      strconv.Itoa(shardIndex),
					Message:    shardSyncStatus.Message,
				})
			}
		}
	}

	return api.DatabaseSynchronizationStatus{
		ShardsTotal:  shardsTotal,
		ShardsInSync: shardsInSync,
		Errors:       errs,
	}
}

// createEndpointStatus creates an api EndpointStatus from the given sync status.
func createEndpointStatus(status client.SyncInfo, outgoingID string) api.EndpointStatus {
	if outgoingID == "" {
		// Incoming DC.
		return createEndpointStatusFromShards(status.Shards)
	}

	for _, o := range status.Outgoing {
		if o.ID == outgoingID {
			// Outgoing DC.
			return createEndpointStatusFromShards(o.Shards)
		}
	}

	return api.EndpointStatus{}
}

// createEndpointStatusFromShards creates an api EndpointStatus from the given list of shard statuses.
func createEndpointStatusFromShards(shards []client.ShardSyncInfo) api.EndpointStatus {
	result := api.EndpointStatus{}

	getDatabase := func(name string) *api.DatabaseStatus {
		for i, d := range result.Databases {
			if d.Name == name {
				return &result.Databases[i]
			}
		}
		// Not found, add it
		result.Databases = append(result.Databases, api.DatabaseStatus{Name: name})
		return &result.Databases[len(result.Databases)-1]
	}

	getCollection := func(db *api.DatabaseStatus, name string) *api.CollectionStatus {
		for i, c := range db.Collections {
			if c.Name == name {
				return &db.Collections[i]
			}
		}
		// Not found, add it
		db.Collections = append(db.Collections, api.CollectionStatus{Name: name})
		return &db.Collections[len(db.Collections)-1]
	}

	// Sort shard by index
	sort.Slice(shards, func(i, j int) bool {
		return shards[i].ShardIndex < shards[j].ShardIndex
	})
	for _, s := range shards {
		db := getDatabase(s.Database)
		col := getCollection(db, s.Collection)

		// Add "missing" shards if needed
		for len(col.Shards) < s.ShardIndex {
			col.Shards = append(col.Shards, api.ShardStatus{Status: ""})
		}

		// Add current shard
		col.Shards = append(col.Shards, api.ShardStatus{Status: string(s.Status)})
	}

	// Sort result
	sort.Slice(result.Databases, func(i, j int) bool { return result.Databases[i].Name < result.Databases[j].Name })
	for i, db := range result.Databases {
		sort.Slice(db.Collections, func(i, j int) bool { return db.Collections[i].Name < db.Collections[j].Name })
		result.Databases[i] = db
	}
	return result
}

//
func (dr *DeploymentReplication) checkDestinationDC(ctx context.Context, destClient client.API) (bool, error) {
	destStatus, err := destClient.Master().Status(ctx)
	if err != nil {
		return false, errors.WithMessage(err, "Failed to fetch status from destination sync master")
	}

	// TODO test what if it is cancelling?
	if !destStatus.Status.IsActive() {
		// Destination has correct source, but is inactive.
		changed := dr.status.Conditions.Update(api.ConditionTypeConfigured, false, api.ConditionTypeConfiguredInactive,
			"Destination sync master is configured correctly but in-active")

		return changed, nil
	}

	spec := dr.apiObject.Spec
	ep, err := dr.createArangoSyncEndpoint(spec.Source)
	if err != nil {
		return false, errors.WithMessage(err, "Failed to create source sync master endpoint")
	}

	// Check whether destination source endpoints and spec endpoints contain each other.
	// If intersection is not empty it means that endpoints match.
	isIncomingEndpoint := !destStatus.Source.Intersection(ep).IsEmpty()
	if isIncomingEndpoint {
		destArangosyncVersion, err := destClient.Version(ctx)
		if err != nil {
			return false, errors.WithMessage(err, "Failed to get destination arangosync version")
		}

		// Destination is correctly configured
		dr.status.Conditions.Update(api.ConditionTypeConfigured, true, api.ConditionTypeConfiguredActive,
			"Destination sync master is configured correctly and active")
		dr.status.Destination = createEndpointStatus(destStatus, "")
		dr.status.IncomingSynchronization = dr.inspectIncomingSynchronizationStatus(ctx, destClient,
			driver.Version(destArangosyncVersion.Version), destStatus.Shards)
	} else {
		// Sync is active, but from different source.
		dr.log.
			Strs("spec-endpoints", ep...).
			Strs("sync-endpoints", destStatus.Source...).
			Warn("Destination sync master is configured for different source")

		if !dr.status.Conditions.Update(api.ConditionTypeConfigured, false, api.ConditionTypeConfiguredInvalid,
			"Destination sync master is configured for different source") {
			return false, nil
		}
	}

	return true, nil
}

// setSourceEndpoints sets source endpoints in the ArangoDeploymentReplication status.
func (dr *DeploymentReplication) setSourceEndpoints(ctx context.Context, destClient client.API) (bool, error) {
	spec := dr.apiObject.Spec

	sourceClient, err := dr.createSyncMasterClient(spec.Source)
	if err != nil {
		return false, errors.WithMessage(err, "Failed to create source sync master client")
	}

	sourceStatus, err := sourceClient.Master().Status(ctx)
	if err != nil {
		return false, errors.WithMessage(err, "Failed to fetch status from source sync master")
	}

	destEndpoint, err := destClient.Master().GetEndpoints(ctx)
	if err != nil {
		return false, errors.WithMessage(err, "Failed to fetch endpoints from destination sync master")
	}

	outgoingID, hasOutgoingEndpoint, err := dr.hasOutgoingEndpoint(sourceStatus, spec.Destination, destEndpoint)
	if err != nil {
		return false, errors.WithMessage(err, "Failed to check has-outgoing-endpoint")
	} else if !hasOutgoingEndpoint {
		return false, errors.New("Destination not yet known in source sync masters")
	}

	dr.status.Source = createEndpointStatus(sourceStatus, outgoingID)
	// TODO check if source is the some and if it is then return false, nil
	return true, nil
}

// isConfigurationNeeded returns true when configuration.
func (dr *DeploymentReplication) isConfigurationNeeded() bool {
	cond, ok := dr.status.Conditions.Get(api.ConditionTypeConfigured)
	if !ok {
		return false
	}

	if cond.Reason == api.ConditionTypeConfiguredInactive && cond.Status == core.ConditionFalse {
		return true
	}

	return false
}

// configureSynchronization turns on synchronization between two deployments.
func (dr *DeploymentReplication) configureSynchronization(ctx context.Context, destClient client.API) error {
	spec := dr.apiObject.Spec
	source, err := dr.createArangoSyncEndpoint(spec.Source)
	if err != nil {
		return errors.WithMessage(err, "Failed to create source sync master endpoint")
	}

	auth, err := dr.createArangoSyncTLSAuthentication(spec)
	if err != nil {
		return errors.WithMessage(err, "Failed to configure synchronization authentication")
	}

	req := client.SynchronizationRequest{
		Source:         source,
		Authentication: auth,
	}
	dr.log.Info("Configuring synchronization")
	if err := destClient.Master().Synchronize(ctx, req); err != nil {
		return errors.WithMessage(err, "Failed to configure synchronization")
	}

	dr.log.Info("Configured synchronization")
	return nil
}
