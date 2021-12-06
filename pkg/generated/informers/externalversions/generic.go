//
// DISCLAIMER
//
// Copyright 2016-2021 ArangoDB GmbH, Cologne, Germany
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

// Code generated by informer-gen. DO NOT EDIT.

package externalversions

import (
	"fmt"

	v1 "github.com/arangodb/kube-arangodb/pkg/apis/apps/v1"
	backupv1 "github.com/arangodb/kube-arangodb/pkg/apis/backup/v1"
	deploymentv1 "github.com/arangodb/kube-arangodb/pkg/apis/deployment/v1"
	v2alpha1 "github.com/arangodb/kube-arangodb/pkg/apis/deployment/v2alpha1"
	replicationv1 "github.com/arangodb/kube-arangodb/pkg/apis/replication/v1"
	replicationv2alpha1 "github.com/arangodb/kube-arangodb/pkg/apis/replication/v2alpha1"
	v1alpha "github.com/arangodb/kube-arangodb/pkg/apis/storage/v1alpha"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	cache "k8s.io/client-go/tools/cache"
)

// GenericInformer is type of SharedIndexInformer which will locate and delegate to other
// sharedInformers based on type
type GenericInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() cache.GenericLister
}

type genericInformer struct {
	informer cache.SharedIndexInformer
	resource schema.GroupResource
}

// Informer returns the SharedIndexInformer.
func (f *genericInformer) Informer() cache.SharedIndexInformer {
	return f.informer
}

// Lister returns the GenericLister.
func (f *genericInformer) Lister() cache.GenericLister {
	return cache.NewGenericLister(f.Informer().GetIndexer(), f.resource)
}

// ForResource gives generic access to a shared informer of the matching type
// TODO extend this to unknown resources with a client pool
func (f *sharedInformerFactory) ForResource(resource schema.GroupVersionResource) (GenericInformer, error) {
	switch resource {
	// Group=apps.arangodb.com, Version=v1
	case v1.SchemeGroupVersion.WithResource("arangojobs"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Apps().V1().ArangoJobs().Informer()}, nil

		// Group=backup.arangodb.com, Version=v1
	case backupv1.SchemeGroupVersion.WithResource("arangobackups"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Backup().V1().ArangoBackups().Informer()}, nil
	case backupv1.SchemeGroupVersion.WithResource("arangobackuppolicies"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Backup().V1().ArangoBackupPolicies().Informer()}, nil

		// Group=database.arangodb.com, Version=v1
	case deploymentv1.SchemeGroupVersion.WithResource("arangodeployments"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Database().V1().ArangoDeployments().Informer()}, nil
	case deploymentv1.SchemeGroupVersion.WithResource("arangomembers"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Database().V1().ArangoMembers().Informer()}, nil

		// Group=database.arangodb.com, Version=v2alpha1
	case v2alpha1.SchemeGroupVersion.WithResource("arangodeployments"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Database().V2alpha1().ArangoDeployments().Informer()}, nil
	case v2alpha1.SchemeGroupVersion.WithResource("arangomembers"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Database().V2alpha1().ArangoMembers().Informer()}, nil

		// Group=replication.database.arangodb.com, Version=v1
	case replicationv1.SchemeGroupVersion.WithResource("arangodeploymentreplications"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Replication().V1().ArangoDeploymentReplications().Informer()}, nil

		// Group=replication.database.arangodb.com, Version=v2alpha1
	case replicationv2alpha1.SchemeGroupVersion.WithResource("arangodeploymentreplications"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Replication().V2alpha1().ArangoDeploymentReplications().Informer()}, nil

		// Group=storage.arangodb.com, Version=v1alpha
	case v1alpha.SchemeGroupVersion.WithResource("arangolocalstorages"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Storage().V1alpha().ArangoLocalStorages().Informer()}, nil

	}

	return nil, fmt.Errorf("no informer found for %v", resource)
}
