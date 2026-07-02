// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2026 Kernloom Contributors

package contractv1

import (
	"context"
	"testing"

	adapterv1 "github.com/kernloom/kernloom-protocol/sdk/go/adapter/v1"
)

type emptyAdapter struct{}

func (emptyAdapter) Describe(context.Context) (*adapterv1.AdapterDescriptor, error) {
	return &adapterv1.AdapterDescriptor{
		AdapterId:       "empty.test",
		Name:            "Empty Contract Adapter",
		ProtocolVersion: adapterv1.ProtocolVersion,
		Facets: []string{
			adapterv1.FacetDescribe,
			adapterv1.FacetHealth,
		},
	}, nil
}

func (emptyAdapter) Health(context.Context) (*adapterv1.HealthResponse, error) {
	return &adapterv1.HealthResponse{Status: adapterv1.HealthServing}, nil
}

func TestEmptyAdapterPassesMinimalContract(t *testing.T) {
	RunMinimalContract(t, emptyAdapter{})
}
