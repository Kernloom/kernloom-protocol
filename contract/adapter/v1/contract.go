// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2026 Kernloom Contributors

package contractv1

import (
	"context"
	"testing"
	"time"

	adapterv1 "github.com/kernloom/kernloom-protocol/sdk/go/adapter/v1"
)

func RunMinimalContract(t *testing.T, impl adapterv1.MinimalAdapter) {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	desc, err := impl.Describe(ctx)
	if err != nil {
		t.Fatalf("Describe returned error: %v", err)
	}
	if err := adapterv1.ValidateDescriptor(desc); err != nil {
		t.Fatalf("Describe returned invalid descriptor: %v", err)
	}

	health, err := impl.Health(ctx)
	if err != nil {
		t.Fatalf("Health returned error: %v", err)
	}
	switch health.Status {
	case adapterv1.HealthServing, adapterv1.HealthDegraded:
	default:
		t.Fatalf("Health returned non-serving status %q", health.Status)
	}
}
