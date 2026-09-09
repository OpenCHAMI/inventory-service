/*
 * Copyright © 2026 OpenCHAMI a Series of LF Projects, LLC
 *
 * SPDX-License-Identifier: MIT
 *
 * Tests for the SMD-compatible component query routes registered in
 * csm_routes.go:
 *   POST /hsm/v2/State/Components/Query
 *   GET  /hsm/v2/State/Components/Query/{xname}
 *
 * Request/response shapes (from csm_component_handlers.go / csm_models.go):
 *   POST body   : ComponentQuery { "ComponentIDs": [...], "type": [...], ... }
 *   POST returns: HTTP 200, ComponentArray { "Components": [ <ComponentSpec>, ... ] }
 *   GET  returns: HTTP 200, ComponentArray with the single matching component
 *
 * Power-control depends on POST /Query to build its component map; a missing
 * or empty response there is what caused integration-sandbox UC6 to fail.
 */

package resttests

import (
	"net/http"
	"testing"
)

// csmComponentQuery mirrors cmd/server.ComponentQuery.
type csmComponentQuery struct {
	ComponentIDs []string `json:"ComponentIDs,omitempty"`
	Type         []string `json:"type,omitempty"`
	State        []string `json:"state,omitempty"`
}

const csmQueryBase = csmBase + "/Query"

// TestQueryComponentsCsmByID verifies POST /Query with an explicit ComponentIDs
// filter returns exactly the requested component.
func TestQueryComponentsCsmByID(t *testing.T) {
	xname := "x9000c0s0b0n0"
	csmCreate(t, &csmComponentSpec{ID: xname, Type: "Node"})
	defer csmDelete(t, xname)

	resp := doRequest(t, http.MethodPost, csmQueryBase, csmComponentQuery{ComponentIDs: []string{xname}})
	requireStatus(t, resp, http.StatusOK)

	var list csmComponentArray
	decodeJSON(t, resp, &list)

	if len(list.Components) != 1 {
		t.Fatalf("expected exactly 1 component, got %d", len(list.Components))
	}
	if list.Components[0].ID != xname {
		t.Errorf("expected ID=%q, got %q", xname, list.Components[0].ID)
	}
}

// TestQueryComponentsCsmEmptyBody verifies POST /Query with no body matches all
// components (SMD parity), including the one just created.
func TestQueryComponentsCsmEmptyBody(t *testing.T) {
	xname := "x9000c0s0b0n1"
	csmCreate(t, &csmComponentSpec{ID: xname, Type: "Node"})
	defer csmDelete(t, xname)

	resp := doRequest(t, http.MethodPost, csmQueryBase, nil)
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
		t.Errorf("component %s not found in empty-query result", xname)
	}
}

// TestQueryComponentsCsmAllWildcard verifies POST /Query with the SMD wildcard
// ComponentIDs=["all"] returns every component (this is exactly what
// power-control sends via FillHSMData(["all"])).
func TestQueryComponentsCsmAllWildcard(t *testing.T) {
	xname := "x9000c0s0b0n3"
	csmCreate(t, &csmComponentSpec{ID: xname, Type: "Node"})
	defer csmDelete(t, xname)

	resp := doRequest(t, http.MethodPost, csmQueryBase, csmComponentQuery{ComponentIDs: []string{"all"}})
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
		t.Errorf("component %s not found in all-wildcard query result", xname)
	}
}

// TestQueryComponentsCsmUnknownXname verifies POST /Query for an xname that does
// not exist returns HTTP 200 with an empty (non-error) result.
func TestQueryComponentsCsmUnknownXname(t *testing.T) {
	resp := doRequest(t, http.MethodPost, csmQueryBase, csmComponentQuery{ComponentIDs: []string{"x9999c9s9b9n9"}})
	requireStatus(t, resp, http.StatusOK)

	var list csmComponentArray
	decodeJSON(t, resp, &list)

	if len(list.Components) != 0 {
		t.Errorf("expected 0 components for unknown xname, got %d", len(list.Components))
	}
}

// TestQueryComponentByXnameCsm verifies GET /Query/{xname} returns a
// ComponentArray with the single matching component.
func TestQueryComponentByXnameCsm(t *testing.T) {
	xname := "x9000c0s0b0n2"
	csmCreate(t, &csmComponentSpec{ID: xname, Type: "Node"})
	defer csmDelete(t, xname)

	resp := doRequest(t, http.MethodGet, csmQueryBase+"/"+xname, nil)
	requireStatus(t, resp, http.StatusOK)

	var list csmComponentArray
	decodeJSON(t, resp, &list)

	if len(list.Components) != 1 {
		t.Fatalf("expected exactly 1 component, got %d", len(list.Components))
	}
	if list.Components[0].ID != xname {
		t.Errorf("expected ID=%q, got %q", xname, list.Components[0].ID)
	}
}
