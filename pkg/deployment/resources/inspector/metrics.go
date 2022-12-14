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
	"reflect"
	"sync"

	meta "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/arangodb/kube-arangodb/pkg/generated/metric_descriptions"
	"github.com/arangodb/kube-arangodb/pkg/metrics/collector"
	"github.com/arangodb/kube-arangodb/pkg/util/k8sutil/inspector/constants"
	"github.com/arangodb/kube-arangodb/pkg/util/metrics"
)

func init() {
	collector.GetCollector().RegisterMetric(collectorObject{})
}

type collectorObject struct {
}

func (c collectorObject) CollectMetrics(in metrics.PushMetric) {
	metricsLock.Lock()
	defer metricsLock.Unlock()

	for k, v := range metricsObject {
		for o, c := range v {
			in.Push(metric_descriptions.ArangodbOperatorKubernetesClientRequestsCounter(float64(c), k, string(o)))
		}
	}
}

type metricsDefinition map[string]map[constants.Operation]int

var (
	metricsObject metricsDefinition
	metricsLock   sync.Mutex
)

func logDeleteOperation(kind string, namespace, name string, opts meta.DeleteOptions) {
	op := constants.DeleteOperation

	if o := opts.GracePeriodSeconds; o != nil && *o == 0 {
		op = constants.ForceDeleteOperation
	}

	logOperation(kind, op, namespace, name)
}

func logPatchOperation(kind string, namespace, name string) {
	logOperation(kind, constants.PatchOperation, namespace, name)
}

func logCreateOperation(kind string, object meta.Object) {
	logObjectOperation(kind, constants.CreateOperation, object)
}

func logUpdateStatusOperation(kind string, object meta.Object) {
	logObjectOperation(kind, constants.UpdateStatusOperation, object)
}

func logUpdateOperation(kind string, object meta.Object) {
	logObjectOperation(kind, constants.UpdateOperation, object)
}

func logObjectOperation(kind string, operation constants.Operation, object meta.Object) {
	if reflect.ValueOf(object).IsZero() {
		return
	}

	logOperation(kind, operation, object.GetNamespace(), object.GetName())
}

func logOperation(kind string, operation constants.Operation, namespace, name string) {
	metricsLock.Lock()
	defer metricsLock.Unlock()

	if metricsObject == nil {
		metricsObject = metricsDefinition{}
	}

	q := metricsObject[kind]
	if q == nil {
		q = map[constants.Operation]int{}
	}

	q[operation] += 1

	metricsObject[kind] = q

	l := clientLogger.
		Str("namespace", namespace).
		Str("name", name).
		Str("kind", kind).
		Str("operation", string(operation))

	c := l.Debug

	switch operation {
	case constants.CreateOperation, constants.DeleteOperation:
		c = l.Info
	case constants.ForceDeleteOperation:
		c = l.Warn
	}

	c("Kubernetes Request send")
}
