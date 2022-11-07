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
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/arangodb/kube-arangodb/pkg/generated/metric_descriptions"
	"github.com/arangodb/kube-arangodb/pkg/util/k8sutil"
	"github.com/arangodb/kube-arangodb/pkg/util/k8sutil/inspector/definitions"
	"github.com/arangodb/kube-arangodb/pkg/util/metrics"
)

func init() {
	prometheus.MustRegister(metricsDefinitions)
}

type metricsMapObject struct {
	lock sync.Mutex

	metrics map[definitions.Verb]map[definitions.Component]metricsObject
}

func (m *metricsMapObject) Describe(descs chan<- *prometheus.Desc) {
}

func (m *metricsMapObject) Collect(c chan<- prometheus.Metric) {
	m.lock.Lock()
	defer m.lock.Unlock()

	p := metrics.NewPushMetric(c)

	for component, data := range m.metrics[definitions.Create] {
		p.Push(metric_descriptions.ArangodbOperatorKubernetesClientCreateRequestsCounter(float64(data.executions), string(component)),
			metric_descriptions.ArangodbOperatorKubernetesClientCreateRequestsCounter(float64(data.errors), string(component)))
	}
}

type metricsObject struct {
	executions int
	errors     int
}

var (
	metricsDefinitions = &metricsMapObject{}
)

func registerDeleteCall(component definitions.Component, cname string, opts meta.DeleteOptions, err error, namespace, name string) {
	if k8sutil.IsForceDeletion(opts) {
		metricsDefinitions.registerCall(component, cname, definitions.ForceDelete, err, namespace, name)
	} else {
		metricsDefinitions.registerCall(component, cname, definitions.Delete, err, namespace, name)
	}
}

func registerCreateCall(component definitions.Component, cname string, _ meta.CreateOptions, err error, namespace string, obj meta.Object) {
	if obj != nil {
		metricsDefinitions.registerCall(component, cname, definitions.Create, err, namespace, obj.GetName())
	}
}

func registerPatchCall(component definitions.Component, cname string, _ meta.PatchOptions, err error, namespace string, obj meta.Object) {
	if obj != nil {
		metricsDefinitions.registerCall(component, cname, definitions.Patch, err, namespace, obj.GetName())
	}
}

func registerUpdateCall(component definitions.Component, cname string, _ meta.UpdateOptions, err error, namespace string, obj meta.Object) {
	if obj != nil {
		metricsDefinitions.registerCall(component, cname, definitions.Update, err, namespace, obj.GetName())
	}
}

func (m *metricsMapObject) registerCall(component definitions.Component, cname string, verb definitions.Verb, err error, namespace, name string) {
	m.lock.Lock()
	defer m.lock.Unlock()

	if m.metrics == nil {
		m.metrics = map[definitions.Verb]map[definitions.Component]metricsObject{}
	}

	a, ok := m.metrics[verb]
	if !ok {
		a = map[definitions.Component]metricsObject{}
	}

	b, ok := a[component]
	if !ok {
		b = metricsObject{}
	}

	b.executions++

	l := clientLogger.Wrap(func(in *zerolog.Event) *zerolog.Event {
		return in.Str("component", string(component)).
			Str("verb", string(verb)).
			Str("namespace", namespace).
			Str("name", name).
			Str("cname", cname).Err(err)
	})

	if err == nil {
		f := l.Trace

		switch verb {
		case definitions.ForceDelete:
			f = l.Warn
		case definitions.Delete:
			f = l.Info
		case definitions.Create:
			f = l.Debug
		}

		f("Kubernetes Client request")
	} else {
		b.errors++
		l.Warn("Kubernetes Client request failed")
	}

	a[component] = b
	m.metrics[verb] = a
}
