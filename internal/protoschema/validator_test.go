// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2026 Kernloom Contributors

package protoschema

import "testing"

func TestAdapterV1ProtoHasCanonicalShape(t *testing.T) {
	if err := ValidateAdapterV1("../../proto/kernloom/adapter/v1/adapter.proto"); err != nil {
		t.Fatal(err)
	}
}
