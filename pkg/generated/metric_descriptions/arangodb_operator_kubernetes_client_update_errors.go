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

package metric_descriptions

import "github.com/arangodb/kube-arangodb/pkg/util/metrics"

var (
	arangodbOperatorKubernetesClientUpdateErrors = metrics.NewDescription("arangodb_operator_kubernetes_client_update_errors", "Counter for kubernetes Update request errors", []string{`component`}, nil)
)

func init() {
	registerDescription(arangodbOperatorKubernetesClientUpdateErrors)
}

func ArangodbOperatorKubernetesClientUpdateErrors() metrics.Description {
	return arangodbOperatorKubernetesClientUpdateErrors
}

func ArangodbOperatorKubernetesClientUpdateErrorsCounter(value float64, component string) metrics.Metric {
	return ArangodbOperatorKubernetesClientUpdateErrors().Gauge(value, component)
}
