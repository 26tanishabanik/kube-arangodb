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

// Code generated by client-gen. DO NOT EDIT.

package v2alpha1

import (
	v2alpha1 "github.com/arangodb/kube-arangodb/pkg/apis/deployment/v2alpha1"
	"github.com/arangodb/kube-arangodb/pkg/generated/clientset/versioned/scheme"
	rest "k8s.io/client-go/rest"
)

type DatabaseV2alpha1Interface interface {
	RESTClient() rest.Interface
	ArangoClusterSynchronizationsGetter
	ArangoDeploymentsGetter
	ArangoMembersGetter
}

// DatabaseV2alpha1Client is used to interact with features provided by the database.arangodb.com group.
type DatabaseV2alpha1Client struct {
	restClient rest.Interface
}

func (c *DatabaseV2alpha1Client) ArangoClusterSynchronizations(namespace string) ArangoClusterSynchronizationInterface {
	return newArangoClusterSynchronizations(c, namespace)
}

func (c *DatabaseV2alpha1Client) ArangoDeployments(namespace string) ArangoDeploymentInterface {
	return newArangoDeployments(c, namespace)
}

func (c *DatabaseV2alpha1Client) ArangoMembers(namespace string) ArangoMemberInterface {
	return newArangoMembers(c, namespace)
}

// NewForConfig creates a new DatabaseV2alpha1Client for the given config.
func NewForConfig(c *rest.Config) (*DatabaseV2alpha1Client, error) {
	config := *c
	if err := setConfigDefaults(&config); err != nil {
		return nil, err
	}
	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}
	return &DatabaseV2alpha1Client{client}, nil
}

// NewForConfigOrDie creates a new DatabaseV2alpha1Client for the given config and
// panics if there is an error in the config.
func NewForConfigOrDie(c *rest.Config) *DatabaseV2alpha1Client {
	client, err := NewForConfig(c)
	if err != nil {
		panic(err)
	}
	return client
}

// New creates a new DatabaseV2alpha1Client for the given RESTClient.
func New(c rest.Interface) *DatabaseV2alpha1Client {
	return &DatabaseV2alpha1Client{c}
}

func setConfigDefaults(config *rest.Config) error {
	gv := v2alpha1.SchemeGroupVersion
	config.GroupVersion = &gv
	config.APIPath = "/apis"
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()

	if config.UserAgent == "" {
		config.UserAgent = rest.DefaultKubernetesUserAgent()
	}

	return nil
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *DatabaseV2alpha1Client) RESTClient() rest.Interface {
	if c == nil {
		return nil
	}
	return c.restClient
}
