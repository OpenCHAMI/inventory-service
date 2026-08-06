package v1

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/openchami/fabrica/pkg/fabrica"
	"github.com/openchami/inventory-service/schemas"
)

type Hardware struct {
	APIVersion string           `json:"apiVersion" yaml:"apiVersion"`
	Kind       string           `json:"kind" yaml:"kind"`
	Metadata   fabrica.Metadata `json:"metadata" yaml:"metadata"`
	ID         string           `json:"id,omitempty" yaml:"id,omitempty"`
	Spec       HardwareSpec     `json:"spec" yaml:"spec" validate:"required"`
	Status     HardwareStatus   `json:"status,omitempty" yaml:"status,omitempty"`
}

type HardwareSpec struct {
	Description string `json:"description,omitempty" yaml:"description,omitempty" validate:"max=200"`
	ID          string `json:"ID" yaml:"ID"`
	Type        string `json:"Type" yaml:"Type"`
	Ordinal     int    `json:"Ordinal" yaml:"Ordinal"`
	Status      string `json:"Status" yaml:"Status"`

	HWInventoryByLocationType string `json:"HWInventoryByLocationType" yaml:"HWInventoryByLocationType"`

	HMSCabinetLocationInfo *ChassisLocationInfoRF `json:"CabinetLocationInfo,omitempty" yaml:"CabinetLocationInfo,omitempty"`
	HMSChassisLocationInfo *ChassisLocationInfoRF `json:"ChassisLocationInfo,omitempty" yaml:"ChassisLocationInfo,omitempty"`

	HMSComputeModuleLocationInfo *ChassisLocationInfoRF `json:"ComputeModuleLocationInfo,omitempty" yaml:"ComputeModuleLocationInfo,omitempty"`
	HMSRouterModuleLocationInfo  *ChassisLocationInfoRF `json:"RouterModuleLocationInfo,omitempty" yaml:"RouterModuleLocationInfo,omitempty"`
	HMSNodeEnclosureLocationInfo *ChassisLocationInfoRF `json:"NodeEnclosureLocationInfo,omitempty" yaml:"NodeEnclosureLocationInfo,omitempty"`

	HMSHSNBoardLocationInfo     *ChassisLocationInfoRF `json:"HSNBoardLocationInfo,omitempty" yaml:"HSNBoardLocationInfo,omitempty"`
	HMSMgmtSwitchLocationInfo   *ChassisLocationInfoRF `json:"MgmtSwitchLocationInfo,omitempty" yaml:"MgmtSwitchLocationInfo,omitempty"`
	HMSMgmtHLSwitchLocationInfo *ChassisLocationInfoRF `json:"MgmtHLSwitchLocationInfo,omitempty" yaml:"MgmtHLSwitchLocationInfo,omitempty"`

	HMSCDUMgmtSwitchLocationInfo *ChassisLocationInfoRF `json:"CDUMgmtSwitchLocationInfo,omitempty" yaml:"CDUMgmtSwitchLocationInfo,omitempty"`

	HMSNodeLocationInfo *SystemLocationInfoRF `json:"NodeLocationInfo,omitempty" yaml:"NodeLocationInfo,omitempty"`

	HMSProcessorLocationInfo                *ProcessorLocationInfoRF       `json:"ProcessorLocationInfo,omitempty" yaml:"ProcessorLocationInfo,omitempty"`
	HMSNodeAccelLocationInfo                *ProcessorLocationInfoRF       `json:"NodeAccelLocationInfo,omitempty" yaml:"NodeAccelLocationInfo,omitempty"`
	HMSMemoryLocationInfo                   *MemoryLocationInfoRF          `json:"MemoryLocationInfo,omitempty" yaml:"MemoryLocationInfo,omitempty"`
	HMSDriveLocationInfo                    *DriveLocationInfoRF           `json:"DriveLocationInfo,omitempty" yaml:"DriveLocationInfo,omitempty"`
	HMSHSNNICLocationInfo                   *NALocationInfoRF              `json:"NodeHsnNicLocationInfo,omitempty" yaml:"NodeHsnNicLocationInfo,omitempty"`
	HMSPDULocationInfo                      *PowerDistributionLocationInfo `json:"PDULocationInfo,omitempty" yaml:"PDULocationInfo,omitempty"`
	HMSOutletLocationInfo                   *OutletLocationInfo            `json:"OutletLocationInfo,omitempty" yaml:"OutletLocationInfo,omitempty"`
	HMSCMMRectifierLocationInfo             *PowerSupplyLocationInfoRF     `json:"CMMRectifierLocationInfo,omitempty" yaml:"CMMRectifierLocationInfo,omitempty"`
	HMSNodeEnclosurePowerSupplyLocationInfo *PowerSupplyLocationInfoRF     `json:"NodeEnclosurePowerSupplyLocationInfo,omitempty" yaml:"NodeEnclosurePowerSupplyLocationInfo,omitempty"`
	HMSNodeBMCLocationInfo                  *ManagerLocationInfoRF         `json:"NodeBMCLocationInfo,omitempty" yaml:"NodeBMCLocationInfo,omitempty"`
	HMSRouterBMCLocationInfo                *ManagerLocationInfoRF         `json:"RouterBMCLocationInfo,omitempty" yaml:"RouterBMCLocationInfo,omitempty"`
	HMSNodeAccelRiserLocationInfo           *NodeAccelRiserLocationInfoRF  `json:"NodeAccelRiserLocationInfo,omitempty" yaml:"NodeAccelRiserLocationInfo,omitempty"`
	PopulatedFRU                            *HWInvByFRU                    `json:"PopulatedFRU,omitempty" yaml:"PopulatedFRU,omitempty"`
	hmsTypeArrays
}

type HardwareStatus struct {
	Phase   string `json:"phase,omitempty" yaml:"phase,omitempty"`
	Message string `json:"message,omitempty" yaml:"message,omitempty"`
	Ready   bool   `json:"ready" yaml:"ready"`
}

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

type HWInvByFRU struct {
	FRUID                              string                    `json:"FRUID" yaml:"FRUID"`
	Type                               string                    `json:"Type" yaml:"Type"`
	Subtype                            string                    `json:"Subtype" yaml:"Subtype"`
	HWInventoryByFRUType               string                    `json:"HWInventoryByFRUType" yaml:"HWInventoryByFRUType"`
	HMSCabinetFRUInfo                  *ChassisFRUInfoRF         `json:"CabinetFRUInfo,omitempty" yaml:"CabinetFRUInfo,omitempty"`
	HMSChassisFRUInfo                  *ChassisFRUInfoRF         `json:"ChassisFRUInfo,omitempty" yaml:"ChassisFRUInfo,omitempty"`
	HMSComputeModuleFRUInfo            *ChassisFRUInfoRF         `json:"ComputeModuleFRUInfo,omitempty" yaml:"ComputeModuleFRUInfo,omitempty"`
	HMSRouterModuleFRUInfo             *ChassisFRUInfoRF         `json:"RouterModuleFRUInfo,omitempty" yaml:"RouterModuleFRUInfo,omitempty"`
	HMSNodeEnclosureFRUInfo            *ChassisFRUInfoRF         `json:"NodeEnclosureFRUInfo,omitempty" yaml:"NodeEnclosureFRUInfo,omitempty"`
	HMSHSNBoardFRUInfo                 *ChassisFRUInfoRF         `json:"HSNBoardFRUInfo,omitempty" yaml:"HSNBoardFRUInfo,omitempty"`
	HMSMgmtSwitchFRUInfo               *ChassisFRUInfoRF         `json:"MgmtSwitchFRUInfo,omitempty" yaml:"MgmtSwitchFRUInfo,omitempty"`
	HMSMgmtHLSwitchFRUInfo             *ChassisFRUInfoRF         `json:"MgmtHLSwitchFRUInfo,omitempty" yaml:"MgmtHLSwitchFRUInfo,omitempty"`
	HMSCDUMgmtSwitchFRUInfo            *ChassisFRUInfoRF         `json:"CDUMgmtSwitchFRUInfo,omitempty" yaml:"CDUMgmtSwitchFRUInfo,omitempty"`
	HMSNodeFRUInfo                     *SystemFRUInfoRF          `json:"NodeFRUInfo,omitempty" yaml:"NodeFRUInfo,omitempty"`
	HMSProcessorFRUInfo                *ProcessorFRUInfoRF       `json:"ProcessorFRUInfo,omitempty" yaml:"ProcessorFRUInfo,omitempty"`
	HMSNodeAccelFRUInfo                *ProcessorFRUInfoRF       `json:"NodeAccelFRUInfo,omitempty" yaml:"NodeAccelFRUInfo,omitempty"`
	HMSMemoryFRUInfo                   *MemoryFRUInfoRF          `json:"MemoryFRUInfo,omitempty" yaml:"MemoryFRUInfo,omitempty"`
	HMSDriveFRUInfo                    *DriveFRUInfoRF           `json:"DriveFRUInfo,omitempty" yaml:"DriveFRUInfo,omitempty"`
	HMSHSNNICFRUInfo                   *NAFRUInfoRF              `json:"NodeHsnNicFRUInfo,omitempty" yaml:"NodeHsnNicFRUInfo,omitempty"`
	HMSPDUFRUInfo                      *PowerDistributionFRUInfo `json:"PDUFRUInfo,omitempty" yaml:"PDUFRUInfo,omitempty"`
	HMSOutletFRUInfo                   *OutletFRUInfo            `json:"OutletFRUInfo,omitempty" yaml:"OutletFRUInfo,omitempty"`
	HMSCMMRectifierFRUInfo             *PowerSupplyFRUInfoRF     `json:"CMMRectifierFRUInfo,omitempty" yaml:"CMMRectifierFRUInfo,omitempty"`
	HMSNodeEnclosurePowerSupplyFRUInfo *PowerSupplyFRUInfoRF     `json:"NodeEnclosurePowerSupplyFRUInfo,omitempty" yaml:"NodeEnclosurePowerSupplyFRUInfo,omitempty"`
	HMSNodeBMCFRUInfo                  *ManagerFRUInfoRF         `json:"NodeBMCFRUInfo,omitempty" yaml:"NodeBMCFRUInfo,omitempty"`
	HMSRouterBMCFRUInfo                *ManagerFRUInfoRF         `json:"RouterBMCFRUInfo,omitempty" yaml:"RouterBMCFRUInfo,omitempty"`
	HMSNodeAccelRiserFRUInfo           *NodeAccelRiserFRUInfoRF  `json:"NodeAccelRiserFRUInfo,omitempty" yaml:"NodeAccelRiserFRUInfo,omitempty"`
}
type ChassisLocationInfoRF struct {
	Id          string `json:"Id" yaml:"Id"`
	Name        string `json:"Name" yaml:"Name"`
	Description string `json:"Description" yaml:"Description"`
	Hostname    string `json:"HostName" yaml:"HostName"`
}
type ComputerSystemProcessorSummary struct {
	Count json.Number `json:"Count" yaml:"Count"`
	Model string      `json:"Model" yaml:"Model"`
}
type ComputerSystemMemorySummary struct {
	TotalSystemMemoryGiB json.Number `json:"TotalSystemMemoryGiB" yaml:"TotalSystemMemoryGiB"`
}
type SystemLocationInfoRF struct {
	Id               string                         `json:"Id" yaml:"Id"`
	Name             string                         `json:"Name" yaml:"Name"`
	Description      string                         `json:"Description" yaml:"Description"`
	Hostname         string                         `json:"HostName" yaml:"HostName"`
	ProcessorSummary ComputerSystemProcessorSummary `json:"ProcessorSummary" yaml:"ProcessorSummary"`
	MemorySummary    ComputerSystemMemorySummary    `json:"MemorySummary" yaml:"MemorySummary"`
}
type ProcessorLocationInfoRF struct {
	Id          string `json:"Id" yaml:"Id"`
	Name        string `json:"Name" yaml:"Name"`
	Description string `json:"Description" yaml:"Description"`
	Socket      string `json:"Socket" yaml:"Socket"`
}
type MemoryLocationInfoRF struct {
	Id             string           `json:"Id" yaml:"Id"`
	Name           string           `json:"Name" yaml:"Name"`
	Description    string           `json:"Description" yaml:"Description"`
	MemoryLocation MemoryLocationRF `json:"MemoryLocation" yaml:"MemoryLocation"`
}
type MemoryLocationRF struct {
	Socket           json.Number `json:"Socket" yaml:"Socket"`
	MemoryController json.Number `json:"MemoryController" yaml:"MemoryController"`
	Channel          json.Number `json:"Channel" yaml:"Channel"`
	Slot             json.Number `json:"Slot" yaml:"Slot"`
}
type ManagerLocationInfoRF struct {
	DateTime            string `json:"DateTime" yaml:"DateTime"`
	DateTimeLocalOffset string `json:"DateTimeLocalOffset" yaml:"DateTimeLocalOffset"`
	Description         string `json:"Description" yaml:"Description"`
	FirmwareVersion     string `json:"FirmwareVersion" yaml:"FirmwareVersion"`
	Id                  string `json:"Id" yaml:"Id"`
	Name                string `json:"Name" yaml:"Name"`
}
type NodeAccelRiserLocationInfoRF struct {
	Name        string `json:"Name" yaml:"Name"`
	Description string `json:"Description" yaml:"Description"`
}
type PowerSupplyLocationInfoRF struct {
	Name            string `json:"Name" yaml:"Name"`
	FirmwareVersion string `json:"FirmwareVersion" yaml:"FirmwareVersion"`
}
type DriveLocationInfoRF struct {
	Id          string `json:"Id" yaml:"Id"`
	Name        string `json:"Name" yaml:"Name"`
	Description string `json:"Description" yaml:"Description"`
}
type NALocationInfoRF struct {
	Id          string `json:"Id" yaml:"Id"`
	Name        string `json:"Name" yaml:"Name"`
	Description string `json:"Description" yaml:"Description"`
}
type PowerDistributionLocationInfo struct {
	Id          string    `json:"Id" yaml:"Id"`
	Description string    `json:"Description" yaml:"Description"`
	Name        string    `json:"Name" yaml:"Name"`
	UUID        string    `json:"UUID" yaml:"UUID"`
	Location    *Location `json:"Location,omitempty" yaml:"Location,omitempty"`
}
type Location struct {
	ContactInfo   *ContactInfo   `json:"ContactInfo,omitempty" yaml:"ContactInfo,omitempty"`
	Latitude      json.Number    `json:"Latitude,omitempty" yaml:"Latitude,omitempty"`
	Longitude     json.Number    `json:"Longitude,omitempty" yaml:"Longitude,omitempty"`
	PartLocation  *PartLocation  `json:"PartLocation,omitempty" yaml:"PartLocation,omitempty"`
	Placement     *Placement     `json:"Placement,omitempty" yaml:"Placement,omitempty"`
	PostalAddress *PostalAddress `json:"PostalAddress,omitempty" yaml:"PostalAddress,omitempty"`
}
type ContactInfo struct {
	ContactName  string `json:"ContactName" yaml:"ContactName"`
	EmailAddress string `json:"EmailAddress" yaml:"EmailAddress"`
	PhoneNumber  string `json:"PhoneNumber,omitempty" yaml:"PhoneNumber,omitempty"`
}
type PartLocation struct {
	LocationOrdinalValue json.Number `json:"LocationOrdinalValue,omitempty" yaml:"LocationOrdinalValue,omitempty"`
	LocationType         string      `json:"LocationType" yaml:"LocationType"`
	Orientation          string      `json:"Orientation" yaml:"Orientation"`
	Reference            string      `json:"Reference" yaml:"Reference"`
	ServiceLabel         string      `json:"ServiceLabel" yaml:"ServiceLabel"`
}
type PostalAddress struct {
	Country    string `json:"Country" yaml:"Country"`
	Territory  string `json:"Territory" yaml:"Territory"`
	City       string `json:"City" yaml:"City"`
	Street     string `json:"Street" yaml:"Street"`
	Name       string `json:"Name" yaml:"Name"`
	PostalCode string `json:"PostalCode" yaml:"PostalCode"`
	Building   string `json:"Building" yaml:"Building"`
	Floor      string `json:"Floor" yaml:"Floor"`
	Room       string `json:"Room" yaml:"Room"`
}
type Placement struct {
	AdditionalInfo  string      `json:"AdditionalInfo,omitempty" yaml:"AdditionalInfo,omitempty"`
	Rack            string      `json:"Rack,omitempty" yaml:"Rack,omitempty"`
	RackOffset      json.Number `json:"RackOffset,omitempty" yaml:"RackOffset,omitempty"`
	RackOffsetUnits string      `json:"RackOffsetUnits,omitempty" yaml:"RackOffsetUnits,omitempty"`
	Row             string      `json:"Row,omitempty" yaml:"Row,omitempty"`
}
type OutletLocationInfo struct {
	Id          string `json:"Id" yaml:"Id"`
	Description string `json:"Description" yaml:"Description"`
	Name        string `json:"Name" yaml:"Name"`
}
type ChassisFRUInfoRF struct {
	AssetTag     string `json:"AssetTag" yaml:"AssetTag"`
	ChassisType  string `json:"ChassisType" yaml:"ChassisType"`
	Model        string `json:"Model" yaml:"Model"`
	Manufacturer string `json:"Manufacturer" yaml:"Manufacturer"`
	PartNumber   string `json:"PartNumber" yaml:"PartNumber"`
	SerialNumber string `json:"SerialNumber" yaml:"SerialNumber"`
	SKU          string `json:"SKU" yaml:"SKU"`
}
type SystemFRUInfoRF struct {
	AssetTag     string `json:"AssetTag" yaml:"AssetTag"`
	BiosVersion  string `json:"BiosVersion" yaml:"BiosVersion"`
	Model        string `json:"Model" yaml:"Model"`
	Manufacturer string `json:"Manufacturer" yaml:"Manufacturer"`
	PartNumber   string `json:"PartNumber" yaml:"PartNumber"`
	SerialNumber string `json:"SerialNumber" yaml:"SerialNumber"`
	SKU          string `json:"SKU" yaml:"SKU"`
	SystemType   string `json:"SystemType" yaml:"SystemType"`
	UUID         string `json:"UUID" yaml:"UUID"`
}
type ProcessorFRUInfoRF struct {
	InstructionSet        string        `json:"InstructionSet" yaml:"InstructionSet"`
	Manufacturer          string        `json:"Manufacturer" yaml:"Manufacturer"`
	MaxSpeedMHz           json.Number   `json:"MaxSpeedMHz" yaml:"MaxSpeedMHz"`
	Model                 string        `json:"Model" yaml:"Model"`
	SerialNumber          string        `json:"SerialNumber" yaml:"SerialNumber"`
	PartNumber            string        `json:"PartNumber" yaml:"PartNumber"`
	ProcessorArchitecture string        `json:"ProcessorArchitecture" yaml:"ProcessorArchitecture"`
	ProcessorId           ProcessorIdRF `json:"ProcessorId" yaml:"ProcessorId"`
	ProcessorType         string        `json:"ProcessorType" yaml:"ProcessorType"`
	TotalCores            json.Number   `json:"TotalCores" yaml:"TotalCores"`
	TotalThreads          json.Number   `json:"TotalThreads" yaml:"TotalThreads"`
	Oem                   *ProcessorOEM `json:"Oem" yaml:"Oem"`
}
type ProcessorOEM struct {
	GBTProcessorOemProperty *GBTProcessorOem `json:"GBTProcessorOemProperty,omitempty" yaml:"GBTProcessorOemProperty,omitempty"`
}
type GBTProcessorOem struct {
	ProcessorSerialNumber string `json:"Processor Serial Number,omitempty" yaml:"Processor Serial Number,omitempty"`
}
type ProcessorIdRF struct {
	EffectiveFamily         string `json:"EffectiveFamily" yaml:"EffectiveFamily"`
	EffectiveModel          string `json:"EffectiveModel" yaml:"EffectiveModel"`
	IdentificationRegisters string `json:"IdentificationRegisters" yaml:"IdentificationRegisters"`
	MicrocodeInfo           string `json:"MicrocodeInfo" yaml:"MicrocodeInfo"`
	Step                    string `json:"Step" yaml:"Step"`
	VendorID                string `json:"VendorID" yaml:"VendorID"`
}
type MemoryFRUInfoRF struct {
	BaseModuleType    string      `json:"BaseModuleType,omitempty" yaml:"BaseModuleType,omitempty"`
	BusWidthBits      json.Number `json:"BusWidthBits,omitempty" yaml:"BusWidthBits,omitempty"`
	CapacityMiB       json.Number `json:"CapacityMiB" yaml:"CapacityMiB"`
	DataWidthBits     json.Number `json:"DataWidthBits,omitempty" yaml:"DataWidthBits,omitempty"`
	ErrorCorrection   string      `json:"ErrorCorrection,omitempty" yaml:"ErrorCorrection,omitempty"`
	Manufacturer      string      `json:"Manufacturer,omitempty" yaml:"Manufacturer,omitempty"`
	MemoryType        string      `json:"MemoryType,omitempty" yaml:"MemoryType,omitempty"`
	MemoryDeviceType  string      `json:"MemoryDeviceType,omitempty" yaml:"MemoryDeviceType,omitempty"`
	OperatingSpeedMhz json.Number `json:"OperatingSpeedMhz" yaml:"OperatingSpeedMhz"`
	PartNumber        string      `json:"PartNumber,omitempty" yaml:"PartNumber,omitempty"`
	RankCount         json.Number `json:"RankCount,omitempty" yaml:"RankCount,omitempty"`
	SerialNumber      string      `json:"SerialNumber" yaml:"SerialNumber"`
}
type hmsTypeArrays struct {
	Nodes                      *[]*HardwareSpec `json:"Nodes,omitempty" yaml:"Nodes,omitempty"`
	Cabinets                   *[]*HardwareSpec `json:"Cabinets,omitempty" yaml:"Cabinets,omitempty"`
	Chassis                    *[]*HardwareSpec `json:"Chassis,omitempty" yaml:"Chassis,omitempty"`
	ComputeModules             *[]*HardwareSpec `json:"ComputeModules,omitempty" yaml:"ComputeModules,omitempty"`
	RouterModules              *[]*HardwareSpec `json:"RouterModules,omitempty" yaml:"RouterModules,omitempty"`
	NodeEnclosures             *[]*HardwareSpec `json:"NodeEnclosures,omitempty" yaml:"NodeEnclosures,omitempty"`
	HSNBoards                  *[]*HardwareSpec `json:"HSNBoards,omitempty" yaml:"HSNBoards,omitempty"`
	Processors                 *[]*HardwareSpec `json:"Processors,omitempty" yaml:"Processors,omitempty"`
	Memory                     *[]*HardwareSpec `json:"Memory,omitempty" yaml:"Memory,omitempty"`
	Drives                     *[]*HardwareSpec `json:"Drives,omitempty" yaml:"Drives,omitempty"`
	CabinetPDUs                *[]*HardwareSpec `json:"CabinetPDUs,omitempty" yaml:"CabinetPDUs,omitempty"`
	CabinetPDUOutlets          *[]*HardwareSpec `json:"CabinetPDUPowerConnectors,omitempty" yaml:"CabinetPDUPowerConnectors,omitempty"`
	CMMRectifiers              *[]*HardwareSpec `json:"CMMRectifiers,omitempty" yaml:"CMMRectifiers,omitempty"`
	NodeAccels                 *[]*HardwareSpec `json:"NodeAccels,omitempty" yaml:"NodeAccels,omitempty"`
	NodeAccelRisers            *[]*HardwareSpec `json:"NodeAccelRisers,omitempty" yaml:"NodeAccelRisers,omitempty"`
	NodeEnclosurePowerSupplies *[]*HardwareSpec `json:"NodeEnclosurePowerSupplies,omitempty" yaml:"NodeEnclosurePowerSupplies,omitempty"`
	NodeHsnNICs                *[]*HardwareSpec `json:"NodeHsnNics,omitempty" yaml:"NodeHsnNics,omitempty"`
	CECs                       *[]*HardwareSpec `json:"CECs,omitempty" yaml:"CECs,omitempty"`
	CDUs                       *[]*HardwareSpec `json:"CDUs,omitempty" yaml:"CDUs,omitempty"`
	CabinetCDUs                *[]*HardwareSpec `json:"CabinetCDUs,omitempty" yaml:"CabinetCDUs,omitempty"`
	CMMFpgas                   *[]*HardwareSpec `json:"CMMFpgas,omitempty" yaml:"CMMFpgas,omitempty"`
	NodeFpgas                  *[]*HardwareSpec `json:"NodeFpgas,omitempty" yaml:"NodeFpgas,omitempty"`
	RouterFpgas                *[]*HardwareSpec `json:"RouterFpgas,omitempty" yaml:"RouterFpgas,omitempty"`
	RouterTORFpgas             *[]*HardwareSpec `json:"RouterTORFpgas,omitempty" yaml:"RouterTORFpgas,omitempty"`
	HSNAsics                   *[]*HardwareSpec `json:"HSNAsics,omitempty" yaml:"HSNAsics,omitempty"`
	CabinetBMCs                *[]*HardwareSpec `json:"CabinetBMCs,omitempty" yaml:"CabinetBMCs,omitempty"`
	CabinetPDUControllers      *[]*HardwareSpec `json:"CabinetPDUControllers,omitempty" yaml:"CabinetPDUControllers,omitempty"`
	ChassisBMCs                *[]*HardwareSpec `json:"ChassisBMCs,omitempty" yaml:"ChassisBMCs,omitempty"`
	NodeBMCs                   *[]*HardwareSpec `json:"NodeBMCs,omitempty" yaml:"NodeBMCs,omitempty"`
	RouterBMCs                 *[]*HardwareSpec `json:"RouterBMCs,omitempty" yaml:"RouterBMCs,omitempty"`
	CabinetPDUNics             *[]*HardwareSpec `json:"CabinetPDUNics,omitempty" yaml:"CabinetPDUNics,omitempty"`
	NodePowerConnectors        *[]*HardwareSpec `json:"NodePowerConnectors,omitempty" yaml:"NodePowerConnectors,omitempty"`
	NodeBMCNics                *[]*HardwareSpec `json:"NodeBMCNics,omitempty" yaml:"NodeBMCNics,omitempty"`
	NodeNICs                   *[]*HardwareSpec `json:"NodeNICs,omitempty" yaml:"NodeNICs,omitempty"`
	RouterBMCNics              *[]*HardwareSpec `json:"RouterBMCNics,omitempty" yaml:"RouterBMCNics,omitempty"`
	MgmtSwitches               *[]*HardwareSpec `json:"MgmtSwitches,omitempty" yaml:"MgmtSwitches,omitempty"`
	MgmtHLSwitches             *[]*HardwareSpec `json:"MgmtHLSwitches,omitempty" yaml:"MgmtHLSwitches,omitempty"`
	CDUMgmtSwitches            *[]*HardwareSpec `json:"CDUMgmtSwitches,omitempty" yaml:"CDUMgmtSwitches,omitempty"`
	SMSBoxes                   *[]*HardwareSpec `json:"SMSBoxes,omitempty" yaml:"SMSBoxes,omitempty"`
	HSNLinks                   *[]*HardwareSpec `json:"HSNLinks,omitempty" yaml:"HSNLinks,omitempty"`
	HSNConnectors              *[]*HardwareSpec `json:"HSNConnectors,omitempty" yaml:"HSNConnectors,omitempty"`
	HSNConnectorPorts          *[]*HardwareSpec `json:"HSNConnectorPorts,omitempty" yaml:"HSNConnectorPorts,omitempty"`
	MgmtSwitchConnectors       *[]*HardwareSpec `json:"MgmtSwitchConnectors,omitempty" yaml:"MgmtSwitchConnectors,omitempty"`
}
type DriveFRUInfoRF struct {
	Manufacturer                  string      `json:"Manufacturer" yaml:"Manufacturer"`
	SerialNumber                  string      `json:"SerialNumber" yaml:"SerialNumber"`
	PartNumber                    string      `json:"PartNumber" yaml:"PartNumber"`
	Model                         string      `json:"Model" yaml:"Model"`
	SKU                           string      `json:"SKU" yaml:"SKU"`
	CapacityBytes                 json.Number `json:"CapacityBytes" yaml:"CapacityBytes"`
	Protocol                      string      `json:"Protocol" yaml:"Protocol"`
	MediaType                     string      `json:"MediaType" yaml:"MediaType"`
	RotationSpeedRPM              json.Number `json:"RotationSpeedRPM" yaml:"RotationSpeedRPM"`
	BlockSizeBytes                json.Number `json:"BlockSizeBytes" yaml:"BlockSizeBytes"`
	CapableSpeedGbs               json.Number `json:"CapableSpeedGbs" yaml:"CapableSpeedGbs"`
	FailurePredicted              bool        `json:"FailurePredicted" yaml:"FailurePredicted"`
	EncryptionAbility             string      `json:"EncryptionAbility" yaml:"EncryptionAbility"`
	EncryptionStatus              string      `json:"EncryptionStatus" yaml:"EncryptionStatus"`
	NegotiatedSpeedGbs            json.Number `json:"NegotiatedSpeedGbs" yaml:"NegotiatedSpeedGbs"`
	PredictedMediaLifeLeftPercent json.Number `json:"PredictedMediaLifeLeftPercent" yaml:"PredictedMediaLifeLeftPercent"`
}
type NAFRUInfoRF struct {
	Manufacturer string `json:"Manufacturer" yaml:"Manufacturer"`
	Model        string `json:"Model" yaml:"Model"`
	PartNumber   string `json:"PartNumber" yaml:"PartNumber"`
	SKU          string `json:"SKU,omitempty" yaml:"SKU,omitempty"`
	SerialNumber string `json:"SerialNumber" yaml:"SerialNumber"`
}
type PowerDistributionFRUInfo struct {
	AssetTag          string         `json:"AssetTag" yaml:"AssetTag"`
	DateOfManufacture string         `json:"DateOfManufacture,omitempty" yaml:"DateOfManufacture,omitempty"`
	EquipmentType     string         `json:"EquipmentType" yaml:"EquipmentType"`
	FirmwareVersion   string         `json:"FirmwareVersion" yaml:"FirmwareVersion"`
	HardwareRevision  string         `json:"HardwareRevision" yaml:"HardwareRevision"`
	Manufacturer      string         `json:"Manufacturer" yaml:"Manufacturer"`
	Model             string         `json:"Model" yaml:"Model"`
	PartNumber        string         `json:"PartNumber" yaml:"PartNumber"`
	SerialNumber      string         `json:"SerialNumber" yaml:"SerialNumber"`
	CircuitSummary    CircuitSummary `json:"CircuitSummary" yaml:"CircuitSummary"`
}
type CircuitSummary struct {
	ControlledOutlets json.Number `json:"ControlledOutlets,omitempty" yaml:"ControlledOutlets,omitempty"`
	MonitoredBranches json.Number `json:"MonitoredBranches,omitempty" yaml:"MonitoredBranches,omitempty"`
	MonitoredOutlets  json.Number `json:"MonitoredOutlets,omitempty" yaml:"MonitoredOutlets,omitempty"`
	MonitoredPhases   json.Number `json:"MonitoredPhases,omitempty" yaml:"MonitoredPhases,omitempty"`
	TotalBranches     json.Number `json:"TotalBranches,omitempty" yaml:"TotalBranches,omitempty"`
	TotalOutlets      json.Number `json:"TotalOutlets,omitempty" yaml:"TotalOutlets,omitempty"`
	TotalPhases       json.Number `json:"TotalPhases,omitempty" yaml:"TotalPhases,omitempty"`
}
type OutletFRUInfo struct {
	NominalVoltage   string         `json:"NominalVoltage,omitempty" yaml:"NominalVoltage,omitempty"`
	OutletType       string         `json:"OutletType" yaml:"OutletType"`
	EnergySensor     *SensorExcerpt `json:"EnergySensor,omitempty" yaml:"EnergySensor,omitempty"`
	FrequencySensor  *SensorExcerpt `json:"FrequencySensor,omitempty" yaml:"FrequencySensor,omitempty"`
	PhaseWiringType  string         `json:"PhaseWiringType,omitempty" yaml:"PhaseWiringType,omitempty"`
	PowerEnabled     *bool          `json:"PowerEnabled,omitempty" yaml:"PowerEnabled,omitempty"`
	RatedCurrentAmps json.Number    `json:"RatedCurrentAmps,omitempty" yaml:"RatedCurrentAmps,omitempty"`
	VoltageType      string         `json:"VoltageType,omitempty" yaml:"VoltageType,omitempty"`
}
type SensorExcerpt struct {
	DataSourceUri      string      `json:"DataSourceUri" yaml:"DataSourceUri"`
	Name               string      `json:"Name" yaml:"Name"`
	PeakReading        json.Number `json:"PeakReading,omitempty" yaml:"PeakReading,omitempty"`
	PhysicalContext    string      `json:"PhysicalContext,omitempty" yaml:"PhysicalContext,omitempty"`
	PhysicalSubContext string      `json:"PhysicalSubContext,omitempty" yaml:"PhysicalSubContext,omitempty"`
	Reading            json.Number `json:"Reading,omitempty" yaml:"Reading,omitempty"`
	ReadingUnits       string      `json:"ReadingUnits,omitempty" yaml:"ReadingUnits,omitempty"`
	Status             StatusRF    `json:"Status,omitempty" yaml:"Status,omitempty"`
}
type PowerSupplyFRUInfoRF struct {
	Manufacturer       string      `json:"Manufacturer" yaml:"Manufacturer"`
	SerialNumber       string      `json:"SerialNumber" yaml:"SerialNumber"`
	Model              string      `json:"Model" yaml:"Model"`
	PartNumber         string      `json:"PartNumber" yaml:"PartNumber"`
	PowerCapacityWatts int         `json:"PowerCapacityWatts" yaml:"PowerCapacityWatts"`
	PowerInputWatts    int         `json:"PowerInputWatts" yaml:"PowerInputWatts"`
	PowerOutputWatts   interface{} `json:"PowerOutputWatts" yaml:"PowerOutputWatts"`
	PowerSupplyType    string      `json:"PowerSupplyType" yaml:"PowerSupplyType"`
}
type ManagerFRUInfoRF struct {
	ManagerType  string `json:"ManagerType" yaml:"ManagerType"`
	Model        string `json:"Model" yaml:"Model"`
	Manufacturer string `json:"Manufacturer" yaml:"Manufacturer"`
	PartNumber   string `json:"PartNumber" yaml:"PartNumber"`
	SerialNumber string `json:"SerialNumber" yaml:"SerialNumber"`
}
type NodeAccelRiserFRUInfoRF struct {
	PhysicalContext        string             `json:"PhysicalContext" yaml:"PhysicalContext"`
	Producer               string             `json:"Producer" yaml:"Producer"`
	SerialNumber           string             `json:"SerialNumber" yaml:"SerialNumber"`
	PartNumber             string             `json:"PartNumber" yaml:"PartNumber"`
	Model                  string             `json:"Model" yaml:"Model"`
	ProductionDate         string             `json:"ProductionDate" yaml:"ProductionDate"`
	Version                string             `json:"Version" yaml:"Version"`
	EngineeringChangeLevel string             `json:"EngineeringChangeLevel" yaml:"EngineeringChangeLevel"`
	OEM                    *NodeAccelRiserOEM `json:"Oem,omitempty" yaml:"Oem,omitempty"`
}
type NodeAccelRiserOEM struct {
	PCBSerialNumber string `json:"PCBSerialNumber" yaml:"PCBSerialNumber"`
}
