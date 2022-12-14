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

	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	api "github.com/arangodb/kube-arangodb/pkg/apis/deployment/v1"
	typedApi "github.com/arangodb/kube-arangodb/pkg/generated/clientset/versioned/typed/deployment/v1"
	arangoclustersynchronizationv1 "github.com/arangodb/kube-arangodb/pkg/util/k8sutil/inspector/arangoclustersynchronization/v1"
	"github.com/arangodb/kube-arangodb/pkg/util/k8sutil/inspector/constants"
)

func (p arangoClusterSynchronizationMod) V1() arangoclustersynchronizationv1.ModInterface {
	return arangoClusterSynchronizationModV1(p)
}

type arangoClusterSynchronizationModV1 struct {
	i *inspectorState
}

func (p arangoClusterSynchronizationModV1) client() typedApi.ArangoClusterSynchronizationInterface {
	return p.i.Client().Arango().DatabaseV1().ArangoClusterSynchronizations(p.i.Namespace())
}

func (p arangoClusterSynchronizationModV1) Create(ctx context.Context, arangoClusterSynchronization *api.ArangoClusterSynchronization, opts meta.CreateOptions) (*api.ArangoClusterSynchronization, error) {
	logCreateOperation(constants.ArangoClusterSynchronizationKind, arangoClusterSynchronization)

	if arangoClusterSynchronization, err := p.client().Create(ctx, arangoClusterSynchronization, opts); err != nil {
		return arangoClusterSynchronization, err
	} else {
		p.i.GetThrottles().ArangoClusterSynchronization().Invalidate()
		return arangoClusterSynchronization, err
	}
}

func (p arangoClusterSynchronizationModV1) Update(ctx context.Context, arangoClusterSynchronization *api.ArangoClusterSynchronization, opts meta.UpdateOptions) (*api.ArangoClusterSynchronization, error) {
	logUpdateOperation(constants.ArangoClusterSynchronizationKind, arangoClusterSynchronization)

	if arangoClusterSynchronization, err := p.client().Update(ctx, arangoClusterSynchronization, opts); err != nil {
		return arangoClusterSynchronization, err
	} else {
		p.i.GetThrottles().ArangoClusterSynchronization().Invalidate()
		return arangoClusterSynchronization, err
	}
}

func (p arangoClusterSynchronizationModV1) UpdateStatus(ctx context.Context, arangoClusterSynchronization *api.ArangoClusterSynchronization, opts meta.UpdateOptions) (*api.ArangoClusterSynchronization, error) {
	logUpdateStatusOperation(constants.ArangoClusterSynchronizationKind, arangoClusterSynchronization)

	if arangoClusterSynchronization, err := p.client().UpdateStatus(ctx, arangoClusterSynchronization, opts); err != nil {
		return arangoClusterSynchronization, err
	} else {
		p.i.GetThrottles().ArangoClusterSynchronization().Invalidate()
		return arangoClusterSynchronization, err
	}
}

func (p arangoClusterSynchronizationModV1) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts meta.PatchOptions, subresources ...string) (result *api.ArangoClusterSynchronization, err error) {
	logPatchOperation(constants.ArangoClusterSynchronizationKind, p.i.Namespace(), name)

	if arangoClusterSynchronization, err := p.client().Patch(ctx, name, pt, data, opts, subresources...); err != nil {
		return arangoClusterSynchronization, err
	} else {
		p.i.GetThrottles().ArangoClusterSynchronization().Invalidate()
		return arangoClusterSynchronization, err
	}
}

func (p arangoClusterSynchronizationModV1) Delete(ctx context.Context, name string, opts meta.DeleteOptions) error {
	logDeleteOperation(constants.ArangoClusterSynchronizationKind, p.i.Namespace(), name, opts)

	if err := p.client().Delete(ctx, name, opts); err != nil {
		return err
	} else {
		p.i.GetThrottles().ArangoClusterSynchronization().Invalidate()
		return err
	}
}
