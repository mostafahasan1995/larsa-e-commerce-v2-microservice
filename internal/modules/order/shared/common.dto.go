package order_shared

import "larsa-e-commerce-microservice-v2/internal/common/types"

type AdditionalInformationDto struct {
	Image *types.FileFieldDto `json:"image,omitempty"`
	Logo  *types.FileFieldDto `json:"logo,omitempty"`
	Note  *string             `json:"note,omitempty"`
}
