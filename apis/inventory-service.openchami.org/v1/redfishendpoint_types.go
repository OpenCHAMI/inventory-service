package v1

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/OpenCHAMI/inventory-service/schemas"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/openchami/fabrica/pkg/fabrica"
)

type RedfishEndpoint struct {
	APIVersion string                `json:"apiVersion"`
	Kind       string                `json:"kind"`
	Metadata   fabrica.Metadata      `json:"metadata"`
	ID         string                `json:"id,omitempty"`
	Spec       RedfishEndpointSpec   `json:"spec" validate:"required"`
	Status     RedfishEndpointStatus `json:"status,omitempty"`
}

type RedfishEndpointSpec struct {
	Description string `json:"description,omitempty" validate:"max=200"`
	ID          string `json:"ID"`

	Type     string `json:"Type"`
	Name     string `json:"Name,omitempty"`
	Hostname string `json:"Hostname"`
	Domain   string `json:"Domain"`
	FQDN     string `json:"FQDN"`
	Enabled  bool   `json:"Enabled"`
	UUID     string `json:"UUID,omitempty"`
	User     string `json:"User"`
	Password string `json:"Password"`

	UseSSDP     bool `json:"UseSSDP,omitempty"`
	MACRequired bool `json:"MACRequired,omitempty"`

	MACAddr            string `json:"MACAddr,omitempty"`
	IPAddress          string `json:"IPAddress,omitempty"`
	RedsicoverOnUpdate bool   `json:"RediscoverOnUpdate"`
	TemplateID         string `json:"TemplateID,omitempty"`

	DiscoveryInfo DiscoveryInfo `json:"DiscoveryInfo"`
}

type RedfishEndpointStatus struct {
	Phase   string `json:"phase,omitempty"`
	Message string `json:"message,omitempty"`
	Ready   bool   `json:"ready"`
}

func (r *RedfishEndpoint) Validate(ctx context.Context) error {
	var schema jsonschema.Schema
	if err := json.Unmarshal(schemas.RedfishEndpointSchema, &schema); err != nil {
		return fmt.Errorf("loading redfish endpoint schema: %w", err)
	}

	resolved, err := schema.Resolve(nil)
	if err != nil {
		return fmt.Errorf("resolving redfish endpoint schema: %w", err)
	}

	resourceJSON, err := json.Marshal(r)
	if err != nil {
		return fmt.Errorf("marshaling resource for validation: %w", err)
	}

	var instance any
	if err := json.Unmarshal(resourceJSON, &instance); err != nil {
		return fmt.Errorf("unmarshaling resource for validation: %w", err)
	}

	return resolved.Validate(instance)
}

func (r *RedfishEndpoint) GetKind() string {
	return "RedfishEndpoint"
}

func (r *RedfishEndpoint) GetName() string {
	return r.Metadata.Name
}

func (r *RedfishEndpoint) GetUID() string {
	return r.Metadata.UID
}

func (r *RedfishEndpoint) IsHub() {}

type DiscoveryInfo struct {
	LastAttempt    string `json:"LastDiscoveryAttempt,omitempty"`
	LastStatus     string `json:"LastDiscoveryStatus"`
	RedfishVersion string `json:"RedfishVersion,omitempty"`
}
