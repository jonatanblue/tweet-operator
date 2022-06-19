/*
Copyright The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	examplecomv1 "github.com/jonatanblue/tweet-operator/pkg/apis/example.com/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeTweets implements TweetInterface
type FakeTweets struct {
	Fake *FakeExampleV1
	ns   string
}

var tweetsResource = schema.GroupVersionResource{Group: "example.com", Version: "v1", Resource: "tweets"}

var tweetsKind = schema.GroupVersionKind{Group: "example.com", Version: "v1", Kind: "Tweet"}

// Get takes name of the tweet, and returns the corresponding tweet object, and an error if there is any.
func (c *FakeTweets) Get(ctx context.Context, name string, options v1.GetOptions) (result *examplecomv1.Tweet, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(tweetsResource, c.ns, name), &examplecomv1.Tweet{})

	if obj == nil {
		return nil, err
	}
	return obj.(*examplecomv1.Tweet), err
}

// List takes label and field selectors, and returns the list of Tweets that match those selectors.
func (c *FakeTweets) List(ctx context.Context, opts v1.ListOptions) (result *examplecomv1.TweetList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(tweetsResource, tweetsKind, c.ns, opts), &examplecomv1.TweetList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &examplecomv1.TweetList{ListMeta: obj.(*examplecomv1.TweetList).ListMeta}
	for _, item := range obj.(*examplecomv1.TweetList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested tweets.
func (c *FakeTweets) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(tweetsResource, c.ns, opts))

}

// Create takes the representation of a tweet and creates it.  Returns the server's representation of the tweet, and an error, if there is any.
func (c *FakeTweets) Create(ctx context.Context, tweet *examplecomv1.Tweet, opts v1.CreateOptions) (result *examplecomv1.Tweet, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(tweetsResource, c.ns, tweet), &examplecomv1.Tweet{})

	if obj == nil {
		return nil, err
	}
	return obj.(*examplecomv1.Tweet), err
}

// Update takes the representation of a tweet and updates it. Returns the server's representation of the tweet, and an error, if there is any.
func (c *FakeTweets) Update(ctx context.Context, tweet *examplecomv1.Tweet, opts v1.UpdateOptions) (result *examplecomv1.Tweet, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(tweetsResource, c.ns, tweet), &examplecomv1.Tweet{})

	if obj == nil {
		return nil, err
	}
	return obj.(*examplecomv1.Tweet), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeTweets) UpdateStatus(ctx context.Context, tweet *examplecomv1.Tweet, opts v1.UpdateOptions) (*examplecomv1.Tweet, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(tweetsResource, "status", c.ns, tweet), &examplecomv1.Tweet{})

	if obj == nil {
		return nil, err
	}
	return obj.(*examplecomv1.Tweet), err
}

// Delete takes name of the tweet and deletes it. Returns an error if one occurs.
func (c *FakeTweets) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(tweetsResource, c.ns, name), &examplecomv1.Tweet{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeTweets) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(tweetsResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &examplecomv1.TweetList{})
	return err
}

// Patch applies the patch and returns the patched tweet.
func (c *FakeTweets) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *examplecomv1.Tweet, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(tweetsResource, c.ns, name, pt, data, subresources...), &examplecomv1.Tweet{})

	if obj == nil {
		return nil, err
	}
	return obj.(*examplecomv1.Tweet), err
}
