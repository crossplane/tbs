/*

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

package controllers

import (
	"context"
	"fmt"
	"strings"

	"github.com/nlopes/slack"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	runtimev1alpha1 "github.com/crossplaneio/crossplane-runtime/apis/core/v1alpha1"
	"github.com/crossplaneio/crossplane-runtime/pkg/meta"
	"github.com/crossplaneio/crossplane-runtime/pkg/resource"

	slackv1 "github.com/crossplaneio/tbs/episodes/1/assets/api/v1"
	"github.com/crossplaneio/tbs/episodes/1/assets/clients"
)

var (
	errNotMessage    = "managed kind was not a message"
	errMessageCreate = "unable to create message"
	errNewClient     = "unable to create a new slack client"
)

// +kubebuilder:rbac:groups=slack.crossplane.io,resources=messages,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=slack.crossplane.io,resources=messages/status,verbs=get;update;patch

// MessageController is responsible for adding the Message
// controller and its corresponding reconciler to the manager with any runtime configuration.
type MessageController struct{}

// SetupWithManager creates a new Message Controller and adds it to the
// Manager with default RBAC. The Manager will set fields on the Controller and
// start it when the Manager is Started.
func (c *MessageController) SetupWithManager(mgr ctrl.Manager) error {
	r := resource.NewManagedReconciler(mgr,
		resource.ManagedKind(slackv1.MessageGroupVersionKind),
		resource.WithManagedConnectionPublishers(),
		resource.WithExternalConnecter(&connecter{client: mgr.GetClient()}))

	name := strings.ToLower(fmt.Sprintf("%s.%s", slackv1.MessageKind, slackv1.Group))

	return ctrl.NewControllerManagedBy(mgr).
		Named(name).
		For(&slackv1.Message{}).
		Complete(r)
}

type connecter struct {
	client      client.Client
	newClientFn func(credentials []byte) *slack.Client
}

func (c *connecter) Connect(ctx context.Context, mg resource.Managed) (resource.ExternalClient, error) {
	m, ok := mg.(*slackv1.Message)
	if !ok {
		return nil, errors.New(errNotMessage)
	}

	p := &slackv1.Provider{}
	n := meta.NamespacedNameOf(m.Spec.ProviderReference)
	if err := c.client.Get(ctx, n, p); err != nil {
		return nil, errors.Wrapf(err, "cannot get provider %s", n)
	}

	s := &corev1.Secret{}
	n = types.NamespacedName{Namespace: p.GetNamespace(), Name: p.Spec.Secret.Name}
	if err := c.client.Get(ctx, n, s); err != nil {
		return nil, errors.Wrapf(err, "cannot get provider secret %s", n)
	}
	newClientFn := clients.NewClient
	if c.newClientFn != nil {
		newClientFn = c.newClientFn
	}
	client := newClientFn(s.Data[p.Spec.Secret.Key])
	return &external{client: client}, errors.Wrap(nil, errNewClient)
}

type external struct{ client *slack.Client }

func (e *external) Observe(ctx context.Context, mg resource.Managed) (resource.ExternalObservation, error) {
	m, ok := mg.(*slackv1.Message)
	if !ok {
		return resource.ExternalObservation{}, errors.New(errNotMessage)
	}

	if m.Status.Sent != true {
		return resource.ExternalObservation{
			ResourceExists: false,
		}, nil
	}

	return resource.ExternalObservation{ResourceExists: true}, nil
}

func (e *external) Create(ctx context.Context, mg resource.Managed) (resource.ExternalCreation, error) {
	m, ok := mg.(*slackv1.Message)
	if !ok {
		return resource.ExternalCreation{}, errors.New(errNotMessage)
	}

	m.Status.SetConditions(runtimev1alpha1.Creating())

	_, _, err := e.client.PostMessage(m.Spec.Channel, slack.MsgOptionText(m.Spec.Text, false))
	if err != nil {
		return resource.ExternalCreation{}, errors.Wrap(err, errMessageCreate)
	}

	m.Status.Sent = true

	return resource.ExternalCreation{}, nil
}

func (e *external) Update(ctx context.Context, mg resource.Managed) (resource.ExternalUpdate, error) {
	// m, ok := mg.(*slackv1.Message)
	// if !ok {
	// 	return resource.ExternalUpdate{}, errors.New(errNotMessage)
	// }

	return resource.ExternalUpdate{}, nil
}

func (e *external) Delete(ctx context.Context, mg resource.Managed) error {
	// m, ok := mg.(*slackv1.Message)
	// if !ok {
	// 	return errors.New(errNotMessage)
	// }

	return nil
}
