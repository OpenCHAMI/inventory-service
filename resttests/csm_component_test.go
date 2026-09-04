/*
 * Copyright © 2026 OpenCHAMI a Series of LF Projects, LLC
 *
 * SPDX-License-Identifier: MIT
 *
 * Tests for the /hsm/v2/State/Components routes registered in csm_routes.go.
 * Routes under test:
 *   GET    /hsm/v2/State/Components
 *   POST   /hsm/v2/State/Components
 *   GET    /hsm/v2/State/Components/{id}
 *   PUT    /hsm/v2/State/Components/{id}
 *   DELETE /hsm/v2/State/Components/{id}
 *
 * Notes on request/response shapes (from csm_component_handlers.go):
 *   POST  body   : ComponentArray  { "Components": [ <ComponentSpec>, ... ] }
 *   POST  returns: HTTP 201, no body
 *   GET / returns: ComponentArray  { "Components": [ <ComponentSpec>, ... ] }
 *   GET /{id}:    flat ComponentSpec (fields at top level, SMD-compatible)
 *   PUT  body   : ComponentSpec
 *   PUT  returns: updated ComponentSpec
 *   DELETE /{id}: DeleteResponse { message, uid }
 */

package resttests

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

const csmBase = "/hsm/v2/State/Components"

// ─── CSM-specific request / response shapes ───────────────────────────────────
// These mirror csm_models.go / component_types.go without importing package main.

type csmComponentSpec struct {
	ID       string `json:"ID"`
	Type     string `json:"Type,omitempty"`
	State    string `json:"State,omitempty"`
	Flag     string `json:"Flag,omitempty"`
	Enabled  *bool  `json:"Enabled,omitempty"`
	SwStatus string `json:"SoftwareStatus,omitempty"`
	Role     string `json:"Role,omitempty"`
	SubRole  string `json:"SubRole,omitempty"`
	Subtype  string `json:"Subtype,omitempty"`
	NetType  string `json:"NetType,omitempty"`
	Arch     string `json:"Arch,omitempty"`
	Class    string `json:"Class,omitempty"`
	NID      any    `json:"NID,omitempty"`
}

// csmComponentArray mirrors cmd/server.ComponentArray.
type csmComponentArray struct {
	Components []*csmComponentSpec `json:"Components"`
}

// ─── Helpers ──────────────────────────────────────────────────────────────────

// csmCreate POSTs a batch of components via the CSM endpoint and asserts 201.
func csmCreate(t *testing.T, specs ...*csmComponentSpec) {
	t.Helper()
	body := csmComponentArray{Components: specs}
	resp := doRequest(t, http.MethodPost, csmBase, body)
	requireStatus(t, resp, http.StatusCreated)
	resp.Body.Close()
}

// csmGetOne fetches a single component by xname ID and returns the flat SMD
// ComponentSpec (fields at top level, no fabrica metadata).
func csmGetOne(t *testing.T, xname string) (*csmComponentSpec, int) {
	t.Helper()
	resp := doRequest(t, http.MethodGet, fmt.Sprintf("%s/%s", csmBase, xname), nil)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, resp.StatusCode
	}
	var c csmComponentSpec
	decodeJSON(t, resp, &c)
	return &c, resp.StatusCode
}

// nativeComponentByXname finds a component in the native GET /components list by
// its xname (Spec.ID). The flat CSM response omits fabrica metadata (UID,
// timestamps), so tests that need to observe those read them from here.
func nativeComponentByXname(t *testing.T, xname string) *componentResponse {
	t.Helper()
	resp := doRequest(t, http.MethodGet, "/components", nil)
	requireStatus(t, resp, http.StatusOK)
	var list []componentResponse
	decodeJSON(t, resp, &list)
	for i := range list {
		if list[i].Spec.ID == xname {
			return &list[i]
		}
	}
	t.Fatalf("component %s not found in native /components list", xname)
	return nil
}

// csmDelete deletes a component by xname ID.
func csmDelete(t *testing.T, xname string) {
	t.Helper()
	resp := doRequest(t, http.MethodDelete, fmt.Sprintf("%s/%s", csmBase, xname), nil)
	resp.Body.Close()
}

// ─── Tests ────────────────────────────────────────────────────────────────────

// TestCreateComponentCsm verifies that POST /hsm/v2/State/Components with a
// ComponentArray body returns HTTP 201.
func TestCreateComponentCsm(t *testing.T) {
	xname := "x3000c0s2b0n0"
	csmCreate(t, &csmComponentSpec{ID: xname, Type: "Node"})
	defer csmDelete(t, xname)

	// Verify the component exists after creation
	comp, status := csmGetOne(t, xname)
	if status != http.StatusOK {
		t.Fatalf("expected HTTP 200 for GET after POST, got %d", status)
	}
	if comp.ID != xname {
		t.Errorf("expected ID=%q, got %q", xname, comp.ID)
	}
}

// TestCreateComponentCsmBulk verifies that multiple components can be created
// in a single POST.
func TestCreateComponentCsmBulk(t *testing.T) {
	xnames := []string{"x3000c0s2b0n1", "x3000c0s2b0n2", "x3000c0s2b0n3"}
	specs := make([]*csmComponentSpec, len(xnames))
	for i, x := range xnames {
		specs[i] = &csmComponentSpec{ID: x, Type: "Node"}
	}
	csmCreate(t, specs...)
	defer func() {
		for _, x := range xnames {
			csmDelete(t, x)
		}
	}()

	// Verify each was created
	for _, x := range xnames {
		comp, status := csmGetOne(t, x)
		if status != http.StatusOK {
			t.Errorf("expected HTTP 200 for %s, got %d", x, status)
			continue
		}
		if comp.ID != x {
			t.Errorf("expected ID=%q, got %q", x, comp.ID)
		}
	}
}

// TestGetComponentsCsm verifies that GET /hsm/v2/State/Components returns
// HTTP 200 and a ComponentArray with at least the created component.
func TestGetComponentsCsm(t *testing.T) {
	xname := "x3000c0s2b0n4"
	csmCreate(t, &csmComponentSpec{ID: xname, Type: "Node"})
	defer csmDelete(t, xname)

	resp := doRequest(t, http.MethodGet, csmBase, nil)
	requireStatus(t, resp, http.StatusOK)

	var list csmComponentArray
	decodeJSON(t, resp, &list)

	found := false
	for _, c := range list.Components {
		if c != nil && c.ID == xname {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("component %s not found in GET %s list", xname, csmBase)
	}
}

// TestGetComponentCsm verifies that GET /hsm/v2/State/Components/{id} returns
// HTTP 200 and the correct component.
func TestGetComponentCsm(t *testing.T) {
	xname := "x3000c0s2b0n5"
	csmCreate(t, &csmComponentSpec{ID: xname, Type: "Node"})
	defer csmDelete(t, xname)

	comp, status := csmGetOne(t, xname)
	if status != http.StatusOK {
		t.Fatalf("expected HTTP 200, got %d", status)
	}
	if comp.ID != xname {
		t.Errorf("expected ID=%q, got %q", xname, comp.ID)
	}
}

// TestUpdateComponentCsm verifies that PUT /hsm/v2/State/Components/{id}
// updates the component spec and returns HTTP 200.
func TestUpdateComponentCsm(t *testing.T) {
	xname := "x3000c0s2b0n6"
	csmCreate(t, &csmComponentSpec{ID: xname, Type: "Node"})
	defer csmDelete(t, xname)

	updateSpec := csmComponentSpec{
		ID:    xname,
		Type:  "Node",
		State: "Ready",
		Role:  "Compute",
		Flag:  "OK",
	}
	resp := doRequest(t, http.MethodPut, fmt.Sprintf("%s/%s", csmBase, xname), updateSpec)
	requireStatus(t, resp, http.StatusOK)
	resp.Body.Close()

	// Verify via GET that the update persisted
	comp, status := csmGetOne(t, xname)
	if status != http.StatusOK {
		t.Fatalf("expected HTTP 200 after PUT, got %d", status)
	}
	if comp.State != "Ready" {
		t.Errorf("expected State=Ready after PUT, got %q", comp.State)
	}
	if comp.Role != "Compute" {
		t.Errorf("expected Role=Compute after PUT, got %q", comp.Role)
	}
}

// TestDeleteComponentCsm verifies that DELETE /hsm/v2/State/Components/{id}
// returns HTTP 200 and that a subsequent GET does not return HTTP 200.
func TestDeleteComponentCsm(t *testing.T) {
	xname := "x3000c0s2b0n7"
	csmCreate(t, &csmComponentSpec{ID: xname, Type: "Node"})

	// Delete
	delResp := doRequest(t, http.MethodDelete, fmt.Sprintf("%s/%s", csmBase, xname), nil)
	requireStatus(t, delResp, http.StatusOK)

	var delBody deleteComponentResponse
	decodeJSON(t, delResp, &delBody)
	if delBody.Message == "" {
		t.Error("expected non-empty message in delete response")
	}

	// Confirm it is gone – the CSM GET handler returns non-200 for missing components
	_, status := csmGetOne(t, xname)
	if status == http.StatusOK {
		t.Errorf("expected non-200 after DELETE %s, got 200", xname)
	}
}

// TestCsmComponentLifecycle exercises the full POST → GET → PUT → DELETE cycle
// via the CSM /hsm/v2/State/Components endpoint.
func TestCsmComponentLifecycle(t *testing.T) {
	xname := "x3000c0s2b0n8"

	// ── POST ──────────────────────────────────────────────────────────────────
	csmCreate(t, &csmComponentSpec{ID: xname, Type: "Node"})

	// ── GET by xname ──────────────────────────────────────────────────────────
	comp, status := csmGetOne(t, xname)
	if status != http.StatusOK {
		t.Fatalf("POST→GET: expected HTTP 200, got %d", status)
	}
	if comp.ID != xname {
		t.Errorf("POST→GET: expected ID=%q, got %q", xname, comp.ID)
	}
	t.Logf("Created component UID: %s", nativeComponentByXname(t, xname).Metadata.UID)

	// ── GET all – component must appear ───────────────────────────────────────
	listResp := doRequest(t, http.MethodGet, csmBase, nil)
	requireStatus(t, listResp, http.StatusOK)
	var list csmComponentArray
	decodeJSON(t, listResp, &list)
	found := false
	for _, c := range list.Components {
		if c != nil && c.ID == xname {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("component %s not found in GET list after creation", xname)
	}

	// ── PUT ───────────────────────────────────────────────────────────────────
	putResp := doRequest(t, http.MethodPut, fmt.Sprintf("%s/%s", csmBase, xname), csmComponentSpec{
		ID: xname, Type: "Node", State: "On", Role: "Compute", Flag: "OK",
	})
	requireStatus(t, putResp, http.StatusOK)
	putResp.Body.Close()

	comp, status = csmGetOne(t, xname)
	if status != http.StatusOK {
		t.Fatalf("GET after PUT: expected HTTP 200, got %d", status)
	}
	if comp.State != "On" {
		t.Errorf("PUT: expected State=On, got %q", comp.State)
	}

	// ── DELETE ────────────────────────────────────────────────────────────────
	delResp := doRequest(t, http.MethodDelete, fmt.Sprintf("%s/%s", csmBase, xname), nil)
	requireStatus(t, delResp, http.StatusOK)
	delResp.Body.Close()

	_, status = csmGetOne(t, xname)
	if status == http.StatusOK {
		t.Errorf("expected non-200 after DELETE, still got 200")
	}
}

// TestCreateComponentCsmUpsert verifies that POST /hsm/v2/State/Components with an
// ID that already exists updates the component in place (SMD upsert semantics):
// the POST succeeds, the UID and CreatedAt are preserved, UpdatedAt advances, and
// the spec fields are replaced by the new body.
func TestCreateComponentCsmUpsert(t *testing.T) {
	xname := "x3000c0s3b0n0"
	csmCreate(t, &csmComponentSpec{ID: xname, Type: "Node", State: "On", Role: "Compute"})
	defer csmDelete(t, xname)

	// UID/CreatedAt/UpdatedAt are inventory-service internals that the flat CSM
	// response intentionally omits, so read them from the native /components list.
	before := nativeComponentByXname(t, xname)

	// Ensure a measurable gap so UpdatedAt is observably different.
	time.Sleep(10 * time.Millisecond)

	// Re-POST the same xname with different spec fields.
	resp := doRequest(t, http.MethodPost, csmBase, csmComponentArray{
		Components: []*csmComponentSpec{{ID: xname, Type: "Node", State: "Off", Role: "Service"}},
	})
	requireStatus(t, resp, http.StatusCreated)
	resp.Body.Close()

	after := nativeComponentByXname(t, xname)

	// UID and CreatedAt are stable across the upsert; UpdatedAt advances.
	if after.Metadata.UID != before.Metadata.UID {
		t.Errorf("expected UID to be preserved on upsert: before=%q after=%q",
			before.Metadata.UID, after.Metadata.UID)
	}
	if after.Metadata.CreatedAt != before.Metadata.CreatedAt {
		t.Errorf("expected CreatedAt to be preserved on upsert: before=%q after=%q",
			before.Metadata.CreatedAt, after.Metadata.CreatedAt)
	}
	if after.Metadata.UpdatedAt == before.Metadata.UpdatedAt {
		t.Errorf("expected UpdatedAt to change on upsert, still %q", after.Metadata.UpdatedAt)
	}

	// Spec fields are replaced by the new POST body (observed via the flat CSM GET).
	got, status := csmGetOne(t, xname)
	if status != http.StatusOK {
		t.Fatalf("expected HTTP 200 for CSM GET after upsert, got %d", status)
	}
	if got.State != "Off" {
		t.Errorf("expected State=Off after upsert, got %q", got.State)
	}
	if got.Role != "Service" {
		t.Errorf("expected Role=Service after upsert, got %q", got.Role)
	}
}

// TestCreateComponentCsmFieldMapping verifies that the CSM batch POST
// (POST /hsm/v2/State/Components) persists the full set of SMD component
// fields — not just ID/Type. A component is created with every mappable
// field populated and each is asserted to survive a POST → GET round-trip.
func TestCreateComponentCsmFieldMapping(t *testing.T) {
	xname := "x3000c0s4b0n0"
	enabled := true
	want := &csmComponentSpec{
		ID:       xname,
		Type:     "Node",
		State:    "On",
		Flag:     "OK",
		Enabled:  &enabled,
		SwStatus: "AdminUp",
		Role:     "Compute",
		SubRole:  "Worker",
		Subtype:  "river",
		NetType:  "Sling",
		Arch:     "X86",
		Class:    "River",
		NID:      42,
	}
	csmCreate(t, want)
	defer csmDelete(t, xname)

	got, status := csmGetOne(t, xname)
	if status != http.StatusOK {
		t.Fatalf("expected HTTP 200 for GET after POST, got %d", status)
	}

	if got.ID != want.ID {
		t.Errorf("ID: expected %q, got %q", want.ID, got.ID)
	}
	if got.Type != want.Type {
		t.Errorf("Type: expected %q, got %q", want.Type, got.Type)
	}
	if got.State != want.State {
		t.Errorf("State: expected %q, got %q", want.State, got.State)
	}
	if got.Flag != want.Flag {
		t.Errorf("Flag: expected %q, got %q", want.Flag, got.Flag)
	}
	if got.Enabled == nil || *got.Enabled != enabled {
		t.Errorf("Enabled: expected %v, got %v", enabled, got.Enabled)
	}
	if got.SwStatus != want.SwStatus {
		t.Errorf("SoftwareStatus: expected %q, got %q", want.SwStatus, got.SwStatus)
	}
	if got.Role != want.Role {
		t.Errorf("Role: expected %q, got %q", want.Role, got.Role)
	}
	if got.SubRole != want.SubRole {
		t.Errorf("SubRole: expected %q, got %q", want.SubRole, got.SubRole)
	}
	if got.Subtype != want.Subtype {
		t.Errorf("Subtype: expected %q, got %q", want.Subtype, got.Subtype)
	}
	if got.NetType != want.NetType {
		t.Errorf("NetType: expected %q, got %q", want.NetType, got.NetType)
	}
	if got.Arch != want.Arch {
		t.Errorf("Arch: expected %q, got %q", want.Arch, got.Arch)
	}
	if got.Class != want.Class {
		t.Errorf("Class: expected %q, got %q", want.Class, got.Class)
	}
	if fmt.Sprintf("%v", got.NID) != "42" {
		t.Errorf("NID: expected 42, got %v", got.NID)
	}
}
