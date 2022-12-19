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

	core "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/arangodb/kube-arangodb/pkg/util/k8sutil/inspector/constants"
)

type endpointsInspectorAnonymousV1 struct {
	i *inspectorState
}

func (e *endpointsInspectorAnonymousV1) Get(ctx context.Context, name string, opts meta.GetOptions) (meta.Object, error) {
	return e.i.endpoints.v1.Get(ctx, name, opts)
}

func (e *endpointsInspectorAnonymousV1) Create(ctx context.Context, obj meta.Object, opts meta.CreateOptions) (meta.Object, error) {
	if o, ok := obj.(*core.Endpoints); !ok {
		return nil, newInvalidTypeError(constants.EndpointsGKv1())
	} else {
		return e.i.EndpointsModInterface().V1().Create(ctx, o, opts)
	}
}

func (e *endpointsInspectorAnonymousV1) Update(ctx context.Context, obj meta.Object, opts meta.UpdateOptions) (meta.Object, error) {
	if o, ok := obj.(*core.Endpoints); !ok {
		return nil, newInvalidTypeError(constants.EndpointsGKv1())
	} else {
		return e.i.EndpointsModInterface().V1().Update(ctx, o, opts)
	}
}

func (e *endpointsInspectorAnonymousV1) UpdateStatus(ctx context.Context, obj meta.Object, opts meta.UpdateOptions) (meta.Object, error) {
	return nil, newNotImplementedError(constants.EndpointsGKv1())
}

func (e *endpointsInspectorAnonymousV1) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts meta.PatchOptions, subresources ...string) (result meta.Object, err error) {
	return e.i.EndpointsModInterface().V1().Patch(ctx, name, pt, data, opts, subresources...)
}

func (e *endpointsInspectorAnonymousV1) Delete(ctx context.Context, name string, opts meta.DeleteOptions) error {
	return e.i.EndpointsModInterface().V1().Delete(ctx, name, opts)
}