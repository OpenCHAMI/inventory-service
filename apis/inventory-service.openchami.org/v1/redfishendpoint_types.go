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
	APIVersion string                `json:"apiVersion" yaml:"apiVersion"`
	Kind       string                `json:"kind" yaml:"kind"`
	Metadata   fabrica.Metadata      `json:"metadata" yaml:"metadata"`
	ID         string                `json:"id,omitempty" yaml:"id,omitempty"`
	Spec       RedfishEndpointSpec   `json:"spec" yaml:"spec" validate:"required"`
	Status     RedfishEndpointStatus `json:"status,omitempty" yaml:"status,omitempty"`
}

type RedfishEndpointSpec struct {
	Description string `json:"description,omitempty" yaml:"description,omitempty" validate:"max=200"`
	ID          string `json:"ID" yaml:"ID"`

	Type     string `json:"Type" yaml:"Type"`
	Name     string `json:"Name,omitempty" yaml:"Name,omitempty"`
	Hostname string `json:"Hostname" yaml:"Hostname"`
	Domain   string `json:"Domain" yaml:"Domain"`
	FQDN     string `json:"FQDN" yaml:"FQDN"`
	Enabled  bool   `json:"Enabled" yaml:"Enabled"`
	UUID     string `json:"UUID,omitempty" yaml:"UUID,omitempty"`
	User     string `json:"User" yaml:"User"`
	Password string `json:"Password" yaml:"Password"`

	UseSSDP     bool `json:"UseSSDP,omitempty" yaml:"UseSSDP,omitempty"`
	MACRequired bool `json:"MACRequired,omitempty" yaml:"MACRequired,omitempty"`

	MACAddr            string `json:"MACAddr,omitempty" yaml:"MACAddr,omitempty"`
	IPAddress          string `json:"IPAddress,omitempty" yaml:"IPAddress,omitempty"`
	RedsicoverOnUpdate bool   `json:"RediscoverOnUpdate" yaml:"RediscoverOnUpdate"`
	TemplateID         string `json:"TemplateID,omitempty" yaml:"TemplateID,omitempty"`

	DiscoveryInfo DiscoveryInfo `json:"DiscoveryInfo" yaml:"DiscoveryInfo"`
}

type RedfishEndpointStatus struct {
	Phase   string `json:"phase,omitempty" yaml:"phase,omitempty"`
	Message string `json:"message,omitempty" yaml:"message,omitempty"`
	Ready   bool   `json:"ready" yaml:"ready"`
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
	LastAttempt    string `json:"LastDiscoveryAttempt,omitempty" yaml:"LastDiscoveryAttempt,omitempty"`
	LastStatus     string `json:"LastDiscoveryStatus" yaml:"LastDiscoveryStatus"`
	RedfishVersion string `json:"RedfishVersion,omitempty" yaml:"RedfishVersion,omitempty"`
}
