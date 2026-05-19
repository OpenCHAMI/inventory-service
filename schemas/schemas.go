package schemas

import _ "embed"

// ComponentsSchema holds the JSON Schema for Component.
//
//go:embed default/components_schema.json
var ComponentsSchema []byte

// EthernetInterfaceSchema holds the JSON Schema for EthernetInterface.
//
//go:embed default/ethernet_interface_schema.json
var EthernetInterfaceSchema []byte

// GroupSchema holds the JSON Schema for Group.
//
//go:embed default/group_schema.json
var GroupSchema []byte

// HardwareSchema holds the JSON Schema for Hardware.
//
//go:embed default/hardware_schema.json
var HardwareSchema []byte

// RedfishEndpointSchema holds the JSON Schema for RedfishEndpoint.
//
//go:embed default/redfish_endpoint_schema.json
var RedfishEndpointSchema []byte

// ServiceEndpointSchema holds the JSON Schema for ServiceEndpoint.
//
//go:embed default/service_endpoint_schema.json
var ServiceEndpointSchema []byte
