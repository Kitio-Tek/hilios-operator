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

// Package statuswriter centralises the retry-on-conflict pattern used by the
// HILIOS reconcilers when they update a resource's status subresource. The
// helpers reduce flapping under load by re-reading the resource and re-applying
// the mutator on conflict instead of dropping the update.
package statuswriter

import (
	"context"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/util/retry"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// UpdateStatus runs mutate against a fresh copy of obj and persists the change
// through the status subresource. On conflict the helper retries with
// retry.DefaultRetry. obj must implement client.Object and have a Status field.
func UpdateStatus[T client.Object](ctx context.Context, c client.Client, obj T, mutate func(T) error) error {
	return retry.RetryOnConflict(retry.DefaultRetry, func() error {
		key := client.ObjectKeyFromObject(obj)
		fresh := obj.DeepCopyObject().(T)
		if err := c.Get(ctx, key, fresh); err != nil {
			return err
		}
		if err := mutate(fresh); err != nil {
			return err
		}
		if err := c.Status().Update(ctx, fresh); err != nil {
			if apierrors.IsConflict(err) {
				return err
			}
			return err
		}
		return nil
	})
}
