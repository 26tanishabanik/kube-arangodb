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
	"testing"

	"github.com/arangodb/kube-arangodb/pkg/logging"
	"github.com/arangodb/kube-arangodb/pkg/util/k8sutil/inspector"
	"github.com/arangodb/kube-arangodb/pkg/util/k8sutil/inspector/throttle"
	"github.com/arangodb/kube-arangodb/pkg/util/kclient"
	"github.com/arangodb/kube-arangodb/pkg/util/tests/tlogging"
)

func newInspector() inspector.Inspector {
	c := kclient.NewFakeClient()

	return NewInspector(throttle.NewAlwaysThrottleComponents(), c, "test", "test")
}

func withMetricsTest(t *testing.T, name string, level logging.Level, invoke func(t *testing.T), validate func(t *testing.T, metrics metricsDefinition, log tlogging.LogScanner)) {
	tlogging.WithLogScanner(t, name, func(t *testing.T, s tlogging.LogScanner) {
		c := clientLogger
		defer func() {
			clientLogger = c
			metricsObject = metricsDefinition{}
		}()
		clientLogger = s.Factory().RegisterAndGetLogger("kubernetes-client", level)

		t.Run("Invoke", func(t *testing.T) {
			invoke(t)
		})

		t.Run("Validate", func(t *testing.T) {
			validate(t, metricsObject, s)
		})
	})
}
