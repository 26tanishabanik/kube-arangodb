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

package inspector

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/arangodb/kube-arangodb/pkg/deployment/patch"
	"github.com/arangodb/kube-arangodb/pkg/logging"
	"github.com/arangodb/kube-arangodb/pkg/util"
	inspectorInterface "github.com/arangodb/kube-arangodb/pkg/util/k8sutil/inspector"
	"github.com/arangodb/kube-arangodb/pkg/util/k8sutil/inspector/constants"
	v1 "github.com/arangodb/kube-arangodb/pkg/util/k8sutil/inspector/endpoints/v1"
	"github.com/arangodb/kube-arangodb/pkg/util/tests/tlogging"
)

func Test_Endpoints_ModV1_Metrics(t *testing.T) {
	i := newInspector()

	client := func(i inspectorInterface.Inspector) v1.ModInterface {
		return i.EndpointsModInterface().V1()
	}

	kind := constants.EndpointsKind

	object := func(name string, args ...interface{}) *core.Endpoints {
		return &core.Endpoints{
			ObjectMeta: meta.ObjectMeta{
				Namespace: "test",
				Name:      fmt.Sprintf(name, args...),
			},
		}
	}

	withMetricsTest(t, "Create", logging.Info, func(t *testing.T) {
		_, err := client(i).Create(context.Background(), object("test-1"), meta.CreateOptions{})

		require.NoError(t, err)
	}, func(t *testing.T, metrics metricsDefinition, log tlogging.LogScanner) {
		_, got := log.Get(100 * time.Millisecond)
		require.True(t, got)

		require.Len(t, metrics, 1)
		require.Contains(t, metrics, kind)
		require.Len(t, metrics[kind], 1)
		require.Contains(t, metrics[kind], constants.CreateOperation)
		require.Equal(t, 1, metrics[kind][constants.CreateOperation])
	})

	withMetricsTest(t, "Create - AboveLogLevel", logging.Warn, func(t *testing.T) {
		_, err := client(i).Create(context.Background(), object("test-2"), meta.CreateOptions{})

		require.NoError(t, err)
	}, func(t *testing.T, metrics metricsDefinition, log tlogging.LogScanner) {
		_, got := log.Get(100 * time.Millisecond)
		require.False(t, got)

		require.Len(t, metrics, 1)
		require.Contains(t, metrics, kind)
		require.Len(t, metrics[kind], 1)
		require.Contains(t, metrics[kind], constants.CreateOperation)
		require.Equal(t, 1, metrics[kind][constants.CreateOperation])
	})

	withMetricsTest(t, "Update", logging.Debug, func(t *testing.T) {
		_, err := client(i).Update(context.Background(), object("test-1"), meta.UpdateOptions{})

		require.NoError(t, err)
	}, func(t *testing.T, metrics metricsDefinition, log tlogging.LogScanner) {
		_, got := log.Get(100 * time.Millisecond)
		require.True(t, got)

		require.Len(t, metrics, 1)
		require.Contains(t, metrics, kind)
		require.Len(t, metrics[kind], 1)
		require.Contains(t, metrics[kind], constants.UpdateOperation)
		require.Equal(t, 1, metrics[kind][constants.UpdateOperation])
	})

	withMetricsTest(t, "Update - AboveLogLevel", logging.Info, func(t *testing.T) {
		_, err := client(i).Update(context.Background(), object("test-2"), meta.UpdateOptions{})

		require.NoError(t, err)
	}, func(t *testing.T, metrics metricsDefinition, log tlogging.LogScanner) {
		_, got := log.Get(100 * time.Millisecond)
		require.False(t, got)

		require.Len(t, metrics, 1)
		require.Contains(t, metrics, kind)
		require.Len(t, metrics[kind], 1)
		require.Contains(t, metrics[kind], constants.UpdateOperation)
		require.Equal(t, 1, metrics[kind][constants.UpdateOperation])
	})

	d, err := patch.NewPatch(patch.ItemReplace(patch.NewPath("metadata", "labels"), map[string]string{
		"a": "b",
	})).Marshal()
	require.NoError(t, err)

	withMetricsTest(t, "Patch", logging.Debug, func(t *testing.T) {
		_, err := client(i).Patch(context.Background(), "test-1", types.JSONPatchType, d, meta.PatchOptions{})

		require.NoError(t, err)
	}, func(t *testing.T, metrics metricsDefinition, log tlogging.LogScanner) {
		_, got := log.Get(100 * time.Millisecond)
		require.True(t, got)

		require.Len(t, metrics, 1)
		require.Contains(t, metrics, kind)
		require.Len(t, metrics[kind], 1)
		require.Contains(t, metrics[kind], constants.PatchOperation)
		require.Equal(t, 1, metrics[kind][constants.PatchOperation])
	})

	withMetricsTest(t, "Patch - AboveLogLevel", logging.Info, func(t *testing.T) {
		_, err := client(i).Patch(context.Background(), "test-2", types.JSONPatchType, d, meta.PatchOptions{})

		require.NoError(t, err)
	}, func(t *testing.T, metrics metricsDefinition, log tlogging.LogScanner) {
		_, got := log.Get(100 * time.Millisecond)
		require.False(t, got)

		require.Len(t, metrics, 1)
		require.Contains(t, metrics, kind)
		require.Len(t, metrics[kind], 1)
		require.Contains(t, metrics[kind], constants.PatchOperation)
		require.Equal(t, 1, metrics[kind][constants.PatchOperation])
	})

	withMetricsTest(t, "Delete", logging.Info, func(t *testing.T) {
		err := client(i).Delete(context.Background(), "test-1", meta.DeleteOptions{})

		require.NoError(t, err)
	}, func(t *testing.T, metrics metricsDefinition, log tlogging.LogScanner) {
		_, got := log.Get(100 * time.Millisecond)
		require.True(t, got)

		require.Len(t, metrics, 1)
		require.Contains(t, metrics, kind)
		require.Len(t, metrics[kind], 1)
		require.Contains(t, metrics[kind], constants.DeleteOperation)
		require.Equal(t, 1, metrics[kind][constants.DeleteOperation])
	})

	withMetricsTest(t, "Delete - AboveLogLevel", logging.Warn, func(t *testing.T) {
		err := client(i).Delete(context.Background(), "test-2", meta.DeleteOptions{})

		require.NoError(t, err)
	}, func(t *testing.T, metrics metricsDefinition, log tlogging.LogScanner) {
		_, got := log.Get(100 * time.Millisecond)
		require.False(t, got)

		require.Len(t, metrics, 1)
		require.Contains(t, metrics, kind)
		require.Len(t, metrics[kind], 1)
		require.Contains(t, metrics[kind], constants.DeleteOperation)
		require.Equal(t, 1, metrics[kind][constants.DeleteOperation])
	})

	for j := 0; j < 10; j++ {
		withMetricsTest(t, fmt.Sprintf("Create for Delete %d", j), logging.Info, func(t *testing.T) {
			_, err := client(i).Create(context.Background(), object("test-local-%d", j), meta.CreateOptions{})

			require.NoError(t, err)
		}, func(t *testing.T, metrics metricsDefinition, log tlogging.LogScanner) {
			_, got := log.Get(100 * time.Millisecond)
			require.True(t, got)

			require.Len(t, metrics, 1)
			require.Contains(t, metrics, kind)
			require.Len(t, metrics[kind], 1)
			require.Contains(t, metrics[kind], constants.CreateOperation)
			require.Equal(t, 1, metrics[kind][constants.CreateOperation])
		})
	}

	withMetricsTest(t, "Delete - With non-0 DeletionGracePeriod", logging.Info, func(t *testing.T) {
		err := client(i).Delete(context.Background(), "test-local-1", meta.DeleteOptions{
			GracePeriodSeconds: util.NewInt64(1),
		})

		require.NoError(t, err)
	}, func(t *testing.T, metrics metricsDefinition, log tlogging.LogScanner) {
		_, got := log.Get(100 * time.Millisecond)
		require.True(t, got)

		require.Len(t, metrics, 1)
		require.Contains(t, metrics, kind)
		require.Len(t, metrics[kind], 1)
		require.Contains(t, metrics[kind], constants.DeleteOperation)
		require.Equal(t, 1, metrics[kind][constants.DeleteOperation])
	})

	withMetricsTest(t, "Delete - With non-0 DeletionGracePeriod - Above Log Level", logging.Warn, func(t *testing.T) {
		err := client(i).Delete(context.Background(), "test-local-2", meta.DeleteOptions{
			GracePeriodSeconds: util.NewInt64(1),
		})

		require.NoError(t, err)
	}, func(t *testing.T, metrics metricsDefinition, log tlogging.LogScanner) {
		_, got := log.Get(100 * time.Millisecond)
		require.False(t, got)

		require.Len(t, metrics, 1)
		require.Contains(t, metrics, kind)
		require.Len(t, metrics[kind], 1)
		require.Contains(t, metrics[kind], constants.DeleteOperation)
		require.Equal(t, 1, metrics[kind][constants.DeleteOperation])
	})

	withMetricsTest(t, "Delete - With 0 DeletionGracePeriod", logging.Warn, func(t *testing.T) {
		err := client(i).Delete(context.Background(), "test-local-3", meta.DeleteOptions{
			GracePeriodSeconds: util.NewInt64(0),
		})

		require.NoError(t, err)
	}, func(t *testing.T, metrics metricsDefinition, log tlogging.LogScanner) {
		_, got := log.Get(100 * time.Millisecond)
		require.True(t, got)

		require.Len(t, metrics, 1)
		require.Contains(t, metrics, kind)
		require.Len(t, metrics[kind], 1)
		require.Contains(t, metrics[kind], constants.ForceDeleteOperation)
		require.Equal(t, 1, metrics[kind][constants.ForceDeleteOperation])
	})

	withMetricsTest(t, "Delete - With 0 DeletionGracePeriod - Above Log Level", logging.Error, func(t *testing.T) {
		err := client(i).Delete(context.Background(), "test-local-4", meta.DeleteOptions{
			GracePeriodSeconds: util.NewInt64(0),
		})

		require.NoError(t, err)
	}, func(t *testing.T, metrics metricsDefinition, log tlogging.LogScanner) {
		_, got := log.Get(100 * time.Millisecond)
		require.False(t, got)

		require.Len(t, metrics, 1)
		require.Contains(t, metrics, kind)
		require.Len(t, metrics[kind], 1)
		require.Contains(t, metrics[kind], constants.ForceDeleteOperation)
		require.Equal(t, 1, metrics[kind][constants.ForceDeleteOperation])
	})
}
