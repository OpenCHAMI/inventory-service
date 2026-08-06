package v1

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/openchami/fabrica/pkg/fabrica"
	"github.com/openchami/inventory-service/schemas"
)

type EthernetInterface struct {
	APIVersion string                  `json:"apiVersion" yaml:"apiVersion"`
	Kind       string                  `json:"kind" yaml:"kind"`
	Metadata   fabrica.Metadata        `json:"metadata" yaml:"metadata"`
	ID         string                  `json:"id,omitempty" yaml:"id,omitempty"`
	Spec       EthernetInterfaceSpec   `json:"spec" yaml:"spec" validate:"required"`
	Status     EthernetInterfaceStatus `json:"status,omitempty" yaml:"status,omitempty"`
}

type EthernetInterfaceSpec struct {
	Description string      `json:"Description,omitempty" yaml:"Description,omitempty" validate:"max=200"`
	description string      `json:"description,omitempty" yaml:"description,omitempty" validate:"max=200"`
	ID          string      `json:"ID" yaml:"ID"`
	MACAddr     string      `json:"MACAddress" yaml:"MACAddress"`
	LastUpdate  string      `json:"LastUpdate" yaml:"LastUpdate"`
	CompID      string      `json:"ComponentID" yaml:"ComponentID"`
	Type        string      `json:"Type" yaml:"Type"`
	IPAddresses []IPAddress `json:"IPAddresses" yaml:"IPAddresses"`
}

type EthernetInterfaceStatus struct {
	Phase   string `json:"phase,omitempty" yaml:"phase,omitempty"`
	Message string `json:"message,omitempty" yaml:"message,omitempty"`
	Ready   bool   `json:"ready" yaml:"ready"`
}

func (r *EthernetInterface) Validate(ctx context.Context) error {
	var schema jsonschema.Schema
	if err := json.Unmarshal(schemas.EthernetInterfaceSchema, &schema); err != nil {
		return fmt.Errorf("loading ethernet interface schema: %w", err)
	}

	resolved, err := schema.Resolve(nil)
	if err != nil {
		return fmt.Errorf("resolving ethernet interface schema: %w", err)
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

func (r *EthernetInterface) GetKind() string {
	return "EthernetInterface"
}

func (r *EthernetInterface) GetName() string {
	return r.Metadata.Name
}

func (r *EthernetInterface) GetUID() string {
	return r.Metadata.UID
}

func (r *EthernetInterface) IsHub() {}

type IPAddress struct {
	IPAddress string `json:"IPAddress" yaml:"IPAddress"`
	Network   string `json:"Network,omitempty" yaml:"Network,omitempty"`
}
