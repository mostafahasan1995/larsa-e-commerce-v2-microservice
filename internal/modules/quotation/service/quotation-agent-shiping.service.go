package quotation_service

import (
	"context"
	"fmt"

	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	"larsa-e-commerce-microservice-v2/internal/common/util"
	quotation_dto "larsa-e-commerce-microservice-v2/internal/modules/quotation/dto"
	quotation_enums "larsa-e-commerce-microservice-v2/internal/modules/quotation/enums"
	quotation_models "larsa-e-commerce-microservice-v2/internal/modules/quotation/models"
	quotation_utils "larsa-e-commerce-microservice-v2/internal/modules/quotation/utils"
)

func (s *quotationAdmin) SetShippingCost(ctx context.Context, quotationID string, dto quotation_dto.SetShippingCostDTO) (*quotation_models.QuotationResponse, *helpers.ApiError) {
	quotation, apiErr := s.GetById(ctx, quotationID, "")
	if apiErr != nil {
		return nil, apiErr
	}
	cfg, _ := util.GetReqAppCfg(ctx)

	// Authorize: SuperAdmin OR (shipping company assigned equals caller) OR (from-agent equals caller)
	if cfg != nil && cfg.User != nil && cfg.User.Type != enums.SuperAdmin {
		isShippingCompany := quotation.ShippingCompanyId != nil && quotation.ShippingCompanyId.Hex() == cfg.User.Id.Hex()
		isFromAgent := quotation.AgentFromId != nil && quotation.AgentFromId.Hex() == cfg.User.Id.Hex()
		if !isShippingCompany && !isFromAgent {
			return nil, &helpers.ApiError{StatusCode: 403, Message: "Not authorized to set shipping cost"}
		}
	}

	if quotation.Status != quotation_enums.PENDING_ADMIN {
		return nil, &helpers.ApiError{StatusCode: 400, Message: "Can't set shipping cost for this quotation, due to the quotation is not in pending admin status"}
	}

	update := map[string]interface{}{
		"fee.deliveryFee.freightCharges": quotation_models.QuotationPrice{Amount: dto.FreightCharges.Amount, CurrencyCode: dto.FreightCharges.CurrencyCode},
		"fee.deliveryFee.insurance":      quotation_models.QuotationPrice{Amount: dto.Insurance.Amount, CurrencyCode: dto.Insurance.CurrencyCode},
	}

	updated, err := s.repo.UpdateFieldsById(ctx, quotation.Id, update, cfg.User.Id)
	if err != nil {
		return nil, err
	}

	// Update history
	additionalInfo := "Shipping cost updated: freight charges and insurance set"
	historyEntry := quotation_utils.CreateHistoryEntry(quotation_enums.QUOTATION_UPDATED, additionalInfo)
	newHistory := append(updated.History, historyEntry)

	_, historyErr := s.repo.UpdateFieldsById(ctx, quotation.Id, map[string]interface{}{"history": newHistory}, cfg.User.Id)
	if historyErr != nil {
		fmt.Println("Error updating history:", historyErr)
	}

	return updated, nil
}

func (s *quotationAdmin) SetAgentFee(ctx context.Context, quotationID string, dto quotation_dto.SetAgentFeeDTO) (*quotation_models.QuotationResponse, *helpers.ApiError) {
	cfg, _ := util.GetReqAppCfg(ctx)

	quotation, apiErr := s.GetById(ctx, quotationID, "")
	if apiErr != nil {
		return nil, apiErr
	}
	update := map[string]interface{}{}

	if cfg.User.Id.Hex() == quotation.AgentFromId.Hex() {
		// From agent
		update = s.setFromAgentFee(ctx, dto)
	} else if cfg.User.Id.Hex() == quotation.AgentToId.Hex() {
		// To agent
		update = s.setToAgentFee(ctx, dto)
	}

	if len(update) == 0 {
		return quotation, nil
	}

	if quotation.Status != quotation_enums.PENDING_ADMIN {
		return nil, &helpers.ApiError{StatusCode: 400, Message: "Can't set agent fee for this quotation, due to the quotation is not in pending admin status"}
	}

	updated, err := s.repo.UpdateFieldsById(ctx, quotation.Id, update, cfg.User.Id)
	if err != nil {
		return nil, err
	}

	// Update history
	var additionalInfo string
	if cfg.User.Id.Hex() == quotation.AgentFromId.Hex() {
		additionalInfo = "Agent fee updated: export declaration fees set by from-agent"
	} else if cfg.User.Id.Hex() == quotation.AgentToId.Hex() {
		additionalInfo = "Agent fee updated: customs clearance, import declaration, and delivery fees set by to-agent"
	} else {
		additionalInfo = "Agent fee updated"
	}

	historyEntry := quotation_utils.CreateHistoryEntry(quotation_enums.QUOTATION_UPDATED, additionalInfo)
	newHistory := append(updated.History, historyEntry)

	_, historyErr := s.repo.UpdateFieldsById(ctx, quotation.Id, map[string]interface{}{"history": newHistory}, cfg.User.Id)
	if historyErr != nil {
		fmt.Println("Error updating history:", historyErr)
	}

	return updated, nil
}

func (s *quotationAdmin) setFromAgentFee(ctx context.Context, dto quotation_dto.SetAgentFeeDTO) map[string]interface{} {
	update := map[string]interface{}{
		"fee.deliveryFee.exportDeclarationFees": quotation_models.QuotationPrice{Amount: dto.ExportDeclarationFees.Amount, CurrencyCode: dto.ExportDeclarationFees.CurrencyCode},
	}
	return update
}

func (s *quotationAdmin) setToAgentFee(ctx context.Context, dto quotation_dto.SetAgentFeeDTO) map[string]interface{} {
	update := map[string]interface{}{
		"fee.deliveryFee.customsClearanceFee":   quotation_models.QuotationPrice{Amount: dto.CustomsClearanceFee.Amount, CurrencyCode: dto.CustomsClearanceFee.CurrencyCode},
		"fee.deliveryFee.importDeclarationFees": quotation_models.QuotationPrice{Amount: dto.ImportDeclarationFees.Amount, CurrencyCode: dto.ImportDeclarationFees.CurrencyCode},
		"fee.deliveryFee.deliveryToBuyerFee":    quotation_models.QuotationPrice{Amount: dto.DeliveryToBuyerFee.Amount, CurrencyCode: dto.DeliveryToBuyerFee.CurrencyCode},
	}
	return update
}
