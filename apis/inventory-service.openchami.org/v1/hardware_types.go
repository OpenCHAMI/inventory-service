// Copyright © 2026 OpenCHAMI a Series of LF Projects, LLC
//
// SPDX-License-Identifier: MIT

package v1

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/OpenCHAMI/inventory-service/schemas"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/openchami/fabrica/pkg/fabrica"
)

type Hardware struct {
	APIVersion string           `json:"apiVersion"`
	Kind       string           `json:"kind"`
	Metadata   fabrica.Metadata `json:"metadata"`
	ID         string           `json:"id,omitempty"`
	Spec       HardwareSpec     `json:"spec" validate:"required"`
	Status     HardwareStatus   `json:"status,omitempty"`
}

type HardwareSpec struct {
	Description string `json:"description,omitempty" validate:"max=200"`
	ID          string `json:"ID"`
	Type        string `json:"Type"`
	Ordinal     int    `json:"Ordinal"`
	Status      string `json:"Status"`

	HWInventoryByLocationType string `json:"HWInventoryByLocationType"`

	HMSCabinetLocationInfo       *ChassisLocationInfoRF `json:"CabinetLocationInfo,omitempty"`
	HMSChassisLocationInfo       *ChassisLocationInfoRF `json:"ChassisLocationInfo,omitempty"`
	HMSComputeModuleLocationInfo *ChassisLocationInfoRF `json:"ComputeModuleLocationInfo,omitempty"`
	HMSRouterModuleLocationInfo  *ChassisLocationInfoRF `json:"RouterModuleLocationInfo,omitempty"`
	HMSNodeEnclosureLocationInfo *ChassisLocationInfoRF `json:"NodeEnclosureLocationInfo,omitempty"`
	HMSHSNBoardLocationInfo      *ChassisLocationInfoRF `json:"HSNBoardLocationInfo,omitempty"`
	HMSMgmtSwitchLocationInfo    *ChassisLocationInfoRF `json:"MgmtSwitchLocationInfo,omitempty"`
	HMSMgmtHLSwitchLocationInfo  *ChassisLocationInfoRF `json:"MgmtHLSwitchLocationInfo,omitempty"`
	HMSCDUMgmtSwitchLocationInfo *ChassisLocationInfoRF `json:"CDUMgmtSwitchLocationInfo,omitempty"`

	HMSNodeLocationInfo                     *SystemLocationInfoRF          `json:"NodeLocationInfo,omitempty"`
	HMSProcessorLocationInfo                *ProcessorLocationInfoRF       `json:"ProcessorLocationInfo,omitempty"`
	HMSNodeAccelLocationInfo                *ProcessorLocationInfoRF       `json:"NodeAccelLocationInfo,omitempty"`
	HMSMemoryLocationInfo                   *MemoryLocationInfoRF          `json:"MemoryLocationInfo,omitempty"`
	HMSDriveLocationInfo                    *DriveLocationInfoRF           `json:"DriveLocationInfo,omitempty"`
	HMSHSNNICLocationInfo                   *NALocationInfoRF              `json:"NodeHsnNicLocationInfo,omitempty"`
	HMSPDULocationInfo                      *PowerDistributionLocationInfo `json:"PDULocationInfo,omitempty"`
	HMSOutletLocationInfo                   *OutletLocationInfo            `json:"OutletLocationInfo,omitempty"`
	HMSCMMRectifierLocationInfo             *PowerSupplyLocationInfoRF     `json:"CMMRectifierLocationInfo,omitempty"`
	HMSNodeEnclosurePowerSupplyLocationInfo *PowerSupplyLocationInfoRF     `json:"NodeEnclosurePowerSupplyLocationInfo,omitempty"`
	HMSNodeBMCLocationInfo                  *ManagerLocationInfoRF         `json:"NodeBMCLocationInfo,omitempty"`
	HMSRouterBMCLocationInfo                *ManagerLocationInfoRF         `json:"RouterBMCLocationInfo,omitempty"`
	HMSNodeAccelRiserLocationInfo           *NodeAccelRiserLocationInfoRF  `json:"NodeAccelRiserLocationInfo,omitempty"`
	// TODO: Remaining types in hmsTypeArrays
	PopulatedFRU *HWInvByFRU `json:"PopulatedFRU,omitempty"`
	hmsTypeArrays
}

type HardwareStatus struct {
	Phase   string `json:"phase,omitempty"`
	Message string `json:"message,omitempty"`
	Ready   bool   `json:"ready"`
}

type HWInvByFRU struct {
	FRUID                              string                    `json:"FRUID"`
	Type                               string                    `json:"Type"`
	Subtype                            string                    `json:"Subtype"`
	HWInventoryByFRUType               string                    `json:"HWInventoryByFRUType"`
	HMSCabinetFRUInfo                  *ChassisFRUInfoRF         `json:"CabinetFRUInfo,omitempty"`
	HMSChassisFRUInfo                  *ChassisFRUInfoRF         `json:"ChassisFRUInfo,omitempty"`
	HMSComputeModuleFRUInfo            *ChassisFRUInfoRF         `json:"ComputeModuleFRUInfo,omitempty"`
	HMSRouterModuleFRUInfo             *ChassisFRUInfoRF         `json:"RouterModuleFRUInfo,omitempty"`
	HMSNodeEnclosureFRUInfo            *ChassisFRUInfoRF         `json:"NodeEnclosureFRUInfo,omitempty"`
	HMSHSNBoardFRUInfo                 *ChassisFRUInfoRF         `json:"HSNBoardFRUInfo,omitempty"`
	HMSMgmtSwitchFRUInfo               *ChassisFRUInfoRF         `json:"MgmtSwitchFRUInfo,omitempty"`
	HMSMgmtHLSwitchFRUInfo             *ChassisFRUInfoRF         `json:"MgmtHLSwitchFRUInfo,omitempty"`
	HMSCDUMgmtSwitchFRUInfo            *ChassisFRUInfoRF         `json:"CDUMgmtSwitchFRUInfo,omitempty"`
	HMSNodeFRUInfo                     *SystemFRUInfoRF          `json:"NodeFRUInfo,omitempty"`
	HMSProcessorFRUInfo                *ProcessorFRUInfoRF       `json:"ProcessorFRUInfo,omitempty"`
	HMSNodeAccelFRUInfo                *ProcessorFRUInfoRF       `json:"NodeAccelFRUInfo,omitempty"`
	HMSMemoryFRUInfo                   *MemoryFRUInfoRF          `json:"MemoryFRUInfo,omitempty"`
	HMSDriveFRUInfo                    *DriveFRUInfoRF           `json:"DriveFRUInfo,omitempty"`
	HMSHSNNICFRUInfo                   *NAFRUInfoRF              `json:"NodeHsnNicFRUInfo,omitempty"`
	HMSPDUFRUInfo                      *PowerDistributionFRUInfo `json:"PDUFRUInfo,omitempty"`
	HMSOutletFRUInfo                   *OutletFRUInfo            `json:"OutletFRUInfo,omitempty"`
	HMSCMMRectifierFRUInfo             *PowerSupplyFRUInfoRF     `json:"CMMRectifierFRUInfo,omitempty"`
	HMSNodeEnclosurePowerSupplyFRUInfo *PowerSupplyFRUInfoRF     `json:"NodeEnclosurePowerSupplyFRUInfo,omitempty"`
	HMSNodeBMCFRUInfo                  *ManagerFRUInfoRF         `json:"NodeBMCFRUInfo,omitempty"`
	HMSRouterBMCFRUInfo                *ManagerFRUInfoRF         `json:"RouterBMCFRUInfo,omitempty"`
	HMSNodeAccelRiserFRUInfo           *NodeAccelRiserFRUInfoRF  `json:"NodeAccelRiserFRUInfo,omitempty"`
}
type ChassisLocationInfoRF struct {
	Id          string `json:"Id"`
	Name        string `json:"Name"`
	Description string `json:"Description"`
	Hostname    string `json:"HostName"`
}
type ComputerSystemProcessorSummary struct {
	Count json.Number `json:"Count"`
	Model string      `json:"Model"`
}
type ComputerSystemMemorySummary struct {
	TotalSystemMemoryGiB json.Number `json:"TotalSystemMemoryGiB"`
}
type SystemLocationInfoRF struct {
	Id               string                         `json:"Id"`
	Name             string                         `json:"Name"`
	Description      string                         `json:"Description"`
	Hostname         string                         `json:"HostName"`
	ProcessorSummary ComputerSystemProcessorSummary `json:"ProcessorSummary"`
	MemorySummary    ComputerSystemMemorySummary    `json:"MemorySummary"`
}
type ProcessorLocationInfoRF struct {
	Id          string `json:"Id"`
	Name        string `json:"Name"`
	Description string `json:"Description"`
	Socket      string `json:"Socket"`
}
type MemoryLocationInfoRF struct {
	Id             string           `json:"Id"`
	Name           string           `json:"Name"`
	Description    string           `json:"Description"`
	MemoryLocation MemoryLocationRF `json:"MemoryLocation"`
}
type MemoryLocationRF struct {
	Socket           json.Number `json:"Socket"`
	MemoryController json.Number `json:"MemoryController"`
	Channel          json.Number `json:"Channel"`
	Slot             json.Number `json:"Slot"`
}
type ManagerLocationInfoRF struct {
	DateTime            string `json:"DateTime"`
	DateTimeLocalOffset string `json:"DateTimeLocalOffset"`
	Description         string `json:"Description"`
	FirmwareVersion     string `json:"FirmwareVersion"`
	Id                  string `json:"Id"`
	Name                string `json:"Name"`
}
type NodeAccelRiserLocationInfoRF struct {
	Name        string `json:"Name"`
	Description string `json:"Description"`
}
type PowerSupplyLocationInfoRF struct {
	Name            string `json:"Name"`
	FirmwareVersion string `json:"FirmwareVersion"`
}
type DriveLocationInfoRF struct {
	Id          string `json:"Id"`
	Name        string `json:"Name"`
	Description string `json:"Description"`
}
type NALocationInfoRF struct {
	Id          string `json:"Id"`
	Name        string `json:"Name"`
	Description string `json:"Description"`
}
type PowerDistributionLocationInfo struct {
	Id          string    `json:"Id"`
	Description string    `json:"Description"`
	Name        string    `json:"Name"`
	UUID        string    `json:"UUID"`
	Location    *Location `json:"Location,omitempty"`
}
type Location struct {
	ContactInfo   *ContactInfo   `json:"ContactInfo,omitempty"`
	Latitude      json.Number    `json:"Latitude,omitempty"`
	Longitude     json.Number    `json:"Longitude,omitempty"`
	PartLocation  *PartLocation  `json:"PartLocation,omitempty"`
	Placement     *Placement     `json:"Placement,omitempty"`
	PostalAddress *PostalAddress `json:"PostalAddress,omitempty"`
}
type ContactInfo struct {
	ContactName  string `json:"ContactName"`
	EmailAddress string `json:"EmailAddress"`
	PhoneNumber  string `json:"PhoneNumber,omitempty"`
}
type PartLocation struct {
	LocationOrdinalValue json.Number `json:"LocationOrdinalValue,omitempty"`
	LocationType         string      `json:"LocationType"`
	Orientation          string      `json:"Orientation"`
	Reference            string      `json:"Reference"`
	ServiceLabel         string      `json:"ServiceLabel"`
}
type PostalAddress struct {
	Country    string `json:"Country"`
	Territory  string `json:"Territory"`
	City       string `json:"City"`
	Street     string `json:"Street"`
	Name       string `json:"Name"`
	PostalCode string `json:"PostalCode"`
	Building   string `json:"Building"`
	Floor      string `json:"Floor"`
	Room       string `json:"Room"`
}
type Placement struct {
	AdditionalInfo  string      `json:"AdditionalInfo,omitempty"`
	Rack            string      `json:"Rack,omitempty"`
	RackOffset      json.Number `json:"RackOffset,omitempty"`
	RackOffsetUnits string      `json:"RackOffsetUnits,omitempty"`
	Row             string      `json:"Row,omitempty"`
}
type OutletLocationInfo struct {
	Id          string `json:"Id"`
	Description string `json:"Description"`
	Name        string `json:"Name"`
}
type ChassisFRUInfoRF struct {
	AssetTag     string `json:"AssetTag"`
	ChassisType  string `json:"ChassisType"`
	Model        string `json:"Model"`
	Manufacturer string `json:"Manufacturer"`
	PartNumber   string `json:"PartNumber"`
	SerialNumber string `json:"SerialNumber"`
	SKU          string `json:"SKU"`
}
type SystemFRUInfoRF struct {
	AssetTag     string `json:"AssetTag"`
	BiosVersion  string `json:"BiosVersion"`
	Model        string `json:"Model"`
	Manufacturer string `json:"Manufacturer"`
	PartNumber   string `json:"PartNumber"`
	SerialNumber string `json:"SerialNumber"`
	SKU          string `json:"SKU"`
	SystemType   string `json:"SystemType"`
	UUID         string `json:"UUID"`
}
type ProcessorFRUInfoRF struct {
	InstructionSet        string        `json:"InstructionSet"`
	Manufacturer          string        `json:"Manufacturer"`
	MaxSpeedMHz           json.Number   `json:"MaxSpeedMHz"`
	Model                 string        `json:"Model"`
	SerialNumber          string        `json:"SerialNumber"`
	PartNumber            string        `json:"PartNumber"`
	ProcessorArchitecture string        `json:"ProcessorArchitecture"`
	ProcessorId           ProcessorIdRF `json:"ProcessorId"`
	ProcessorType         string        `json:"ProcessorType"`
	TotalCores            json.Number   `json:"TotalCores"`
	TotalThreads          json.Number   `json:"TotalThreads"`
	Oem                   *ProcessorOEM `json:"Oem"`
}
type ProcessorOEM struct {
	GBTProcessorOemProperty *GBTProcessorOem `json:"GBTProcessorOemProperty,omitempty"`
}
type GBTProcessorOem struct {
	ProcessorSerialNumber string `json:"Processor Serial Number,omitempty"`
}
type ProcessorIdRF struct {
	EffectiveFamily         string `json:"EffectiveFamily"`
	EffectiveModel          string `json:"EffectiveModel"`
	IdentificationRegisters string `json:"IdentificationRegisters"`
	MicrocodeInfo           string `json:"MicrocodeInfo"`
	Step                    string `json:"Step"`
	VendorID                string `json:"VendorID"`
}
type MemoryFRUInfoRF struct {
	BaseModuleType    string      `json:"BaseModuleType,omitempty"`
	BusWidthBits      json.Number `json:"BusWidthBits,omitempty"`
	CapacityMiB       json.Number `json:"CapacityMiB"`
	DataWidthBits     json.Number `json:"DataWidthBits,omitempty"`
	ErrorCorrection   string      `json:"ErrorCorrection,omitempty"`
	Manufacturer      string      `json:"Manufacturer,omitempty"`
	MemoryType        string      `json:"MemoryType,omitempty"`
	MemoryDeviceType  string      `json:"MemoryDeviceType,omitempty"`
	OperatingSpeedMhz json.Number `json:"OperatingSpeedMhz"`
	PartNumber        string      `json:"PartNumber,omitempty"`
	RankCount         json.Number `json:"RankCount,omitempty"`
	SerialNumber      string      `json:"SerialNumber"`
}
type hmsTypeArrays struct {
	Nodes                      *[]*HardwareSpec `json:"Nodes,omitempty"`
	Cabinets                   *[]*HardwareSpec `json:"Cabinets,omitempty"`
	Chassis                    *[]*HardwareSpec `json:"Chassis,omitempty"`
	ComputeModules             *[]*HardwareSpec `json:"ComputeModules,omitempty"`
	RouterModules              *[]*HardwareSpec `json:"RouterModules,omitempty"`
	NodeEnclosures             *[]*HardwareSpec `json:"NodeEnclosures,omitempty"`
	HSNBoards                  *[]*HardwareSpec `json:"HSNBoards,omitempty"`
	Processors                 *[]*HardwareSpec `json:"Processors,omitempty"`
	Memory                     *[]*HardwareSpec `json:"Memory,omitempty"`
	Drives                     *[]*HardwareSpec `json:"Drives,omitempty"`
	CabinetPDUs                *[]*HardwareSpec `json:"CabinetPDUs,omitempty"`
	CabinetPDUOutlets          *[]*HardwareSpec `json:"CabinetPDUPowerConnectors,omitempty"`
	CMMRectifiers              *[]*HardwareSpec `json:"CMMRectifiers,omitempty"`
	NodeAccels                 *[]*HardwareSpec `json:"NodeAccels,omitempty"`
	NodeAccelRisers            *[]*HardwareSpec `json:"NodeAccelRisers,omitempty"`
	NodeEnclosurePowerSupplies *[]*HardwareSpec `json:"NodeEnclosurePowerSupplies,omitempty"`
	NodeHsnNICs                *[]*HardwareSpec `json:"NodeHsnNics,omitempty"`
	CECs                       *[]*HardwareSpec `json:"CECs,omitempty"`
	CDUs                       *[]*HardwareSpec `json:"CDUs,omitempty"`
	CabinetCDUs                *[]*HardwareSpec `json:"CabinetCDUs,omitempty"`
	CMMFpgas                   *[]*HardwareSpec `json:"CMMFpgas,omitempty"`
	NodeFpgas                  *[]*HardwareSpec `json:"NodeFpgas,omitempty"`
	RouterFpgas                *[]*HardwareSpec `json:"RouterFpgas,omitempty"`
	RouterTORFpgas             *[]*HardwareSpec `json:"RouterTORFpgas,omitempty"`
	HSNAsics                   *[]*HardwareSpec `json:"HSNAsics,omitempty"`
	CabinetBMCs                *[]*HardwareSpec `json:"CabinetBMCs,omitempty"`
	CabinetPDUControllers      *[]*HardwareSpec `json:"CabinetPDUControllers,omitempty"`
	ChassisBMCs                *[]*HardwareSpec `json:"ChassisBMCs,omitempty"`
	NodeBMCs                   *[]*HardwareSpec `json:"NodeBMCs,omitempty"`
	RouterBMCs                 *[]*HardwareSpec `json:"RouterBMCs,omitempty"`
	CabinetPDUNics             *[]*HardwareSpec `json:"CabinetPDUNics,omitempty"`
	NodePowerConnectors        *[]*HardwareSpec `json:"NodePowerConnectors,omitempty"`
	NodeBMCNics                *[]*HardwareSpec `json:"NodeBMCNics,omitempty"`
	NodeNICs                   *[]*HardwareSpec `json:"NodeNICs,omitempty"`
	RouterBMCNics              *[]*HardwareSpec `json:"RouterBMCNics,omitempty"`
	MgmtSwitches               *[]*HardwareSpec `json:"MgmtSwitches,omitempty"`
	MgmtHLSwitches             *[]*HardwareSpec `json:"MgmtHLSwitches,omitempty"`
	CDUMgmtSwitches            *[]*HardwareSpec `json:"CDUMgmtSwitches,omitempty"`
	SMSBoxes                   *[]*HardwareSpec `json:"SMSBoxes,omitempty"`
	HSNLinks                   *[]*HardwareSpec `json:"HSNLinks,omitempty"`
	HSNConnectors              *[]*HardwareSpec `json:"HSNConnectors,omitempty"`
	HSNConnectorPorts          *[]*HardwareSpec `json:"HSNConnectorPorts,omitempty"`
	MgmtSwitchConnectors       *[]*HardwareSpec `json:"MgmtSwitchConnectors,omitempty"`
}
type DriveFRUInfoRF struct {
	Manufacturer                  string      `json:"Manufacturer"`
	SerialNumber                  string      `json:"SerialNumber"`
	PartNumber                    string      `json:"PartNumber"`
	Model                         string      `json:"Model"`
	SKU                           string      `json:"SKU"`
	CapacityBytes                 json.Number `json:"CapacityBytes"`
	Protocol                      string      `json:"Protocol"`
	MediaType                     string      `json:"MediaType"`
	RotationSpeedRPM              json.Number `json:"RotationSpeedRPM"`
	BlockSizeBytes                json.Number `json:"BlockSizeBytes"`
	CapableSpeedGbs               json.Number `json:"CapableSpeedGbs"`
	FailurePredicted              bool        `json:"FailurePredicted"`
	EncryptionAbility             string      `json:"EncryptionAbility"`
	EncryptionStatus              string      `json:"EncryptionStatus"`
	NegotiatedSpeedGbs            json.Number `json:"NegotiatedSpeedGbs"`
	PredictedMediaLifeLeftPercent json.Number `json:"PredictedMediaLifeLeftPercent"`
}
type NAFRUInfoRF struct {
	Manufacturer string `json:"Manufacturer"`
	Model        string `json:"Model"`
	PartNumber   string `json:"PartNumber"`
	SKU          string `json:"SKU,omitempty"`
	SerialNumber string `json:"SerialNumber"`
}
type PowerDistributionFRUInfo struct {
	AssetTag          string         `json:"AssetTag"`
	DateOfManufacture string         `json:"DateOfManufacture,omitempty"`
	EquipmentType     string         `json:"EquipmentType"`
	FirmwareVersion   string         `json:"FirmwareVersion"`
	HardwareRevision  string         `json:"HardwareRevision"`
	Manufacturer      string         `json:"Manufacturer"`
	Model             string         `json:"Model"`
	PartNumber        string         `json:"PartNumber"`
	SerialNumber      string         `json:"SerialNumber"`
	CircuitSummary    CircuitSummary `json:"CircuitSummary"`
}
type CircuitSummary struct {
	ControlledOutlets json.Number `json:"ControlledOutlets,omitempty"`
	MonitoredBranches json.Number `json:"MonitoredBranches,omitempty"`
	MonitoredOutlets  json.Number `json:"MonitoredOutlets,omitempty"`
	MonitoredPhases   json.Number `json:"MonitoredPhases,omitempty"`
	TotalBranches     json.Number `json:"TotalBranches,omitempty"`
	TotalOutlets      json.Number `json:"TotalOutlets,omitempty"`
	TotalPhases       json.Number `json:"TotalPhases,omitempty"`
}
type OutletFRUInfo struct {
	NominalVoltage   string         `json:"NominalVoltage,omitempty"`
	OutletType       string         `json:"OutletType"`
	EnergySensor     *SensorExcerpt `json:"EnergySensor,omitempty"`
	FrequencySensor  *SensorExcerpt `json:"FrequencySensor,omitempty"`
	PhaseWiringType  string         `json:"PhaseWiringType,omitempty"`
	PowerEnabled     *bool          `json:"PowerEnabled,omitempty"`
	RatedCurrentAmps json.Number    `json:"RatedCurrentAmps,omitempty"`
	VoltageType      string         `json:"VoltageType,omitempty"`
}
type SensorExcerpt struct {
	DataSourceUri      string      `json:"DataSourceUri"`
	Name               string      `json:"Name"`
	PeakReading        json.Number `json:"PeakReading,omitempty"`
	PhysicalContext    string      `json:"PhysicalContext,omitempty"`
	PhysicalSubContext string      `json:"PhysicalSubContext,omitempty"`
	Reading            json.Number `json:"Reading,omitempty"`
	ReadingUnits       string      `json:"ReadingUnits,omitempty"`
	Status             StatusRF    `json:"Status,omitempty"`
}
type PowerSupplyFRUInfoRF struct {
	Manufacturer       string      `json:"Manufacturer"`
	SerialNumber       string      `json:"SerialNumber"`
	Model              string      `json:"Model"`
	PartNumber         string      `json:"PartNumber"`
	PowerCapacityWatts int         `json:"PowerCapacityWatts"`
	PowerInputWatts    int         `json:"PowerInputWatts"`
	PowerOutputWatts   interface{} `json:"PowerOutputWatts"`
	PowerSupplyType    string      `json:"PowerSupplyType"`
}
type ManagerFRUInfoRF struct {
	ManagerType  string `json:"ManagerType"`
	Model        string `json:"Model"`
	Manufacturer string `json:"Manufacturer"`
	PartNumber   string `json:"PartNumber"`
	SerialNumber string `json:"SerialNumber"`
}
type NodeAccelRiserFRUInfoRF struct {
	PhysicalContext        string             `json:"PhysicalContext"`
	Producer               string             `json:"Producer"`
	SerialNumber           string             `json:"SerialNumber"`
	PartNumber             string             `json:"PartNumber"`
	Model                  string             `json:"Model"`
	ProductionDate         string             `json:"ProductionDate"`
	Version                string             `json:"Version"`
	EngineeringChangeLevel string             `json:"EngineeringChangeLevel"`
	OEM                    *NodeAccelRiserOEM `json:"Oem,omitempty"`
}
type NodeAccelRiserOEM struct {
	PCBSerialNumber string `json:"PCBSerialNumber"`
}

// Validate implements custom validation logic for Hardware
func (r *Hardware) Validate(ctx context.Context) error {
	var schema jsonschema.Schema
	if err := json.Unmarshal(schemas.HardwareSchema, &schema); err != nil {
		return fmt.Errorf("loading hardware schema: %w", err)
	}

	resolved, err := schema.Resolve(nil)
	if err != nil {
		return fmt.Errorf("resolving hardware schema: %w", err)
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

func (r *Hardware) GetKind() string {
	return "Hardware"
}

func (r *Hardware) GetName() string {
	return r.Metadata.Name
}

func (r *Hardware) GetUID() string {
	return r.Metadata.UID
}

func (r *Hardware) IsHub() {}
