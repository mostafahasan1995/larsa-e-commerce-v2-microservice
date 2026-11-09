package order_models

import (
	"time"

	"larsa-e-commerce-microservice-v2/internal/common/types"
	order_enums "larsa-e-commerce-microservice-v2/internal/modules/order/enums"
)

type ShipmentDetails struct {
	ShipmentMethod string                     `bson:"shipmentMethod" json:"shipmentMethod"`
	ShipmentDate   *time.Time                 `bson:"shipmentDate,omitempty" json:"shipmentDate,omitempty"`
	ShipmentStatus order_enums.ShipmentStatus `bson:"shipmentStatus" json:"shipmentStatus"`

	EstimationDeliveryDate *time.Time `bson:"estimationDeliveryDate,omitempty" json:"estimationDeliveryDate,omitempty"`
	DorToDorShipmentCost   float64    `bson:"dorToDorShipmentCost" json:"dorToDorShipmentCost"`

	ShipmentComplianceAndGoodsDetails *ShipmentComplianceAndGoodsDetails `bson:"otherShipmentDetails,omitempty" json:"otherShipmentDetails,omitempty"`

	ShipmentOrigin      ShipmentDestination `bson:"shipmentOrigin" json:"shipmentOrigin"`
	ShipmentDestination ShipmentDestination `bson:"shipmentDestination" json:"shipmentDestination"`
}

type ShipmentDestination struct {
	Country   *string  `bson:"country" json:"country,omitempty"`
	City      *string  `bson:"city" json:"city,omitempty"`
	Area      *string  `bson:"area,omitempty" json:"area,omitempty"`
	Address   *string  `bson:"address,omitempty" json:"address,omitempty"`
	Notes     *string  `bson:"notes,omitempty" json:"notes,omitempty"`
	Latitude  *float64 `bson:"latitude,omitempty" json:"latitude,omitempty"`
	Longitude *float64 `bson:"longitude,omitempty" json:"longitude,omitempty"`
}

type ShipmentComplianceAndGoodsDetails struct {
	Packing               []Packing              `bson:"packing" json:"packing"`
	CertificateOfOrigin   *CertificateOfOrigin   `bson:"certificateOfOrigin,omitempty" json:"certificateOfOrigin,omitempty"`
	ProductSpecifications *ProductSpecifications `bson:"productSpecifications,omitempty" json:"productSpecifications,omitempty"`
	QualityComplianceCert *types.FileFieldDto    `bson:"qualityComplianceCert,omitempty" json:"qualityComplianceCert,omitempty"`
	ExportLicense         *ExportLicense         `bson:"exportLicense,omitempty" json:"exportLicense,omitempty"`
}

type Packing struct {
	Name        string           `bson:"name" json:"name"`
	Content     string           `bson:"content" json:"content"`
	Dimensions  string           `bson:"dimensions" json:"dimensions"`
	Unit        order_enums.Unit `bson:"unit" json:"unit"`
	TotalWeight float64          `bson:"totalWeight" json:"totalWeight"`
}

type CertificateOfOrigin struct {
	File            types.FileFieldDto `bson:"file" json:"file"`
	CountryOfOrigin string             `bson:"countryOfOrigin" json:"countryOfOrigin"`
}

type ProductSpecifications struct {
	Details string             `bson:"details" json:"details"`
	File    types.FileFieldDto `bson:"file" json:"file"`
}

type ExportLicense struct {
	File            types.FileFieldDto `bson:"file" json:"file"`
	LicenseRequired bool               `bson:"licenseRequired" json:"licenseRequired"`
}
