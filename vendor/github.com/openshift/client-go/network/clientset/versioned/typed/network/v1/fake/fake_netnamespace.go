package fake

import (
	network_v1 "github.com/openshift/api/network/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeNetNamespaces implements NetNamespaceInterface
type FakeNetNamespaces struct {
	Fake *FakeNetworkV1
}

var netnamespacesResource = schema.GroupVersionResource{Group: "network.openshift.io", Version: "v1", Resource: "netnamespaces"}

var netnamespacesKind = schema.GroupVersionKind{Group: "network.openshift.io", Version: "v1", Kind: "NetNamespace"}

// Get takes name of the netNamespace, and returns the corresponding netNamespace object, and an error if there is any.
func (c *FakeNetNamespaces) Get(name string, options v1.GetOptions) (result *network_v1.NetNamespace, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(netnamespacesResource, name), &network_v1.NetNamespace{})
	if obj == nil {
		return nil, err
	}
	return obj.(*network_v1.NetNamespace), err
}

// List takes label and field selectors, and returns the list of NetNamespaces that match those selectors.
func (c *FakeNetNamespaces) List(opts v1.ListOptions) (result *network_v1.NetNamespaceList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(netnamespacesResource, netnamespacesKind, opts), &network_v1.NetNamespaceList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &network_v1.NetNamespaceList{}
	for _, item := range obj.(*network_v1.NetNamespaceList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested netNamespaces.
func (c *FakeNetNamespaces) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(netnamespacesResource, opts))
}

// Create takes the representation of a netNamespace and creates it.  Returns the server's representation of the netNamespace, and an error, if there is any.
func (c *FakeNetNamespaces) Create(netNamespace *network_v1.NetNamespace) (result *network_v1.NetNamespace, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(netnamespacesResource, netNamespace), &network_v1.NetNamespace{})
	if obj == nil {
		return nil, err
	}
	return obj.(*network_v1.NetNamespace), err
}

// Update takes the representation of a netNamespace and updates it. Returns the server's representation of the netNamespace, and an error, if there is any.
func (c *FakeNetNamespaces) Update(netNamespace *network_v1.NetNamespace) (result *network_v1.NetNamespace, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(netnamespacesResource, netNamespace), &network_v1.NetNamespace{})
	if obj == nil {
		return nil, err
	}
	return obj.(*network_v1.NetNamespace), err
}

// Delete takes name of the netNamespace and deletes it. Returns an error if one occurs.
func (c *FakeNetNamespaces) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteAction(netnamespacesResource, name), &network_v1.NetNamespace{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeNetNamespaces) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(netnamespacesResource, listOptions)

	_, err := c.Fake.Invokes(action, &network_v1.NetNamespaceList{})
	return err
}

// Patch applies the patch and returns the patched netNamespace.
func (c *FakeNetNamespaces) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *network_v1.NetNamespace, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(netnamespacesResource, name, data, subresources...), &network_v1.NetNamespace{})
	if obj == nil {
		return nil, err
	}
	return obj.(*network_v1.NetNamespace), err
}
