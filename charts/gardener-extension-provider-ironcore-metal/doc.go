// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and IronCore contributors
// SPDX-License-Identifier: Apache-2.0

//go:generate sh -c "bash $GARDENER_HACK_DIR/generate-controller-registration.sh provider-ironcore-metal . $(cat ../../VERSION) ../../example/controller-registration.yaml ControlPlane:ironcore-metal Infrastructure:ironcore-metal Worker:ironcore-metal"

// Package chart enables go:generate support for generating the correct controller registration.
package chart
