package quotation_utils

import (
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	quotation_enums "larsa-e-commerce-microservice-v2/internal/modules/quotation/enums"
	quotation_models "larsa-e-commerce-microservice-v2/internal/modules/quotation/models"
	"time"
)

// Check if update quotation status is valid based on the current status
func IsUpdateQuotationStatusValid(quotation *quotation_models.QuotationResponse, newStatus quotation_enums.QuotationStatus) bool {
	switch quotation.Status {
	case quotation_enums.PENDING_ADMIN:
		return newStatus == quotation_enums.OPEN || newStatus == quotation_enums.READY_TO_SEND
	case quotation_enums.READY_TO_SEND:
		return newStatus == quotation_enums.SENT_TO_BUYER || newStatus == quotation_enums.DRAFT
	case quotation_enums.SENT_TO_BUYER:
		return newStatus == quotation_enums.APPROVED || newStatus == quotation_enums.REJECTED
	case quotation_enums.APPROVED:
		return newStatus == quotation_enums.CONVERT_TO_ORDER
	case quotation_enums.CONVERT_TO_ORDER:
		return newStatus == quotation_enums.CONVERTED_TO_ORDER
	}
	return false
}

func GetHistoryEntry(quotation *quotation_models.QuotationResponse, newStatus quotation_enums.QuotationStatus) *quotation_models.QuotationHistory {
	// Determine the appropriate history action based on the status change
	var action quotation_enums.QuotationHistoryAction
	switch newStatus {
	case quotation_enums.READY_TO_SEND:
		action = quotation_enums.SUBMITTED_TO_BEEBUY
	case quotation_enums.SENT_TO_BUYER:
		action = quotation_enums.QUOTATION_SENT
	case quotation_enums.REJECTED:
		action = quotation_enums.QUOTATION_REJECTED
	case quotation_enums.APPROVED:
		action = quotation_enums.QUOTATION_ACCEPTED
	case quotation_enums.CONVERTED_TO_ORDER:
		action = quotation_enums.ORDER_CONVERSION
	default:
		action = quotation_enums.QUOTATION_UPDATED
	}

	additionalInfo := "Status changed from " + string(quotation.Status) + " to " + string(newStatus)
	historyEntry := &quotation_models.QuotationHistory{
		Date:           time.Now(),
		Action:         action,
		AdditionalInfo: &additionalInfo,
	}
	return historyEntry
}

// Check if user is allowed to update quotation status
func IsUserAllowedToUpdateQuotationStatus(quotation *quotation_models.QuotationResponse, newStatus quotation_enums.QuotationStatus, user *types.User) bool {
	userId := user.Id.Hex()
	vendorProfileId := quotation.VendorProfileId.Hex()
	buyerProfileId := quotation.QuotationBuyerInfo.BuyerProfileId.Hex()

	if newStatus == quotation_enums.READY_TO_SEND || newStatus == quotation_enums.OPEN {
		//Super admin (READY_TO_SEND)
		return user.Type == enums.SuperAdmin

	} else if newStatus == quotation_enums.PENDING_ADMIN {
		return vendorProfileId == userId
	} else if newStatus == quotation_enums.SENT_TO_BUYER {
		return vendorProfileId == userId
	} else if newStatus == quotation_enums.CONVERT_TO_ORDER {
		return vendorProfileId == userId
	} else if newStatus == quotation_enums.DRAFT {
		return vendorProfileId == userId

	} else if newStatus == quotation_enums.APPROVED {
		return buyerProfileId == userId
	} else if newStatus == quotation_enums.REJECTED {
		return buyerProfileId == userId
	} else if newStatus == quotation_enums.CONVERTED_TO_ORDER {
		return true
	}

	return false
}

// CreateHistoryEntry creates a new history entry for quotation updates
func CreateHistoryEntry(action quotation_enums.QuotationHistoryAction, additionalInfo string) quotation_models.QuotationHistory {
	return quotation_models.QuotationHistory{
		Date:           time.Now(),
		Action:         action,
		AdditionalInfo: &additionalInfo,
	}
}
