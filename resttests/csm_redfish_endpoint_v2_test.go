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
 * These tests focus on EthernetInterface upsert-by-MAC: EthernetInterfaces are keyed
 * by MAC (colons stripped), so a MAC shared across systems/endpoints must UPDATE the
 * existing resource in place — preserving its UID and CreatedAt — rather than abort
 * discovery on the unique resource ID. This mirrors SMD's upsert behaviour.
 */

package resttests

import (
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
// systems share a single MAC does not abort sub-resource creation: every node's
// ComponentEndpoint is created and the shared EthernetInterface is upserted to the
// last system that referenced it (last-writer-wins) instead of failing on the
// unique resource ID.
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
	requireStatus(t, resp, http.StatusCreated)
	resp.Body.Close()

	defer func() {
		csmREDelete(t, reID)
		csmCEDelete(t, node0)
		csmCEDelete(t, node1)
		csmEIDelete(t, macID)
		csmDelete(t, node0)
		csmDelete(t, node1)
	}()

	// Both nodes' ComponentEndpoints must exist — discovery did not abort.
	if _, st := csmCEGetOne(t, node0); st != http.StatusOK {
		t.Errorf("expected ComponentEndpoint %s to exist (HTTP 200), got %d", node0, st)
	}
	if _, st := csmCEGetOne(t, node1); st != http.StatusOK {
		t.Errorf("expected ComponentEndpoint %s to exist (HTTP 200), got %d", node1, st)
	}

	// The shared EthernetInterface exists once, owned by the last system.
	spec, st := csmEIGetOne(t, macID)
	if st != http.StatusOK {
		t.Fatalf("expected EthernetInterface %s to exist (HTTP 200), got %d", macID, st)
	}
	if spec.ComponentID != node1 {
		t.Errorf("expected shared EthernetInterface ComponentID=%q (last writer), got %q", node1, spec.ComponentID)
	}
}

// TestCreateRedfishEndpointCsmV2EthernetInterfaceUpsertPreservesIdentity verifies
// that re-discovering the same MAC under a different endpoint updates the existing
// EthernetInterface in place: its UID and CreatedAt are preserved while its
// ComponentID is reassigned to the newly-discovered node.
func TestCreateRedfishEndpointCsmV2EthernetInterfaceUpsertPreservesIdentity(t *testing.T) {
	const (
		reA   = "x9000c7s1b0"
		reB   = "x9000c7s2b0"
		nodeA = reA + "n0"
		nodeB = reB + "n0"
		mac   = "de:ca:fc:0f:fe:02"
		macID = "decafc0ffe02"
	)

	discover := func(reID, uri string) {
		body := csmV2RedfishEndpoint{
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
		resp := doRequest(t, http.MethodPost, csmREBase, body)
		requireStatus(t, resp, http.StatusCreated)
		resp.Body.Close()
	}

	discover(reA, "/redfish/v1/Systems/0/EthernetInterfaces/0")
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

	// Re-discover the same MAC under a different endpoint.
	discover(reB, "/redfish/v1/Systems/0/EthernetInterfaces/0")

	second, ok := findNativeEIByID(t, macID)
	if !ok {
		t.Fatalf("expected EthernetInterface %s to still exist after second discovery", macID)
	}
	if second.Metadata.UID != first.Metadata.UID {
		t.Errorf("expected UID preserved across upsert: was %q, got %q", first.Metadata.UID, second.Metadata.UID)
	}
	if second.Metadata.CreatedAt != first.Metadata.CreatedAt {
		t.Errorf("expected CreatedAt preserved across upsert: was %q, got %q", first.Metadata.CreatedAt, second.Metadata.CreatedAt)
	}
	if second.Spec.ComponentID != nodeB {
		t.Errorf("expected ComponentID reassigned to %q after re-discovery, got %q", nodeB, second.Spec.ComponentID)
	}
}
