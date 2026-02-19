package v1

import (
	"context"

	"github.com/openchami/fabrica/pkg/fabrica"
)

type EthernetInterface struct {
	APIVersion string                  `json:"apiVersion"`
	Kind       string                  `json:"kind"`
	Metadata   fabrica.Metadata        `json:"metadata"`
	Spec       EthernetInterfaceSpec   `json:"spec" validate:"required"`
	Status     EthernetInterfaceStatus `json:"status,omitempty"`
}

type EthernetInterfaceSpec struct {
	Description string `json:"description,omitempty" validate:"max=200"`
	ID          string `json:"ID"`
	MACAddr     string `json:"MACAddress"`
	IPAddr      string `json:"IPAddress"`
	LastUpdate  string `json:"LastUpdate"`
	CompID      string `json:"ComponentID"`
	Type        string `json:"Type"`
}

type EthernetInterfaceStatus struct {
	Phase   string `json:"phase,omitempty"`
	Message string `json:"message,omitempty"`
	Ready   bool   `json:"ready"`
}

func (r *EthernetInterface) Validate(ctx context.Context) error {

	return nil
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
