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

package patcher

import (
	"context"

	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"

	"github.com/arangodb/kube-arangodb/pkg/deployment/patch"
	"github.com/arangodb/kube-arangodb/pkg/util/globals"
	v1 "github.com/arangodb/kube-arangodb/pkg/util/k8sutil/inspector/service/v1"
)

type ServicePatch func(in *core.Service) []patch.Item

func ServicePatcher(ctx context.Context, client v1.ModInterface, in *core.Service, opts meta.PatchOptions, functions ...ServicePatch) (bool, error) {
	if in == nil {
		return false, nil
	}

	if in.GetName() == "" {
		return false, nil
	}

	var items []patch.Item

	for id := range functions {
		if f := functions[id]; f != nil {
			items = append(items, f(in)...)
		}
	}

	if len(items) == 0 {
		return false, nil
	}

	data, err := patch.NewPatch(items...).Marshal()
	if err != nil {
		return false, err
	}

	nctx, c := globals.GetGlobals().Timeouts().Kubernetes().WithTimeout(ctx)
	defer c()

	if _, err := client.Patch(nctx, in.GetName(), types.JSONPatchType, data, opts); err != nil {
		return false, err
	}

	return true, nil
}

func PatchServicePorts(ports []core.ServicePort) ServicePatch {
	return func(in *core.Service) []patch.Item {
		if len(ports) == len(in.Spec.Ports) && equality.Semantic.DeepDerivative(ports, in.Spec.Ports) {
			return nil
		}

		return []patch.Item{
			patch.ItemReplace(patch.NewPath("spec", "ports"), ports),
		}
	}
}

func PatchServiceSelector(selector map[string]string) ServicePatch {
	return func(in *core.Service) []patch.Item {
		if equality.Semantic.DeepEqual(in.Spec.Selector, selector) {
			return nil
		}

		return []patch.Item{
			patch.ItemReplace(patch.NewPath("spec", "selector"), selector),
		}
	}
}

func PatchServiceType(t core.ServiceType) ServicePatch {
	return func(in *core.Service) []patch.Item {
		if in.Spec.Type == t {
			return nil
		}

		return []patch.Item{
			patch.ItemReplace(patch.NewPath("spec", "type"), t),
		}
	}
}

func PatchServicePublishNotReadyAddresses(publishNotReadyAddresses bool) ServicePatch {
	return func(in *core.Service) []patch.Item {
		if in.Spec.PublishNotReadyAddresses == publishNotReadyAddresses {
			return nil
		}

		return []patch.Item{
			patch.ItemReplace(patch.NewPath("spec", "publishNotReadyAddresses"), publishNotReadyAddresses),
		}
	}
}