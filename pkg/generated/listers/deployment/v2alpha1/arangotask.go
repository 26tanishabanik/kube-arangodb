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

// Code generated by lister-gen. DO NOT EDIT.

package v2alpha1

import (
	v2alpha1 "github.com/arangodb/kube-arangodb/pkg/apis/deployment/v2alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// ArangoTaskLister helps list ArangoTasks.
// All objects returned here must be treated as read-only.
type ArangoTaskLister interface {
	// List lists all ArangoTasks in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v2alpha1.ArangoTask, err error)
	// ArangoTasks returns an object that can list and get ArangoTasks.
	ArangoTasks(namespace string) ArangoTaskNamespaceLister
	ArangoTaskListerExpansion
}

// arangoTaskLister implements the ArangoTaskLister interface.
type arangoTaskLister struct {
	indexer cache.Indexer
}

// NewArangoTaskLister returns a new ArangoTaskLister.
func NewArangoTaskLister(indexer cache.Indexer) ArangoTaskLister {
	return &arangoTaskLister{indexer: indexer}
}

// List lists all ArangoTasks in the indexer.
func (s *arangoTaskLister) List(selector labels.Selector) (ret []*v2alpha1.ArangoTask, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v2alpha1.ArangoTask))
	})
	return ret, err
}

// ArangoTasks returns an object that can list and get ArangoTasks.
func (s *arangoTaskLister) ArangoTasks(namespace string) ArangoTaskNamespaceLister {
	return arangoTaskNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// ArangoTaskNamespaceLister helps list and get ArangoTasks.
// All objects returned here must be treated as read-only.
type ArangoTaskNamespaceLister interface {
	// List lists all ArangoTasks in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v2alpha1.ArangoTask, err error)
	// Get retrieves the ArangoTask from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v2alpha1.ArangoTask, error)
	ArangoTaskNamespaceListerExpansion
}

// arangoTaskNamespaceLister implements the ArangoTaskNamespaceLister
// interface.
type arangoTaskNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all ArangoTasks in the indexer for a given namespace.
func (s arangoTaskNamespaceLister) List(selector labels.Selector) (ret []*v2alpha1.ArangoTask, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v2alpha1.ArangoTask))
	})
	return ret, err
}

// Get retrieves the ArangoTask from the indexer for a given namespace and name.
func (s arangoTaskNamespaceLister) Get(name string) (*v2alpha1.ArangoTask, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v2alpha1.Resource("arangotask"), name)
	}
	return obj.(*v2alpha1.ArangoTask), nil
}
