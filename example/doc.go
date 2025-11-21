// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

//go:generate sh -c "bash $GARDENER_HACK_DIR/generate-crds.sh -p 20-crd- extensions.gardener.cloud resources.gardener.cloud"
//go:generate sh -c "$TOOLS_BIN_DIR/extension-generator --name=provider-ironcore-metal --provider-type=ironcore-metal --component-category=provider-extension --extension-oci-repository=ghcr.io/ironcore-dev/charts/gardener-extension-provider-ironcore-metal:$(cat ../VERSION) --admission-runtime-oci-repository=ghcr.io/ironcore-dev/charts/gardener-extension-admission-ironcore-metal-runtime:$(cat ../VERSION) --admission-application-oci-repository=ghcr.io/ironcore-dev/charts/gardener-extension-admission-ironcore-metal-application:$(cat ../VERSION) --destination=./extension/base/extension.yaml"
//go:generate sh -c "$TOOLS_BIN_DIR/kustomize build ./extension -o ./extension.yaml"

// Package example contains generated manifests for all CRDs and other examples.
// Useful for development purposes.
package example
