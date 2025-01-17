//
// DISCLAIMER
//
// Copyright 2023 ArangoDB GmbH, Cologne, Germany
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

package fake

import (
	"context"

	v1alpha1 "github.com/arangodb/kube-arangodb/pkg/apis/ml/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeArangoMLBatchJobs implements ArangoMLBatchJobInterface
type FakeArangoMLBatchJobs struct {
	Fake *FakeMlV1alpha1
	ns   string
}

var arangomlbatchjobsResource = schema.GroupVersionResource{Group: "ml.arangodb.com", Version: "v1alpha1", Resource: "arangomlbatchjobs"}

var arangomlbatchjobsKind = schema.GroupVersionKind{Group: "ml.arangodb.com", Version: "v1alpha1", Kind: "ArangoMLBatchJob"}

// Get takes name of the arangoMLBatchJob, and returns the corresponding arangoMLBatchJob object, and an error if there is any.
func (c *FakeArangoMLBatchJobs) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.ArangoMLBatchJob, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(arangomlbatchjobsResource, c.ns, name), &v1alpha1.ArangoMLBatchJob{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ArangoMLBatchJob), err
}

// List takes label and field selectors, and returns the list of ArangoMLBatchJobs that match those selectors.
func (c *FakeArangoMLBatchJobs) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.ArangoMLBatchJobList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(arangomlbatchjobsResource, arangomlbatchjobsKind, c.ns, opts), &v1alpha1.ArangoMLBatchJobList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.ArangoMLBatchJobList{ListMeta: obj.(*v1alpha1.ArangoMLBatchJobList).ListMeta}
	for _, item := range obj.(*v1alpha1.ArangoMLBatchJobList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested arangoMLBatchJobs.
func (c *FakeArangoMLBatchJobs) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(arangomlbatchjobsResource, c.ns, opts))

}

// Create takes the representation of a arangoMLBatchJob and creates it.  Returns the server's representation of the arangoMLBatchJob, and an error, if there is any.
func (c *FakeArangoMLBatchJobs) Create(ctx context.Context, arangoMLBatchJob *v1alpha1.ArangoMLBatchJob, opts v1.CreateOptions) (result *v1alpha1.ArangoMLBatchJob, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(arangomlbatchjobsResource, c.ns, arangoMLBatchJob), &v1alpha1.ArangoMLBatchJob{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ArangoMLBatchJob), err
}

// Update takes the representation of a arangoMLBatchJob and updates it. Returns the server's representation of the arangoMLBatchJob, and an error, if there is any.
func (c *FakeArangoMLBatchJobs) Update(ctx context.Context, arangoMLBatchJob *v1alpha1.ArangoMLBatchJob, opts v1.UpdateOptions) (result *v1alpha1.ArangoMLBatchJob, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(arangomlbatchjobsResource, c.ns, arangoMLBatchJob), &v1alpha1.ArangoMLBatchJob{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ArangoMLBatchJob), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeArangoMLBatchJobs) UpdateStatus(ctx context.Context, arangoMLBatchJob *v1alpha1.ArangoMLBatchJob, opts v1.UpdateOptions) (*v1alpha1.ArangoMLBatchJob, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(arangomlbatchjobsResource, "status", c.ns, arangoMLBatchJob), &v1alpha1.ArangoMLBatchJob{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ArangoMLBatchJob), err
}

// Delete takes name of the arangoMLBatchJob and deletes it. Returns an error if one occurs.
func (c *FakeArangoMLBatchJobs) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteActionWithOptions(arangomlbatchjobsResource, c.ns, name, opts), &v1alpha1.ArangoMLBatchJob{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeArangoMLBatchJobs) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(arangomlbatchjobsResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.ArangoMLBatchJobList{})
	return err
}

// Patch applies the patch and returns the patched arangoMLBatchJob.
func (c *FakeArangoMLBatchJobs) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.ArangoMLBatchJob, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(arangomlbatchjobsResource, c.ns, name, pt, data, subresources...), &v1alpha1.ArangoMLBatchJob{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.ArangoMLBatchJob), err
}
