package quotation_utils

import (
	"context"
	"time"

	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	"larsa-e-commerce-microservice-v2/internal/modules/base/models"
	quotation_dto "larsa-e-commerce-microservice-v2/internal/modules/quotation/dto"
	quotation_enums "larsa-e-commerce-microservice-v2/internal/modules/quotation/enums"
	quotation_models "larsa-e-commerce-microservice-v2/internal/modules/quotation/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ConvertCreateQuotationDTOToModel(dto *quotation_dto.CreateQuotationDTO, ctx context.Context) (quotation_models.Quotation, *helpers.ApiError) {
	cfg, err := util.GetReqAppCfg(ctx)
	if err != nil {
		return quotation_models.Quotation{}, &helpers.ApiError{
			StatusCode: 500,
			Message:    "GetReqAppCfg went wrong",
		}
	}

	now := time.Now()

	// Convert products
	products := make([]quotation_models.QuotationProduct, len(dto.Products))
	for i, product := range dto.Products {
		products[i] = ConvertQuotationProductDTOToModel(&product)
	}

	toLocation := quotation_models.QuotationDeliveryLocation{
		CountryIso3: dto.Delivery.ToLocation.CountryIso3,
		City:        dto.Delivery.ToLocation.City,
		Area:        dto.Delivery.ToLocation.Area,
		Address:     dto.Delivery.ToLocation.Address,
		Latitude:    dto.Delivery.ToLocation.Latitude,
		Longitude:   dto.Delivery.ToLocation.Longitude,
	}

	fromLocation := quotation_models.QuotationDeliveryLocation{
		CountryIso3: dto.Delivery.FromLocation.CountryIso3,
		City:        dto.Delivery.FromLocation.City,
		Area:        dto.Delivery.FromLocation.Area,
		Address:     dto.Delivery.FromLocation.Address,
		Latitude:    dto.Delivery.FromLocation.Latitude,
		Longitude:   dto.Delivery.FromLocation.Longitude,
	}
	// Convert delivery
	delivery := quotation_models.Delivery{
		ProposedDeliveryDate: dto.Delivery.ProposedDeliveryDate,
		DeliveryTerms:        dto.Delivery.DeliveryTerms,
		ToLocation:           toLocation,
		ShippingMethod:       dto.Delivery.ShippingMethod,
		FromLocation:         fromLocation,
	}

	// Convert payment
	payment := quotation_models.Payment{
		PreferredPaymentMethod: dto.Payment.PreferredPaymentMethod,
		PaymentSchedule:        dto.Payment.PaymentSchedule,
		AdvancePaymentPercent:  &dto.Payment.AdvancePaymentPercent,
	}

	// Convert warranty and validity
	warrantyAndValidity := quotation_models.WarrantyAndValidity{
		Details:                   dto.WarrantyAndValidity.Details,
		ReturnReplacementPolicies: dto.WarrantyAndValidity.ReturnReplacementPolicies,
		QuotationValidityPeriod:   dto.WarrantyAndValidity.QuotationValidityPeriod,
	}
	var additionalNotes *quotation_models.AdditionalNotes
	if dto.AdditionalNotes != nil {
		additionalNotes = &quotation_models.AdditionalNotes{
			RemarksComments: dto.AdditionalNotes.RemarksComments,
			Commitment:      dto.AdditionalNotes.Commitment,
			SpecialNotes:    dto.AdditionalNotes.SpecialNotes,
		}
	}

	// Convert signature
	var signature *quotation_models.Signature
	if dto.Signature != nil {
		signature = &quotation_models.Signature{
			OfferResponsibleName: dto.Signature.OfferResponsibleName,
			Position:             dto.Signature.Position,
			OfferSubmissionDate:  dto.Signature.OfferSubmissionDate,
			SignatureDoc:         *dto.Signature.SignatureDoc,
		}
	}

	packagingRequirements := quotation_models.PackagingRequirements{
		PackingDetails:         dto.PackagingRequirements.PackingDetails,
		LabelRequirement:       dto.PackagingRequirements.LabelRequirement,
		RequiredCertifications: dto.PackagingRequirements.RequiredCertifications,
		SampleRequirement:      dto.PackagingRequirements.SampleRequirement,
		SampleDeliveryDeadline: dto.PackagingRequirements.SampleDeliveryDeadline,
	}

	discount := quotation_models.QuotationPrice{}
	if dto.DiscountAmount != nil {
		discount = quotation_models.QuotationPrice{
			Amount:       dto.DiscountAmount.Amount,
			CurrencyCode: dto.DiscountAmount.CurrencyCode,
		}
	}

	quotation := quotation_models.Quotation{
		BaseQuotation: quotation_models.BaseQuotation{
			BaseModel: models.BaseModel{
				Id: primitive.NewObjectID(),
				MetaDataBaseModel: models.MetaDataBaseModel{
					CreatedAt: now,
					CreatedBy: &cfg.User.Id,
					UpdatedAt: now,
					UpdatedBy: &cfg.User.Id,
				},
			},
			Title:                 dto.Title,
			Type:                  dto.Type,
			Deadline:              dto.Deadline,
			Status:                quotation_enums.PENDING_ADMIN,
			Delivery:              delivery,
			Payment:               payment,
			DiscountAmount:        &discount,
			WarrantyAndValidity:   warrantyAndValidity,
			AdditionalNotes:       *additionalNotes,
			Attachments:           dto.Attachments,
			Signature:             signature,
			VendorProfileId:       dto.VendorProfileId,
			PackagingRequirements: packagingRequirements,
		},
		QuotationBuyerInfo: ConvertQuotationBuyerInfoDTOToModel(dto.QuotationBuyerInfo),
		Products:           products,
	}
	quotation.Code = quotation.GenerateCode()
	return quotation, nil
}

func ConvertQuotationBuyerInfoDTOToModel(dto *quotation_dto.QuotationBuyerInfoDTO) *quotation_models.QuotationBuyerInfo {
	var quotationBuyerInfo *quotation_models.QuotationBuyerInfo
	if dto == nil {
		return nil
	}
	quotationBuyerInfo = &quotation_models.QuotationBuyerInfo{
		BaseQuotationBuyerInfo: quotation_models.BaseQuotationBuyerInfo{
			BuyerProfileId:    dto.UserId,
			ContactPerson:     dto.ContactPerson,
			PhoneNumber:       dto.PhoneNumber,
			EmailAddress:      dto.EmailAddress,
			CompanyName:       dto.CompanyName,
			CompanyAddress:    dto.CompanyAddress,
			BusinessRegNo:     dto.BusinessRegNo,
			TaxRegistrationNo: dto.TaxRegistrationNo,
		},
	}
	return quotationBuyerInfo
}

// convert comment dto to comment model
func ConvertCommentDTOToModel(dto *quotation_dto.CommentDTO, user *primitive.ObjectID) quotation_models.Comment {
	return quotation_models.Comment{
		Comment:     dto.Comment,
		CreatedAt:   time.Now(),
		CreatedById: user,
	}
}

func ConvertQuotationProductDTOToModel(dto *quotation_dto.QuotationProductDTO) quotation_models.QuotationProduct {
	attributes := make([]quotation_models.QuotationProductAttribute, len(dto.Attributes))
	for i, attribute := range dto.Attributes {
		attributes[i] = ConvertQuotationProductAttributeDTOToModel(&attribute)
	}
	return quotation_models.QuotationProduct{
		BaseQuotationProduct: quotation_models.BaseQuotationProduct{
			ProductId:  dto.ProductId,
			Quantity:   dto.Quantity,
			Attributes: attributes,
			Price: quotation_models.QuotationPrice{
				Amount:       dto.Price.Amount,
				CurrencyCode: dto.Price.CurrencyCode,
			},
			DiscountPercent: dto.DiscountPercent,
		},
	}
}

func ConvertQuotationProductAttributeDTOToModel(dto *quotation_dto.QuotationProductAttributeDTO) quotation_models.QuotationProductAttribute {
	return quotation_models.QuotationProductAttribute{
		Code:  dto.Code,
		Value: dto.Value,
		Type:  dto.Type,
	}
}
