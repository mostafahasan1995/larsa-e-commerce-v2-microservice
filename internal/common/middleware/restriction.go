package middleware

import (
	"larsa-e-commerce-microservice-v2/internal/common/enums"
)

func CreateRestriction(module enums.ModuleCapability, operation enums.OperationType) Restriction {
	return Restriction{
		SuperAdminCapability: getCapability(module.SuperAdmin, operation),
		OwnerCapability:      getCapability(module.Vendor, operation),
		MemberCapability:     getCapability(module.Member, operation),
	}
}

func getCapability(userType enums.UserTypeOps, operation enums.OperationType) string {
	switch operation {
	case enums.Post:
		return userType.Post.String()
	case enums.Get:
		return userType.Get.String()
	case enums.GetOne:
		return userType.GetOne.String()
	case enums.Put:
		return userType.Put.String()
	case enums.Delete:
		return userType.Delete.String()
	case enums.DeleteMany:
		return userType.DeleteMany.String()
	default:
		return ""
	}
}
