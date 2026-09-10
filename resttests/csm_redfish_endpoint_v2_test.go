/*
 * Copyright © 2026 OpenCHAMI a Series of LF Projects, LLC
 *
 * SPDX-License-Identifier: MIT
 *
 * Tests for V2-format RedfishEndpoint discovery via POST /hsm/v2/Inventory/RedfishEndpoints.
 * A V2 body carries "Systems"/"Managers" inventory arrays and the handler derives
 * Component, ComponentEndpoint, and EthernetInterface sub-resources
 * (createV2SubResources in csm_redfish_endpoints.go).
 *
 * These tests focus on EthernetInterface duplicate-MAC handling. EthernetInterfaces
 * are keyed by MAC (colons stripped). Matching SMD, a MAC that already exists causes
 * a 409 Conflict during non-forced discovery (POST), while the forceUpdate path (PUT
 * on an existing endpoint) updates the interface in place, preserving its UID and
 * CreatedAt.
 */

package resttests

import (
	"fmt"
	"net/http"
	"testing"
)

// ─── V2 request shapes (mirror RedfishEndpointV2Request in csm_models.go) ──────

type csmV2Eth struct {
	URI         string `json:"uri,omitempty"`
	MAC         string `json:"mac,omitempty"`
	IP          string `json:"ip,omitempty"`
	Description string `json:"description,omitempty"`
}

type csmV2System struct {
	URI                string     `json:"uri,omitempty"`
	UUID               string     `json:"uuid,omitempty"`
	Name               string     `json:"name,omitempty"`
	SystemType         string     `json:"system_type,omitempty"`
	EthernetInterfaces []csmV2Eth `json:"ethernet_interfaces,omitempty"`
}

// csmV2RedfishEndpoint embeds the standard endpoint spec and adds a Systems array,
// producing the V2 discovery body the handler branches on.
type csmV2RedfishEndpoint struct {
	csmRedfishEndpointSpec
	Systems []csmV2System `json:"Systems,omitempty"`
}

// findNativeEIByID returns the native EthernetInterface resource (with metadata)
// whose Spec.ID matches id, or ok=false if none exists.
func findNativeEIByID(t *testing.T, id string) (ethernetInterfaceResponse, bool) {
	t.Helper()
	resp := doRequest(t, http.MethodGet, "/ethernetinterfaces", nil)
	requireStatus(t, resp, http.StatusOK)
	var list []ethernetInterfaceResponse
	decodeJSON(t, resp, &list)
	for _, ei := range list {
		if ei.Spec.ID == id {
			return ei, true
		}
	}
	return ethernetInterfaceResponse{}, false
}

// ─── Tests ────────────────────────────────────────────────────────────────────

// TestCreateRedfishEndpointCsmV2SharedMAC verifies that a V2 discovery body whose
// systems share a single MAC is rejected with 409, matching SMD: the duplicate MAC
// conflicts with the interface created for the first system, and non-forced discovery
// (POST) does not overwrite it.
func TestCreateRedfishEndpointCsmV2SharedMAC(t *testing.T) {
	const (
		reID  = "x9000c7s0b0"
		node0 = reID + "n0"
		node1 = reID + "n1"
		mac   = "de:ca:fc:0f:fe:01"
		macID = "decafc0ffe01"
	)

	body := csmV2RedfishEndpoint{
		csmRedfishEndpointSpec: newCsmRedfishEndpoint(reID, reID+".example.com"),
		Systems: []csmV2System{
			{
				URI:                "/redfish/v1/Systems/0",
				Name:               "System 0",
				SystemType:         "Physical",
				EthernetInterfaces: []csmV2Eth{{URI: "/redfish/v1/Systems/0/EthernetInterfaces/0", MAC: mac, IP: "10.0.0.10"}},
			},
			{
				URI:                "/redfish/v1/Systems/1",
				Name:               "System 1",
				SystemType:         "Physical",
				EthernetInterfaces: []csmV2Eth{{URI: "/redfish/v1/Systems/1/EthernetInterfaces/0", MAC: mac, IP: "10.0.0.11"}},
			},
		},
	}

	resp := doRequest(t, http.MethodPost, csmREBase, body)
	requireStatus(t, resp, http.StatusConflict)
	resp.Body.Close()

	defer func() {
		csmREDelete(t, reID)
		csmCEDelete(t, node0)
		csmCEDelete(t, node1)
		csmEIDelete(t, macID)
		csmDelete(t, node0)
		csmDelete(t, node1)
	}()
}

// TestCreateRedfishEndpointCsmV2EthernetInterfaceUpsertPreservesIdentity verifies
// SMD-parity duplicate-MAC handling: a second endpoint that re-declares an existing
// MAC via POST is rejected with 409, while re-running discovery for the SAME endpoint
// via PUT (SMD's forceUpdate path) updates the EthernetInterface in place, preserving
// its UID and CreatedAt.
func TestCreateRedfishEndpointCsmV2EthernetInterfaceUpsertPreservesIdentity(t *testing.T) {
	const (
		reA   = "x9000c7s1b0"
		reB   = "x9000c7s2b0"
		nodeA = reA + "n0"
		nodeB = reB + "n0"
		mac   = "de:ca:fc:0f:fe:02"
		macID = "decafc0ffe02"
	)

	makeBody := func(reID, uri string) csmV2RedfishEndpoint {
		return csmV2RedfishEndpoint{
			csmRedfishEndpointSpec: newCsmRedfishEndpoint(reID, reID+".example.com"),
			Systems: []csmV2System{
				{
					URI:                "/redfish/v1/Systems/0",
					Name:               "System 0",
					SystemType:         "Physical",
					EthernetInterfaces: []csmV2Eth{{URI: uri, MAC: mac, IP: "10.0.0.20"}},
				},
			},
		}
	}

	// First discovery via POST creates the interface owned by nodeA.
	resp := doRequest(t, http.MethodPost, csmREBase, makeBody(reA, "/redfish/v1/Systems/0/EthernetInterfaces/0"))
	requireStatus(t, resp, http.StatusCreated)
	resp.Body.Close()

	defer func() {
		csmREDelete(t, reA)
		csmREDelete(t, reB)
		csmCEDelete(t, nodeA)
		csmCEDelete(t, nodeB)
		csmEIDelete(t, macID)
		csmDelete(t, nodeA)
		csmDelete(t, nodeB)
	}()

	first, ok := findNativeEIByID(t, macID)
	if !ok {
		t.Fatalf("expected EthernetInterface %s to exist after first discovery", macID)
	}
	if first.Spec.ComponentID != nodeA {
		t.Fatalf("expected ComponentID=%q after first discovery, got %q", nodeA, first.Spec.ComponentID)
	}

	// A different endpoint that re-declares the same MAC via POST conflicts (409).
	resp = doRequest(t, http.MethodPost, csmREBase, makeBody(reB, "/redfish/v1/Systems/0/EthernetInterfaces/0"))
	requireStatus(t, resp, http.StatusConflict)
	resp.Body.Close()

	// Re-running discovery for the SAME endpoint via PUT uses the forceUpdate path and
	// updates the interface in place, preserving its identity.
	resp = doRequest(t, http.MethodPut, fmt.Sprintf("%s/%s", csmREBase, reA), makeBody(reA, "/redfish/v1/Systems/0/EthernetInterfaces/0"))
	requireStatus(t, resp, http.StatusOK)
	resp.Body.Close()

	second, ok := findNativeEIByID(t, macID)
	if !ok {
		t.Fatalf("expected EthernetInterface %s to still exist after PUT", macID)
	}
	if second.Metadata.UID != first.Metadata.UID {
		t.Errorf("expected UID preserved across forceUpdate: was %q, got %q", first.Metadata.UID, second.Metadata.UID)
	}
	if second.Metadata.CreatedAt != first.Metadata.CreatedAt {
		t.Errorf("expected CreatedAt preserved across forceUpdate: was %q, got %q", first.Metadata.CreatedAt, second.Metadata.CreatedAt)
	}
	if second.Spec.ComponentID != nodeA {
		t.Errorf("expected ComponentID to remain %q after PUT, got %q", nodeA, second.Spec.ComponentID)
	}
}
