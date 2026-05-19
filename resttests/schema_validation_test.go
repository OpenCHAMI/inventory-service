/*
 * Copyright © 2026 OpenCHAMI a Series of LF Projects, LLC
 *
 * SPDX-License-Identifier: MIT
 *
 * Tests that POST requests violating the schemas in schemas/default/*.json are rejected
 * with a 4xx response.  One block per resource type, covering:
 *   - missing top-level "spec" (required by all schemas)
 *   - missing required field inside spec
 *   - empty string for a field with minLength: 1
 *
 * Endpoints under test (non-CSM generic routes):
 *   POST /components          — schemas/default/components_schema.json
 *   POST /ethernetinterfaces  — schemas/default/ethernet_interface_schema.json
 *   POST /groups              — schemas/default/group_schema.json
 *   POST /hardwares           — schemas/default/hardware_schema.json
 *   POST /redfishendpoints    — schemas/default/redfish_endpoint_schema.json
 *   POST /serviceendpoints    — schemas/default/service_endpoint_schema.json
 */

package resttests

import (
	"net/http"
	"testing"
)

// assertSchemaError posts body to path and fails the test if the server returns
// a 2xx status, which would indicate the schema validation was not enforced.
func assertSchemaError(t *testing.T, path string, body interface{}) {
	t.Helper()
	resp := doRequest(t, http.MethodPost, path, body)
	defer resp.Body.Close()
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		t.Errorf("expected 4xx for schema-violating request to %s, got HTTP %d", path, resp.StatusCode)
	}
}

// ─── Component (/components) ──────────────────────────────────────────────────
// Schema: schemas/default/components_schema.json
// Constraints: "spec" required; spec.ID required, minLength 1.

// TestCreateComponent_SchemaViolation_MissingSpec verifies that omitting "spec"
// entirely is rejected (schema requires "spec" at the top level).
func TestCreateComponent_SchemaViolation_MissingSpec(t *testing.T) {
	assertSchemaError(t, "/components", map[string]interface{}{
		"metadata": map[string]interface{}{"name": "x3000c0s0b0n0"},
	})
}

// TestCreateComponent_SchemaViolation_MissingID verifies that omitting spec.ID
// is rejected (spec.ID is required by the schema).
func TestCreateComponent_SchemaViolation_MissingID(t *testing.T) {
	assertSchemaError(t, "/components", map[string]interface{}{
		"metadata": map[string]interface{}{"name": "x3000c0s0b0n0"},
		"spec": map[string]interface{}{
			"Type":  "Node",
			"State": "On",
		},
	})
}

// TestCreateComponent_SchemaViolation_EmptyID verifies that an empty spec.ID
// is rejected (spec.ID has minLength: 1).
func TestCreateComponent_SchemaViolation_EmptyID(t *testing.T) {
	assertSchemaError(t, "/components", map[string]interface{}{
		"metadata": map[string]interface{}{"name": ""},
		"spec": map[string]interface{}{
			"ID":   "",
			"Type": "Node",
		},
	})
}

// ─── EthernetInterface (/ethernetinterfaces) ──────────────────────────────────
// Schema: schemas/default/ethernet_interface_schema.json
// Constraints: "spec" required; spec.ID required, minLength 1.

// TestCreateEthernetInterface_SchemaViolation_MissingSpec verifies that omitting
// "spec" is rejected.
func TestCreateEthernetInterface_SchemaViolation_MissingSpec(t *testing.T) {
	assertSchemaError(t, "/ethernetinterfaces", map[string]interface{}{
		"metadata": map[string]interface{}{"name": "a1:00:00:00:00:01"},
	})
}

// TestCreateEthernetInterface_SchemaViolation_MissingID verifies that omitting
// spec.ID is rejected.
func TestCreateEthernetInterface_SchemaViolation_MissingID(t *testing.T) {
	assertSchemaError(t, "/ethernetinterfaces", map[string]interface{}{
		"metadata": map[string]interface{}{"name": "a1:00:00:00:00:01"},
		"spec": map[string]interface{}{
			"MACAddress": "a1:00:00:00:00:01",
		},
	})
}

// TestCreateEthernetInterface_SchemaViolation_EmptyID verifies that an empty
// spec.ID is rejected (minLength: 1).
func TestCreateEthernetInterface_SchemaViolation_EmptyID(t *testing.T) {
	assertSchemaError(t, "/ethernetinterfaces", map[string]interface{}{
		"metadata": map[string]interface{}{"name": ""},
		"spec": map[string]interface{}{
			"ID":         "",
			"MACAddress": "a1:00:00:00:00:01",
		},
	})
}

// ─── Group (/groups) ──────────────────────────────────────────────────────────
// Schema: schemas/default/group_schema.json
// Constraints: "spec" required; spec.label required, minLength 1.

// TestCreateGroup_SchemaViolation_MissingSpec verifies that omitting "spec"
// is rejected.
func TestCreateGroup_SchemaViolation_MissingSpec(t *testing.T) {
	assertSchemaError(t, "/groups", map[string]interface{}{
		"metadata": map[string]interface{}{"name": "schema-test-group"},
	})
}

// TestCreateGroup_SchemaViolation_MissingLabel verifies that omitting spec.label
// is rejected (spec.label is required by the schema).
func TestCreateGroup_SchemaViolation_MissingLabel(t *testing.T) {
	assertSchemaError(t, "/groups", map[string]interface{}{
		"metadata": map[string]interface{}{"name": "schema-test-group"},
		"spec": map[string]interface{}{
			"members": map[string]interface{}{"ids": []string{}},
		},
	})
}

// TestCreateGroup_SchemaViolation_EmptyLabel verifies that an empty spec.label
// is rejected (minLength: 1).
func TestCreateGroup_SchemaViolation_EmptyLabel(t *testing.T) {
	assertSchemaError(t, "/groups", map[string]interface{}{
		"metadata": map[string]interface{}{"name": ""},
		"spec": map[string]interface{}{
			"label":   "",
			"members": map[string]interface{}{"ids": []string{}},
		},
	})
}

// ─── Hardware (/hardwares) ────────────────────────────────────────────────────
// Schema: schemas/default/hardware_schema.json
// Constraints: "spec" required; spec.ID required, minLength 1.

// TestCreateHardware_SchemaViolation_MissingSpec verifies that omitting "spec"
// is rejected.
func TestCreateHardware_SchemaViolation_MissingSpec(t *testing.T) {
	assertSchemaError(t, "/hardwares", map[string]interface{}{
		"metadata": map[string]interface{}{"name": "x3000c0s5b0n0"},
	})
}

// TestCreateHardware_SchemaViolation_MissingID verifies that omitting spec.ID
// is rejected.
func TestCreateHardware_SchemaViolation_MissingID(t *testing.T) {
	assertSchemaError(t, "/hardwares", map[string]interface{}{
		"metadata": map[string]interface{}{"name": "x3000c0s5b0n0"},
		"spec": map[string]interface{}{
			"Type":   "Node",
			"Status": "Populated",
		},
	})
}

// TestCreateHardware_SchemaViolation_EmptyID verifies that an empty spec.ID
// is rejected (minLength: 1).
func TestCreateHardware_SchemaViolation_EmptyID(t *testing.T) {
	assertSchemaError(t, "/hardwares", map[string]interface{}{
		"metadata": map[string]interface{}{"name": ""},
		"spec": map[string]interface{}{
			"ID":   "",
			"Type": "Node",
		},
	})
}

// ─── RedfishEndpoint (/redfishendpoints) ──────────────────────────────────────
// Schema: schemas/default/redfish_endpoint_schema.json
// Constraints: "spec" required; spec.ID required, minLength 1.

// TestCreateRedfishEndpoint_SchemaViolation_MissingSpec verifies that omitting
// "spec" is rejected.
func TestCreateRedfishEndpoint_SchemaViolation_MissingSpec(t *testing.T) {
	assertSchemaError(t, "/redfishendpoints", map[string]interface{}{
		"metadata": map[string]interface{}{"name": "x3000c0s0b0"},
	})
}

// TestCreateRedfishEndpoint_SchemaViolation_MissingID verifies that omitting
// spec.ID is rejected.
func TestCreateRedfishEndpoint_SchemaViolation_MissingID(t *testing.T) {
	assertSchemaError(t, "/redfishendpoints", map[string]interface{}{
		"metadata": map[string]interface{}{"name": "x3000c0s0b0"},
		"spec": map[string]interface{}{
			"Type":     "NodeBMC",
			"Hostname": "bmc.example.com",
		},
	})
}

// TestCreateRedfishEndpoint_SchemaViolation_EmptyID verifies that an empty
// spec.ID is rejected (minLength: 1).
func TestCreateRedfishEndpoint_SchemaViolation_EmptyID(t *testing.T) {
	assertSchemaError(t, "/redfishendpoints", map[string]interface{}{
		"metadata": map[string]interface{}{"name": ""},
		"spec": map[string]interface{}{
			"ID":       "",
			"Type":     "NodeBMC",
			"Hostname": "bmc.example.com",
		},
	})
}

// ─── ServiceEndpoint (/serviceendpoints) ─────────────────────────────────────
// Schema: schemas/default/service_endpoint_schema.json
// Constraints: "spec" required; spec.RedfishEndpointID and spec.RedfishType
// are both required, each with minLength 1.

// TestCreateServiceEndpoint_SchemaViolation_MissingSpec verifies that omitting
// "spec" is rejected.
func TestCreateServiceEndpoint_SchemaViolation_MissingSpec(t *testing.T) {
	assertSchemaError(t, "/serviceendpoints", map[string]interface{}{
		"metadata": map[string]interface{}{"name": "x3000c0b0"},
	})
}

// TestCreateServiceEndpoint_SchemaViolation_MissingRedfishEndpointID verifies
// that omitting spec.RedfishEndpointID is rejected.
func TestCreateServiceEndpoint_SchemaViolation_MissingRedfishEndpointID(t *testing.T) {
	assertSchemaError(t, "/serviceendpoints", map[string]interface{}{
		"metadata": map[string]interface{}{"name": "x3000c0b0"},
		"spec": map[string]interface{}{
			"RedfishType": "Chassis",
		},
	})
}

// TestCreateServiceEndpoint_SchemaViolation_MissingRedfishType verifies that
// omitting spec.RedfishType is rejected.
func TestCreateServiceEndpoint_SchemaViolation_MissingRedfishType(t *testing.T) {
	assertSchemaError(t, "/serviceendpoints", map[string]interface{}{
		"metadata": map[string]interface{}{"name": "x3000c0b0"},
		"spec": map[string]interface{}{
			"RedfishEndpointID": "x3000c0b0",
		},
	})
}

// TestCreateServiceEndpoint_SchemaViolation_EmptyRedfishEndpointID verifies
// that an empty spec.RedfishEndpointID is rejected (minLength: 1).
func TestCreateServiceEndpoint_SchemaViolation_EmptyRedfishEndpointID(t *testing.T) {
	assertSchemaError(t, "/serviceendpoints", map[string]interface{}{
		"metadata": map[string]interface{}{"name": ""},
		"spec": map[string]interface{}{
			"RedfishEndpointID": "",
			"RedfishType":       "Chassis",
		},
	})
}

// TestCreateServiceEndpoint_SchemaViolation_EmptyRedfishType verifies that
// an empty spec.RedfishType is rejected (minLength: 1).
func TestCreateServiceEndpoint_SchemaViolation_EmptyRedfishType(t *testing.T) {
	assertSchemaError(t, "/serviceendpoints", map[string]interface{}{
		"metadata": map[string]interface{}{"name": "x3000c0b0"},
		"spec": map[string]interface{}{
			"RedfishEndpointID": "x3000c0b0",
			"RedfishType":       "",
		},
	})
}
