/*
 * Copyright © 2026 OpenCHAMI a Series of LF Projects, LLC
 *
 * SPDX-License-Identifier: MIT
 */

package v1

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/openchami/fabrica/pkg/fabrica"
)

func validComponent() *Component {
	enabled := true
	return &Component{
		APIVersion: "v1",
		Kind:       "Component",
		Metadata:   fabrica.Metadata{Name: "x3000c0s0b0n0", UID: "test-uid"},
		Spec: ComponentSpec{
			ID:      "x3000c0s0b0n0",
			Type:    "Node",
			State:   "On",
			Flag:    "OK",
			Enabled: &enabled,
			Role:    "Compute",
			NID:     json.Number("3"),
			Arch:    "X86",
			Class:   "River",
		},
	}
}

func TestComponentValidate_ValidPasses(t *testing.T) {
	c := validComponent()
	if err := c.Validate(context.Background()); err != nil {
		t.Fatalf("expected valid component to pass validation, got error: %v", err)
	}
}

func TestComponentValidate_MissingSpecID(t *testing.T) {
	c := validComponent()
	c.Spec.ID = ""
	if err := c.Validate(context.Background()); err == nil {
		t.Fatal("expected validation error for empty spec.ID, got nil")
	}
}

func TestComponentValidate_InvalidNIDLiteral(t *testing.T) {
	c := validComponent()
	// json.Number with a non-numeric literal produces invalid JSON when marshaled,
	// which should surface as a marshaling error from Validate.
	c.Spec.NID = json.Number("not-a-number")
	if err := c.Validate(context.Background()); err == nil {
		t.Fatal("expected error for invalid NID literal, got nil")
	}
}

func TestComponentValidate_EmptyNIDOmitted(t *testing.T) {
	c := validComponent()
	// Empty NID should be omitted (omitempty) and not cause validation to fail.
	c.Spec.NID = ""
	if err := c.Validate(context.Background()); err != nil {
		t.Fatalf("expected empty NID to be omitted without error, got: %v", err)
	}
}

func TestComponentValidate_MinimalSpecPasses(t *testing.T) {
	// Only the required field (spec.ID) is set; schema does not require
	// any other spec fields, nor top-level apiVersion/kind/metadata.
	c := &Component{Spec: ComponentSpec{ID: "x1000c0s0b0n0"}}
	if err := c.Validate(context.Background()); err != nil {
		t.Fatalf("expected minimal valid component to pass, got: %v", err)
	}
}

func TestComponentValidate_LongDescriptionNotEnforcedByJSONSchema(t *testing.T) {
	// The struct tag validate:"max=200" on Spec.Description is not enforced by
	// Validate(), since Validate() only checks the JSON schema (which has no
	// maxLength constraint on description). This test documents that behavior;
	// if the schema is later updated to enforce this, the test should be updated.
	c := validComponent()
	longDesc := make([]byte, 500)
	for i := range longDesc {
		longDesc[i] = 'a'
	}
	c.Spec.Description = string(longDesc)

	if err := c.Validate(context.Background()); err != nil {
		t.Fatalf("expected long description to pass since JSON schema has no maxLength, got: %v", err)
	}
}

func TestComponentGetters(t *testing.T) {
	c := &Component{
		Metadata: fabrica.Metadata{Name: "x3000c0s0b0n0", UID: "abc-123"},
	}

	if got := c.GetKind(); got != "Component" {
		t.Errorf("GetKind() = %q, want %q", got, "Component")
	}
	if got := c.GetName(); got != "x3000c0s0b0n0" {
		t.Errorf("GetName() = %q, want %q", got, "x3000c0s0b0n0")
	}
	if got := c.GetUID(); got != "abc-123" {
		t.Errorf("GetUID() = %q, want %q", got, "abc-123")
	}

	// IsHub is a marker method; it should simply not panic.
	c.IsHub()
}
