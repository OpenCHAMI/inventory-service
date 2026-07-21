/*
 * Copyright © 2026 OpenCHAMI a Series of LF Projects, LLC
 *
 * SPDX-License-Identifier: MIT
 *
 * Tests for the /hsm/v2/memberships routes registered in csm_routes.go.
 * Routes under test:
 *   GET /hsm/v2/memberships
 *   GET /hsm/v2/memberships/{id}
 *
 * Notes on response shapes (from csm_membership.go):
 *   GET / returns: []*Membership { id, groupLabels, partitionName }
 *   GET /{id}    : Membership (single object)
 *
 * Memberships are derived: a Membership exists for every ComponentEndpoint,
 * and its groupLabels list every Group whose members.ids includes that
 * ComponentEndpoint's ID.
 */

package resttests

import (
	"fmt"
	"net/http"
	"testing"
)

const csmMembershipBase = "/hsm/v2/memberships"

type membership struct {
	ID          string   `json:"id"`
	GroupLabels []string `json:"groupLabels"`
}

// TestGetMembershipsCsm verifies GET /hsm/v2/memberships returns a membership
// entry for a ComponentEndpoint, listing the groups it belongs to.
func TestGetMembershipsCsm(t *testing.T) {
	ceID := "x9100c0s0b0n0"
	groupLabel := "csm-membership-group"

	csmCECreate(t, newCsmComponentEndpoint(ceID, "x9100c0s0b0"))
	defer csmCEDelete(t, ceID)

	csmGroupCreate(t, newCsmGroup(groupLabel, ceID))
	defer csmGroupDelete(t, groupLabel)

	resp := doRequest(t, http.MethodGet, csmMembershipBase, nil)
	requireStatus(t, resp, http.StatusOK)

	var list []membership
	decodeJSON(t, resp, &list)

	var found *membership
	for i := range list {
		if list[i].ID == ceID {
			found = &list[i]
			break
		}
	}
	if found == nil {
		t.Fatalf("expected membership for ID=%q in list, got %d entries", ceID, len(list))
	}
	if len(found.GroupLabels) != 1 || found.GroupLabels[0] != groupLabel {
		t.Errorf("expected GroupLabels=[%q], got %v", groupLabel, found.GroupLabels)
	}
}

// TestGetMembershipCsm verifies GET /hsm/v2/memberships/{id} returns the
// membership for a specific ComponentEndpoint ID.
func TestGetMembershipCsm(t *testing.T) {
	ceID := "x9100c0s0b0n1"
	groupLabel := "csm-membership-group-2"

	csmCECreate(t, newCsmComponentEndpoint(ceID, "x9100c0s0b0"))
	defer csmCEDelete(t, ceID)

	csmGroupCreate(t, newCsmGroup(groupLabel, ceID))
	defer csmGroupDelete(t, groupLabel)

	resp := doRequest(t, http.MethodGet, fmt.Sprintf("%s/%s", csmMembershipBase, ceID), nil)
	requireStatus(t, resp, http.StatusOK)

	var m membership
	decodeJSON(t, resp, &m)
	if m.ID != ceID {
		t.Errorf("expected ID=%q, got %q", ceID, m.ID)
	}
	if len(m.GroupLabels) != 1 || m.GroupLabels[0] != groupLabel {
		t.Errorf("expected GroupLabels=[%q], got %v", groupLabel, m.GroupLabels)
	}
}

// TestGetMembershipCsmNotFound verifies GET /hsm/v2/memberships/{id} returns
// HTTP 404 for an unknown ComponentEndpoint ID.
func TestGetMembershipCsmNotFound(t *testing.T) {
	resp := doRequest(t, http.MethodGet, fmt.Sprintf("%s/%s", csmMembershipBase, "does-not-exist"), nil)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("expected HTTP 404 for unknown membership ID, got %d", resp.StatusCode)
	}
}
