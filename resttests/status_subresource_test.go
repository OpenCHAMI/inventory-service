/*
 * Copyright © 2026 OpenCHAMI a Series of LF Projects, LLC
 *
 * SPDX-License-Identifier: MIT
 *
 * Tests for the status subresource routes registered in routes_generated.go.
 * Every generated resource exposes:
 *   PUT   /{resource}/{uid}/status
 *   PATCH /{resource}/{uid}/status
 *
 * All resource Status types (ComponentStatus, ComponentEndpointStatus, ...)
 * share the same {phase, message, ready} shape, so a single set of shared
 * types is used to exercise every resource's status subresource below.
 */

package resttests

import (
	"net/http"
	"testing"
)

// resourceStatus mirrors the shared *Status shape used by every generated resource.
type resourceStatus struct {
	Phase   string `json:"phase,omitempty"`
	Message string `json:"message,omitempty"`
	Ready   bool   `json:"ready"`
}

// statusResponse decodes just the status portion of any resource response.
type statusResponse struct {
	Status resourceStatus `json:"status"`
}

// checkStatusSubresource exercises PUT and PATCH on the given resource's
// status subresource path and verifies the update is applied and persisted.
func checkStatusSubresource(t *testing.T, statusPath, resourcePath string) {
	t.Helper()

	// PUT replaces the whole status.
	putResp := doRequest(t, http.MethodPut, statusPath, resourceStatus{Phase: "Ready", Message: "all good", Ready: true})
	requireStatus(t, putResp, http.StatusOK)
	var afterPut statusResponse
	decodeJSON(t, putResp, &afterPut)
	if afterPut.Status.Phase != "Ready" {
		t.Errorf("PUT %s: expected Phase=Ready, got %q", statusPath, afterPut.Status.Phase)
	}
	if !afterPut.Status.Ready {
		t.Errorf("PUT %s: expected Ready=true, got false", statusPath)
	}

	// PATCH (JSON Merge Patch, the default for Content-Type: application/json)
	// updates only the given fields and preserves the rest.
	patchResp := doRequest(t, http.MethodPatch, statusPath, map[string]any{"message": "patched"})
	requireStatus(t, patchResp, http.StatusOK)
	var afterPatch statusResponse
	decodeJSON(t, patchResp, &afterPatch)
	if afterPatch.Status.Message != "patched" {
		t.Errorf("PATCH %s: expected Message=patched, got %q", statusPath, afterPatch.Status.Message)
	}
	if afterPatch.Status.Phase != "Ready" {
		t.Errorf("PATCH %s: expected Phase to be preserved as Ready, got %q", statusPath, afterPatch.Status.Phase)
	}

	// Confirm the patched status persisted via a subsequent GET of the resource.
	getResp := doRequest(t, http.MethodGet, resourcePath, nil)
	requireStatus(t, getResp, http.StatusOK)
	var fetched statusResponse
	decodeJSON(t, getResp, &fetched)
	if fetched.Status.Message != "patched" {
		t.Errorf("GET %s after PATCH status: expected Message=patched, got %q", resourcePath, fetched.Status.Message)
	}
}

func TestComponentStatusSubresource(t *testing.T) {
	created := createAndRequire(t, newComponent("x9000c0s0b0n0", "Node"))
	defer func() { doRequest(t, http.MethodDelete, "/components/"+created.Metadata.UID, nil).Body.Close() }()

	resourcePath := "/components/" + created.Metadata.UID
	checkStatusSubresource(t, resourcePath+"/status", resourcePath)
}

func TestComponentEndpointStatusSubresource(t *testing.T) {
	created := createComponentEndpointAndRequire(t, newComponentEndpoint("x9000c0s0b0n1", "x9000c0s0b0"))
	defer func() { doRequest(t, http.MethodDelete, "/componentendpoints/"+created.Metadata.UID, nil).Body.Close() }()

	resourcePath := "/componentendpoints/" + created.Metadata.UID
	checkStatusSubresource(t, resourcePath+"/status", resourcePath)
}

func TestEthernetInterfaceStatusSubresource(t *testing.T) {
	created := createEthernetInterfaceAndRequire(t, newEthernetInterface("a0:00:00:00:00:09", "a0:00:00:00:00:09"))
	defer func() { doRequest(t, http.MethodDelete, "/ethernetinterfaces/"+created.Metadata.UID, nil).Body.Close() }()

	resourcePath := "/ethernetinterfaces/" + created.Metadata.UID
	checkStatusSubresource(t, resourcePath+"/status", resourcePath)
}

func TestGroupStatusSubresource(t *testing.T) {
	created := createGroupAndRequire(t, newGroup("status-test-group"))
	defer func() { doRequest(t, http.MethodDelete, "/groups/"+created.Metadata.UID, nil).Body.Close() }()

	resourcePath := "/groups/" + created.Metadata.UID
	checkStatusSubresource(t, resourcePath+"/status", resourcePath)
}

func TestHardwareStatusSubresource(t *testing.T) {
	created := createHardwareAndRequire(t, newHardware("x9000c0s9b0n0", "Node"))
	defer func() { doRequest(t, http.MethodDelete, "/hardwares/"+created.Metadata.UID, nil).Body.Close() }()

	resourcePath := "/hardwares/" + created.Metadata.UID
	checkStatusSubresource(t, resourcePath+"/status", resourcePath)
}

func TestRedfishEndpointStatusSubresource(t *testing.T) {
	created := createRedfishEndpointAndRequire(t, newRedfishEndpoint("x9000c0s9b0", "bmc9.example.com"))
	defer func() { doRequest(t, http.MethodDelete, "/redfishendpoints/"+created.Metadata.UID, nil).Body.Close() }()

	resourcePath := "/redfishendpoints/" + created.Metadata.UID
	checkStatusSubresource(t, resourcePath+"/status", resourcePath)
}

func TestServiceEndpointStatusSubresource(t *testing.T) {
	created := createServiceEndpointAndRequire(t, newServiceEndpoint("x9000c0b9", "Chassis"))
	defer func() { doRequest(t, http.MethodDelete, "/serviceendpoints/"+created.Metadata.UID, nil).Body.Close() }()

	resourcePath := "/serviceendpoints/" + created.Metadata.UID
	checkStatusSubresource(t, resourcePath+"/status", resourcePath)
}
