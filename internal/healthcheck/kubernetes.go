/*
Copyright 2026.

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

package healthcheck

import (
	"context"
	"fmt"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client"

	resiliencev1alpha1 "github.com/Kitio-Tek/hilios-operator/api/v1alpha1"
)

// RunKubernetes verifies that the Kubernetes resource referenced by the probe
// exists, fetching it through the supplied controller-runtime client. It does
// not assert anything about the object's status — controllers that need
// readiness checks should compose this with a Pod or HTTP probe.
func RunKubernetes(ctx context.Context, c client.Client, hc resiliencev1alpha1.HealthCheck) error {
	if hc.Type != resiliencev1alpha1.ProbeTypeKubernetes {
		return fmt.Errorf("kubernetes probe %q has unexpected type %q", hc.Name, hc.Type)
	}
	if hc.Resource == nil {
		return fmt.Errorf("kubernetes probe %q is missing spec.resource", hc.Name)
	}
	gvk, err := parseGVK(hc.Resource.APIVersion, hc.Resource.Kind)
	if err != nil {
		return fmt.Errorf("kubernetes probe %q: %w", hc.Name, err)
	}
	obj := &unstructured.Unstructured{}
	obj.SetGroupVersionKind(gvk)

	if err := c.Get(ctx, client.ObjectKey{Namespace: hc.Resource.Namespace, Name: hc.Resource.Name}, obj); err != nil {
		if apierrors.IsNotFound(err) {
			return fmt.Errorf("kubernetes probe %q: %s/%s %q not found", hc.Name, gvk.Group, gvk.Kind, hc.Resource.Name)
		}
		return fmt.Errorf("kubernetes probe %q: %w", hc.Name, err)
	}
	return nil
}

func parseGVK(apiVersion, kind string) (schema.GroupVersionKind, error) {
	if apiVersion == "" {
		return schema.GroupVersionKind{}, fmt.Errorf("apiVersion is required")
	}
	gv, err := schema.ParseGroupVersion(apiVersion)
	if err != nil {
		return schema.GroupVersionKind{}, fmt.Errorf("parse apiVersion %q: %w", apiVersion, err)
	}
	if kind == "" {
		return schema.GroupVersionKind{}, fmt.Errorf("kind is required")
	}
	return gv.WithKind(kind), nil
}
