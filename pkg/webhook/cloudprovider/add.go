// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

package cloudprovider

import (
	extensionswebhook "github.com/gardener/gardener/extensions/pkg/webhook"
	"github.com/gardener/gardener/extensions/pkg/webhook/cloudprovider"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	"github.com/ironcore-dev/gardener-extension-provider-ironcore-metal/pkg/metal"
)

var (
	logger = log.Log.WithName("metal-cloudprovider-webhook")

	// DefaultAddOptions are the default AddOptions for AddToManager.
	DefaultAddOptions = AddOptions{}
)

// AddOptions are options to apply when adding the cloudprovider webhook to the manager.
type AddOptions struct {
	// MetalNamespace is the namespace in the Metal-API cluster used as the kubeconfig context namespace.
	MetalNamespace string
}

// AddToManagerWithOptions creates the cloudprovider webhook with the given options and adds it to the manager.
func AddToManagerWithOptions(mgr manager.Manager, opts AddOptions) (*extensionswebhook.Webhook, error) {
	logger.Info("adding webhook to manager")
	return cloudprovider.New(mgr, cloudprovider.Args{
		Provider: metal.Type,
		Mutator:  cloudprovider.NewMutator(mgr, logger, NewEnsurer(logger, mgr, opts.MetalNamespace)),
	})
}

// AddToManager creates the cloudprovider webhook and adds it to the manager.
func AddToManager(mgr manager.Manager) (*extensionswebhook.Webhook, error) {
	return AddToManagerWithOptions(mgr, DefaultAddOptions)
}
