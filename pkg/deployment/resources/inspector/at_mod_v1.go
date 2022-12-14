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
	arangotaskv1 "github.com/arangodb/kube-arangodb/pkg/util/k8sutil/inspector/arangotask/v1"
	"github.com/arangodb/kube-arangodb/pkg/util/k8sutil/inspector/constants"
)

func (p arangoTaskMod) V1() arangotaskv1.ModInterface {
	return arangoTaskModV1(p)
}

type arangoTaskModV1 struct {
	i *inspectorState
}

func (p arangoTaskModV1) client() typedApi.ArangoTaskInterface {
	return p.i.Client().Arango().DatabaseV1().ArangoTasks(p.i.Namespace())
}

func (p arangoTaskModV1) Create(ctx context.Context, arangoTask *api.ArangoTask, opts meta.CreateOptions) (*api.ArangoTask, error) {
	logCreateOperation(constants.ArangoTaskKind, arangoTask)

	if arangoTask, err := p.client().Create(ctx, arangoTask, opts); err != nil {
		return arangoTask, err
	} else {
		p.i.GetThrottles().ArangoTask().Invalidate()
		return arangoTask, err
	}
}

func (p arangoTaskModV1) Update(ctx context.Context, arangoTask *api.ArangoTask, opts meta.UpdateOptions) (*api.ArangoTask, error) {
	logUpdateOperation(constants.ArangoTaskKind, arangoTask)

	if arangoTask, err := p.client().Update(ctx, arangoTask, opts); err != nil {
		return arangoTask, err
	} else {
		p.i.GetThrottles().ArangoTask().Invalidate()
		return arangoTask, err
	}
}

func (p arangoTaskModV1) UpdateStatus(ctx context.Context, arangoTask *api.ArangoTask, opts meta.UpdateOptions) (*api.ArangoTask, error) {
	logUpdateStatusOperation(constants.ArangoTaskKind, arangoTask)

	if arangoTask, err := p.client().UpdateStatus(ctx, arangoTask, opts); err != nil {
		return arangoTask, err
	} else {
		p.i.GetThrottles().ArangoTask().Invalidate()
		return arangoTask, err
	}
}

func (p arangoTaskModV1) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts meta.PatchOptions, subresources ...string) (result *api.ArangoTask, err error) {
	logPatchOperation(constants.ArangoTaskKind, p.i.Namespace(), name)

	if arangoTask, err := p.client().Patch(ctx, name, pt, data, opts, subresources...); err != nil {
		return arangoTask, err
	} else {
		p.i.GetThrottles().ArangoTask().Invalidate()
		return arangoTask, err
	}
}

func (p arangoTaskModV1) Delete(ctx context.Context, name string, opts meta.DeleteOptions) error {
	logDeleteOperation(constants.ArangoTaskKind, p.i.Namespace(), name, opts)

	if err := p.client().Delete(ctx, name, opts); err != nil {
		return err
	} else {
		p.i.GetThrottles().ArangoTask().Invalidate()
		return err
	}
}
