package v1

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/openchami/fabrica/pkg/fabrica"
	"github.com/openchami/inventory-service/schemas"
)

type ComponentEndpoint struct {
	APIVersion string                  `json:"apiVersion" yaml:"apiVersion"`
	Kind       string                  `json:"kind" yaml:"kind"`
	Metadata   fabrica.Metadata        `json:"metadata" yaml:"metadata"`
	ID         string                  `json:"id,omitempty" yaml:"id,omitempty"`
	Spec       ComponentEndpointSpec   `json:"spec" yaml:"spec" validate:"required"`
	Status     ComponentEndpointStatus `json:"status,omitempty" yaml:"status,omitempty"`
}

type ComponentEndpointSpec struct {
	Description string `json:"description,omitempty" yaml:"description,omitempty" validate:"max=200"`
	ID          string `json:"ID" yaml:"ID"`

	Type   string `json:"Type" yaml:"Type"`
	Domain string `json:"Domain,omitempty" yaml:"Domain,omitempty"`

	FQDN                string `json:"FQDN,omitempty" yaml:"FQDN,omitempty"`
	RedfishType         string `json:"RedfishType" yaml:"RedfishType"`
	RedfishSubtype      string `json:"RedfishSubtype" yaml:"RedfishSubtype"`
	MACAddr             string `json:"MACAddr,omitempty" yaml:"MACAddr,omitempty"`
	UUID                string `json:"UUID,omitempty" yaml:"UUID,omitempty"`
	OdataID             string `json:"OdataID" yaml:"OdataID"`
	RfEndpointID        string `json:"RedfishEndpointID" yaml:"RedfishEndpointID"`
	Enabled             bool   `json:"Enabled" yaml:"Enabled"`
	RedfishEndpointFQDN string `json:"RedfishEndpointFQDN,omitempty" yaml:"RedfishEndpointFQDN,omitempty"`
	URL                 string `json:"RedfishURL,omitempty" yaml:"RedfishURL,omitempty"`

	ComponentEndpointType string `json:"ComponentEndpointType" yaml:"ComponentEndpointType"`

	RedfishChassisInfo *ComponentChassisInfo `json:"RedfishChassisInfo,omitempty" yaml:"RedfishChassisInfo,omitempty"`
	RedfishSystemInfo  *ComponentSystemInfo  `json:"RedfishSystemInfo,omitempty" yaml:"RedfishSystemInfo,omitempty"`
	RedfishManagerInfo *ComponentManagerInfo `json:"RedfishManagerInfo,omitempty" yaml:"RedfishManagerInfo,omitempty"`
	RedfishPDUInfo     *ComponentPDUInfo     `json:"RedfishPDUInfo,omitempty" yaml:"RedfishPDUInfo,omitempty"`

	RedfishOutletInfo any `json:"RedfishOutletInfo,omitempty" yaml:"RedfishOutletInfo,omitempty"`
}

type ComponentEndpointStatus struct {
	Phase   string `json:"phase,omitempty" yaml:"phase,omitempty"`
	Message string `json:"message,omitempty" yaml:"message,omitempty"`
	Ready   bool   `json:"ready" yaml:"ready"`
}

func (r *ComponentEndpoint) Validate(ctx context.Context) error {
	var schema jsonschema.Schema
	if err := json.Unmarshal(schemas.ComponentEndpointSchema, &schema); err != nil {
		return fmt.Errorf("loading component endpoint schema: %w", err)
	}

	resolved, err := schema.Resolve(nil)
	if err != nil {
		return fmt.Errorf("resolving component endpoint schema: %w", err)
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

func (r *ComponentEndpoint) GetKind() string {
	return "ComponentEndpoint"
}

func (r *ComponentEndpoint) GetName() string {
	return r.Metadata.Name
}

func (r *ComponentEndpoint) GetUID() string {
	return r.Metadata.UID
}

func (r *ComponentEndpoint) IsHub() {}

type ResourceID struct {
	Oid string `json:"@odata.id" yaml:"@odata.id"`
}
type ComponentChassisInfo struct {
	Name    string          `json:"Name,omitempty" yaml:"Name,omitempty"`
	Actions *ChassisActions `json:"Actions,omitempty" yaml:"Actions,omitempty"`
}
type ComponentSystemInfo struct {
	Name    string                 `json:"Name,omitempty" yaml:"Name,omitempty"`
	Actions *ComputerSystemActions `json:"Actions,omitempty" yaml:"Actions,omitempty"`

	EthNICInfo []*EthernetNICInfo `json:"EthernetNICInfo,omitempty" yaml:"EthernetNICInfo,omitempty"`
	PowerCtlInfo
	Controls      []*Control     `json:"Controls,omitempty" yaml:"Controls,omitempty"`
	SerialConsole *SerialConsole `json:"SerialConsole,omitempty" yaml:"SerialConsole,omitempty"`
}
type PowerCtlInfo struct {
	PowerURL string          `json:"PowerURL,omitempty" yaml:"PowerURL,omitempty"`
	PowerCtl []*PowerControl `json:"PowerControl,omitempty" yaml:"PowerControl,omitempty"`
}
type PowerControl struct {
	ResourceID
	MemberId           string        `json:"MemberId,omitempty" yaml:"MemberId,omitempty"`
	Name               string        `json:"Name,omitempty" yaml:"Name,omitempty"`
	PowerCapacityWatts int           `json:"PowerCapacityWatts,omitempty" yaml:"PowerCapacityWatts,omitempty"`
	PowerConsumedWatts interface{}   `json:"PowerConsumedWatts,omitempty" yaml:"PowerConsumedWatts,omitempty"`
	OEM                *PwrCtlOEM    `json:"OEM,omitempty" yaml:"OEM,omitempty"`
	RelatedItem        []*ResourceID `json:"RelatedItem,omitempty" yaml:"RelatedItem,omitempty"`
}
type PwrCtlOEM struct {
	Cray *PwrCtlOEMCray `json:"Cray,omitempty" yaml:"Cray,omitempty"`
	HPE  *PwrCtlOEMHPE  `json:"HPE,omitempty" yaml:"HPE,omitempty"`
}
type PwrCtlOEMCray struct {
	PowerIdleWatts  int           `json:"PowerIdleWatts,omitempty" yaml:"PowerIdleWatts,omitempty"`
	PowerLimit      *CrayPwrLimit `json:"PowerLimit,omitempty" yaml:"PowerLimit,omitempty"`
	PowerResetWatts int           `json:"PowerResetWatts,omitempty" yaml:"PowerResetWatts,omitempty"`
}
type CrayPwrLimit struct {
	Min int `json:"Min,omitempty" yaml:"Min,omitempty"`
	Max int `json:"Max,omitempty" yaml:"Max,omitempty"`
}
type PwrCtlOEMHPE struct {
	PowerLimit             CrayPwrLimit `json:"PowerLimit" yaml:"PowerLimit"`
	PowerRegulationEnabled bool         `json:"PowerRegulationEnabled" yaml:"PowerRegulationEnabled"`
	Status                 string       `json:"Status" yaml:"Status"`
	Target                 string       `json:"Target" yaml:"Target"`
}
type EthernetNICInfo struct {
	RedfishId           string `json:"RedfishId" yaml:"RedfishId"`
	Oid                 string `json:"@odata.id" yaml:"@odata.id"`
	Description         string `json:"Description,omitempty" yaml:"Description,omitempty"`
	FQDN                string `json:"FQDN,omitempty" yaml:"FQDN,omitempty"`
	Hostname            string `json:"Hostname,omitempty" yaml:"Hostname,omitempty"`
	InterfaceEnabled    *bool  `json:"InterfaceEnabled,omitempty" yaml:"InterfaceEnabled,omitempty"`
	MACAddress          string `json:"MACAddress" yaml:"MACAddress"`
	PermanentMACAddress string `json:"PermanentMACAddress,omitempty" yaml:"PermanentMACAddress,omitempty"`
}
type Control struct {
	URL     string    `json:"URL" yaml:"URL"`
	Control RFControl `json:"Control" yaml:"Control"`
}
type ComponentManagerInfo struct {
	Name         string             `json:"Name,omitempty" yaml:"Name,omitempty"`
	Actions      *ManagerActions    `json:"Actions,omitempty" yaml:"Actions,omitempty"`
	EthNICInfo   []*EthernetNICInfo `json:"EthernetNICInfo,omitempty" yaml:"EthernetNICInfo,omitempty"`
	CommandShell *CommandShell      `json:"CommandShell,omitempty" yaml:"CommandShell,omitempty"`
}
type ComponentPDUInfo struct {
	Name    string                    `json:"Name,omitempty" yaml:"Name,omitempty"`
	Actions *PowerDistributionActions `json:"Actions,omitempty" yaml:"Actions,omitempty"`
}
type ChassisActions struct {
	ChassisReset ActionReset        `json:"#Chassis.Reset" yaml:"#Chassis.Reset"`
	OEM          *ChassisActionsOEM `json:"Oem,omitempty" yaml:"Oem,omitempty"`
}
type ChassisActionsOEM struct {
	ChassisEmergencyPower *ActionReset `json:"#Chassis.EmergencyPower,omitempty" yaml:"#Chassis.EmergencyPower,omitempty"`
}
type ActionReset struct {
	AllowableValues []string `json:"ResetType@Redfish.AllowableValues" yaml:"ResetType@Redfish.AllowableValues"`
	RFActionInfo    string   `json:"@Redfish.ActionInfo" yaml:"@Redfish.ActionInfo"`
	Target          string   `json:"target" yaml:"target"`
	Title           string   `json:"title,omitempty" yaml:"title,omitempty"`
}
type ManagerActionsOEM struct {
	ManagerFactoryReset *ActionFactoryReset `json:"#Manager.FactoryReset,omitempty" yaml:"#Manager.FactoryReset,omitempty"`
	CrayProcessSchedule *ActionNamed        `json:"#CrayProcess.Schedule,omitempty" yaml:"#CrayProcess.Schedule,omitempty"`
}
type ActionFactoryReset struct {
	AllowableValues []string `json:"FactoryResetType@Redfish.AllowableValues" yaml:"FactoryResetType@Redfish.AllowableValues"`
	Target          string   `json:"target" yaml:"target"`
	Title           string   `json:"title,omitempty" yaml:"title,omitempty"`
}
type ActionNamed struct {
	AllowableValues []string `json:"Name@Redfish.AllowableValues" yaml:"Name@Redfish.AllowableValues"`
	Target          string   `json:"target" yaml:"target"`
	Title           string   `json:"title,omitempty" yaml:"title,omitempty"`
}
type ComputerSystemActions struct {
	ComputerSystemReset ActionReset `json:"#ComputerSystem.Reset" yaml:"#ComputerSystem.Reset"`
}
type RFControl struct {
	ControlDelaySeconds int      `json:"ControlDelaySeconds" yaml:"ControlDelaySeconds"`
	ControlMode         string   `json:"ControlMode" yaml:"ControlMode"`
	ControlType         string   `json:"ControlType" yaml:"ControlType"`
	Id                  string   `json:"Id" yaml:"Id"`
	Name                string   `json:"Name" yaml:"Name"`
	PhysicalContext     string   `json:"PhysicalContext" yaml:"PhysicalContext"`
	SetPoint            int      `json:"SetPoint" yaml:"SetPoint"`
	SetPointUnits       string   `json:"SetPointUnits" yaml:"SetPointUnits"`
	SettingRangeMax     int      `json:"SettingRangeMax" yaml:"SettingRangeMax"`
	SettingRangeMin     int      `json:"SettingRangeMin" yaml:"SettingRangeMin"`
	Status              StatusRF `json:"Status" yaml:"Status"`
}
type StatusRF struct {
	Health       string `json:"Health" yaml:"Health"`
	HealthRollUp string `json:"HealthRollUp,omitempty" yaml:"HealthRollUp,omitempty"`
	State        string `json:"State,omitempty" yaml:"State,omitempty"`
}
type SerialConsole struct {
	MaxConcurrentSessions int                    `json:"MaxConcurrentSessions" yaml:"MaxConcurrentSessions"`
	IPMI                  *SerialConsoleProtocol `json:"IPMI,omitempty" yaml:"IPMI,omitempty"`
	SSH                   *SerialConsoleProtocol `json:"SSH,omitempty" yaml:"SSH,omitempty"`
	Telnet                *SerialConsoleProtocol `json:"Telnet,omitempty" yaml:"Telnet,omitempty"`
	WebSocket             *WebSocketConsole      `json:"WebSocket,omitempty" yaml:"WebSocket,omitempty"`
}
type SerialConsoleProtocol struct {
	ServiceEnabled        bool   `json:"ServiceEnabled" yaml:"ServiceEnabled"`
	Port                  int    `json:"Port,omitempty" yaml:"Port,omitempty"`
	HotKeySequenceDisplay string `json:"HotKeySequenceDisplay,omitempty" yaml:"HotKeySequenceDisplay,omitempty"`
	SharedWithManagerCLI  bool   `json:"SharedWithManagerCLI,omitempty" yaml:"SharedWithManagerCLI,omitempty"`
	ConsoleEntryCommand   string `json:"ConsoleEntryCommand,omitempty" yaml:"ConsoleEntryCommand,omitempty"`
}
type WebSocketConsole struct {
	ServiceEnabled bool   `json:"ServiceEnabled" yaml:"ServiceEnabled"`
	Interactive    bool   `json:"Interactive" yaml:"Interactive"`
	ConsoleURI     string `json:"ConsoleURI" yaml:"ConsoleURI"`
}
type ManagerActions struct {
	ManagerReset ActionReset        `json:"#Manager.Reset" yaml:"#Manager.Reset"`
	OEM          *ManagerActionsOEM `json:"Oem,omitempty" yaml:"Oem,omitempty"`
}
type CommandShell struct {
	ServiceEnabled        bool     `json:"ServiceEnabled" yaml:"ServiceEnabled"`
	MaxConcurrentSessions int      `json:"MaxConcurrentSessions" yaml:"MaxConcurrentSessions"`
	ConnectTypesSupported []string `json:"ConnectTypesSupported" yaml:"ConnectTypesSupported"`
}
type PowerDistributionActions struct {
	OEM *json.RawMessage `json:"Oem,omitempty" yaml:"Oem,omitempty"`
}
