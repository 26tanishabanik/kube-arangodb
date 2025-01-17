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

package globals

import "time"

const (
	DefaultKubernetesTimeout     = 2 * time.Second
	DefaultArangoDTimeout        = time.Second * 5
	DefaultArangoDAgencyTimeout  = time.Second * 10
	DefaultArangoDCheckTimeout   = time.Second * 2
	DefaultReconciliationTimeout = time.Minute

	BackupDefaultArangoClientTimeout = 30 * time.Second
	BackupUploadArangoClientTimeout  = 300 * time.Second

	// DefaultOutSyncedShardRebuildTimeout
	// timeout after which particular out-synced shard is considered as failed and rebuild is triggered
	DefaultOutSyncedShardRebuildTimeout = time.Minute * 60
	// DefaultOutSyncedShardRebuildRetryTimeout timeout after which rebuild shards retry flow is triggered
	DefaultOutSyncedShardRebuildRetryTimeout = time.Hour * 4

	DefaultKubernetesRequestBatchSize = 256

	DefaultBackupConcurrentUploads = 4

	// Retry

	DefaultOperatorUpdateRetryCount = 25
	DefaultOperatorUpdateRetryDelay = time.Second
)

var globalObj = &globals{
	timeouts: &globalTimeouts{
		requests:                        NewTimeout(DefaultKubernetesTimeout),
		arangod:                         NewTimeout(DefaultArangoDTimeout),
		arangodCheck:                    NewTimeout(DefaultArangoDCheckTimeout),
		reconciliation:                  NewTimeout(DefaultReconciliationTimeout),
		agency:                          NewTimeout(DefaultArangoDAgencyTimeout),
		shardRebuild:                    NewTimeout(DefaultOutSyncedShardRebuildTimeout),
		shardRebuildRetry:               NewTimeout(DefaultOutSyncedShardRebuildRetryTimeout),
		backupArangoClientTimeout:       NewTimeout(BackupDefaultArangoClientTimeout),
		backupArangoClientUploadTimeout: NewTimeout(BackupUploadArangoClientTimeout),
	},
	kubernetes: &globalKubernetes{
		requestBatchSize: NewInt64(DefaultKubernetesRequestBatchSize),
	},
	backup: &globalBackup{
		concurrentUploads: NewInt(DefaultBackupConcurrentUploads),
	},
	retry: &globalRetry{
		operatorUpdateRetryCount: NewInt(DefaultOperatorUpdateRetryCount),
		operatorUpdateRetryDelay: NewTimeout(DefaultOperatorUpdateRetryDelay),
	},
}

func GetGlobals() Globals {
	return globalObj
}

func GetGlobalTimeouts() GlobalTimeouts {
	return globalObj.timeouts
}

type Globals interface {
	Timeouts() GlobalTimeouts
	Kubernetes() GlobalKubernetes
	Backup() GlobalBackup
	Retry() GlobalRetry
}

type globals struct {
	timeouts   *globalTimeouts
	kubernetes *globalKubernetes
	backup     *globalBackup
	retry      *globalRetry
}

func (g globals) Retry() GlobalRetry {
	return g.retry
}

func (g globals) Backup() GlobalBackup {
	return g.backup
}

func (g globals) Kubernetes() GlobalKubernetes {
	return g.kubernetes
}

func (g globals) Timeouts() GlobalTimeouts {
	return g.timeouts
}

type GlobalKubernetes interface {
	RequestBatchSize() Int64
}

type globalKubernetes struct {
	requestBatchSize Int64
}

func (g *globalKubernetes) RequestBatchSize() Int64 {
	return g.requestBatchSize
}

type GlobalBackup interface {
	ConcurrentUploads() Int
}

type globalBackup struct {
	concurrentUploads Int
}

func (g *globalBackup) ConcurrentUploads() Int {
	return g.concurrentUploads
}

type GlobalTimeouts interface {
	Reconciliation() Timeout
	ShardRebuild() Timeout
	ShardRebuildRetry() Timeout

	Kubernetes() Timeout
	ArangoD() Timeout
	ArangoDCheck() Timeout
	Agency() Timeout

	BackupArangoClientTimeout() Timeout
	BackupArangoClientUploadTimeout() Timeout
}

type globalTimeouts struct {
	requests, arangod, reconciliation, arangodCheck, agency, shardRebuild, shardRebuildRetry Timeout
	backupArangoClientTimeout                                                                Timeout
	backupArangoClientUploadTimeout                                                          Timeout
}

func (g *globalTimeouts) Agency() Timeout {
	return g.agency
}

func (g *globalTimeouts) ArangoDCheck() Timeout {
	return g.arangodCheck
}

func (g *globalTimeouts) Reconciliation() Timeout {
	return g.reconciliation
}

func (g *globalTimeouts) ShardRebuild() Timeout {
	return g.shardRebuild
}

func (g *globalTimeouts) ShardRebuildRetry() Timeout {
	return g.shardRebuildRetry
}

func (g *globalTimeouts) ArangoD() Timeout {
	return g.arangod
}

func (g *globalTimeouts) Kubernetes() Timeout {
	return g.requests
}

func (g *globalTimeouts) BackupArangoClientTimeout() Timeout {
	return g.backupArangoClientTimeout
}

func (g *globalTimeouts) BackupArangoClientUploadTimeout() Timeout {
	return g.backupArangoClientUploadTimeout
}

type GlobalRetry interface {
	OperatorUpdateRetryCount() Int
	OperatorUpdateRetryDelay() Timeout
}

type globalRetry struct {
	operatorUpdateRetryCount Int
	operatorUpdateRetryDelay Timeout
}

func (g *globalRetry) OperatorUpdateRetryCount() Int {
	return g.operatorUpdateRetryCount
}

func (g *globalRetry) OperatorUpdateRetryDelay() Timeout {
	return g.operatorUpdateRetryDelay
}
