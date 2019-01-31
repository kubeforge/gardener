/*
Copyright (c) 2019 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file

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

package v1alpha1

import (
	"time"

	scheme "github.com/gardener/gardener/pkg/client/machine/clientset/versioned/scheme"
	v1alpha1 "github.com/gardener/machine-controller-manager/pkg/apis/machine/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// AlicloudMachineClassesGetter has a method to return a AlicloudMachineClassInterface.
// A group's client should implement this interface.
type AlicloudMachineClassesGetter interface {
	AlicloudMachineClasses(namespace string) AlicloudMachineClassInterface
}

// AlicloudMachineClassInterface has methods to work with AlicloudMachineClass resources.
type AlicloudMachineClassInterface interface {
	Create(*v1alpha1.AlicloudMachineClass) (*v1alpha1.AlicloudMachineClass, error)
	Update(*v1alpha1.AlicloudMachineClass) (*v1alpha1.AlicloudMachineClass, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha1.AlicloudMachineClass, error)
	List(opts v1.ListOptions) (*v1alpha1.AlicloudMachineClassList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.AlicloudMachineClass, err error)
	AlicloudMachineClassExpansion
}

// alicloudMachineClasses implements AlicloudMachineClassInterface
type alicloudMachineClasses struct {
	client rest.Interface
	ns     string
}

// newAlicloudMachineClasses returns a AlicloudMachineClasses
func newAlicloudMachineClasses(c *MachineV1alpha1Client, namespace string) *alicloudMachineClasses {
	return &alicloudMachineClasses{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the alicloudMachineClass, and returns the corresponding alicloudMachineClass object, and an error if there is any.
func (c *alicloudMachineClasses) Get(name string, options v1.GetOptions) (result *v1alpha1.AlicloudMachineClass, err error) {
	result = &v1alpha1.AlicloudMachineClass{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("alicloudmachineclasses").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of AlicloudMachineClasses that match those selectors.
func (c *alicloudMachineClasses) List(opts v1.ListOptions) (result *v1alpha1.AlicloudMachineClassList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.AlicloudMachineClassList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("alicloudmachineclasses").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested alicloudMachineClasses.
func (c *alicloudMachineClasses) Watch(opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("alicloudmachineclasses").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch()
}

// Create takes the representation of a alicloudMachineClass and creates it.  Returns the server's representation of the alicloudMachineClass, and an error, if there is any.
func (c *alicloudMachineClasses) Create(alicloudMachineClass *v1alpha1.AlicloudMachineClass) (result *v1alpha1.AlicloudMachineClass, err error) {
	result = &v1alpha1.AlicloudMachineClass{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("alicloudmachineclasses").
		Body(alicloudMachineClass).
		Do().
		Into(result)
	return
}

// Update takes the representation of a alicloudMachineClass and updates it. Returns the server's representation of the alicloudMachineClass, and an error, if there is any.
func (c *alicloudMachineClasses) Update(alicloudMachineClass *v1alpha1.AlicloudMachineClass) (result *v1alpha1.AlicloudMachineClass, err error) {
	result = &v1alpha1.AlicloudMachineClass{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("alicloudmachineclasses").
		Name(alicloudMachineClass.Name).
		Body(alicloudMachineClass).
		Do().
		Into(result)
	return
}

// Delete takes name of the alicloudMachineClass and deletes it. Returns an error if one occurs.
func (c *alicloudMachineClasses) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("alicloudmachineclasses").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *alicloudMachineClasses) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	var timeout time.Duration
	if listOptions.TimeoutSeconds != nil {
		timeout = time.Duration(*listOptions.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("alicloudmachineclasses").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Timeout(timeout).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched alicloudMachineClass.
func (c *alicloudMachineClasses) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.AlicloudMachineClass, err error) {
	result = &v1alpha1.AlicloudMachineClass{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("alicloudmachineclasses").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
