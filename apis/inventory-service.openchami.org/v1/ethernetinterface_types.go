package v1

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/OpenCHAMI/inventory-service/schemas"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/openchami/fabrica/pkg/fabrica"
)

type EthernetInterface struct {
	APIVersion string                  `json:"apiVersion"`
	Kind       string                  `json:"kind"`
	Metadata   fabrica.Metadata        `json:"metadata"`
	ID         string                  `json:"id,omitempty"`
	Spec       EthernetInterfaceSpec   `json:"spec" validate:"required"`
	Status     EthernetInterfaceStatus `json:"status,omitempty"`
}

type EthernetInterfaceSpec struct {
	// todo SMD uses Description and fabrica uses description
	// One solution might be to maintain Description == description
	// and zero out one of them depending on if the request is through the
	// fabrica api or through the smd style api
	Description string      `json:"Description,omitempty" validate:"max=200"`
	description string      `json:"description,omitempty" validate:"max=200"`
	ID          string      `json:"ID"`
	MACAddr     string      `json:"MACAddress"`
	LastUpdate  string      `json:"LastUpdate"`
	CompID      string      `json:"ComponentID"`
	Type        string      `json:"Type"`
	IPAddresses []IPAddress `json:"IPAddresses"`
}

type EthernetInterfaceStatus struct {
	Phase   string `json:"phase,omitempty"`
	Message string `json:"message,omitempty"`
	Ready   bool   `json:"ready"`
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
	IPAddress string `json:"IPAddress"`
	Network   string `json:"Network,omitempty"`
}
