package contact_filters

import (
	"context"
	"encoding/json"
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	base_filters "larsa-e-commerce-microservice-v2/internal/modules/base/filters"

	"go.mongodb.org/mongo-driver/bson"
)

type AgentFilters struct {
	base_filters.BaseFilters
	FullName                     *string   `json:"fullName,omitempty"`
	CompanyName                  *string   `json:"companyName,omitempty"`
	LicenseNumber                *string   `json:"licenseNumber,omitempty"`
	YearsOfExperience            *int      `json:"yearsOfExperience,omitempty"`
	CommercialRegistrationNumber *string   `json:"commercialRegistrationNumber,omitempty"`
	TaxIdentificationNumber      *string   `json:"taxIdentificationNumber,omitempty"`
	PhoneNumber                  *string   `json:"phoneNumber,omitempty"`
	HeadOfficeAddress            *string   `json:"headOfficeAddress,omitempty"`
	OfficialEmail                *string   `json:"officialEmail,omitempty"`
	AuthorizedCountryCodes       *[]string `json:"authorizedCountryCodes,omitempty"`
	WebsiteLink                  *string   `json:"websiteLink,omitempty"`

	// Bank Account Information Filters
	BankName          *string `json:"bankName,omitempty"`
	AccountNumber     *string `json:"accountNumber,omitempty"`
	AccountHolderName *string `json:"accountHolderName,omitempty"`
	Iban              *string `json:"iban,omitempty"`
	SwiftCode         *string `json:"swiftCode,omitempty"`
}

func NewAgentFilters(query string) (*AgentFilters, error) {
	agentFilters := AgentFilters{}
	if query != "" {
		if err := json.Unmarshal([]byte(query), &agentFilters); err != nil {
			return nil, err
		}
	}

	return &agentFilters, nil
}

func (agentFilters *AgentFilters) BuildPipeline(matchStage bson.M, ctx context.Context, addExtraOperations bool) []bson.M {
	if agentFilters == nil {
		return []bson.M{{"$match": matchStage}}
	}

	pipeline := make([]bson.M, 0)

	// Apply base filters
	baseFilters := &agentFilters.BaseFilters
	baseFilters.ApplyBaseFilters(ctx, matchStage, enums.AGENT)

	// Add match stage
	pipeline = append(pipeline, bson.M{"$match": matchStage})

	// Basic info filters
	if agentFilters.FullName != nil {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"basicInformation.fullName": *agentFilters.FullName}})
	}

	if agentFilters.CompanyName != nil {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"basicInformation.companyName": *agentFilters.CompanyName}})
	}

	if agentFilters.LicenseNumber != nil {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"basicInformation.licenseNumber": *agentFilters.LicenseNumber}})
	}

	if agentFilters.YearsOfExperience != nil {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"basicInformation.yearsOfExperience": *agentFilters.YearsOfExperience}})
	}

	// Business & contact info filters
	if agentFilters.CommercialRegistrationNumber != nil {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"businessInfo.commercialRegistrationNumber": *agentFilters.CommercialRegistrationNumber}})
	}

	if agentFilters.TaxIdentificationNumber != nil {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"businessInfo.taxIdentificationNumber": *agentFilters.TaxIdentificationNumber}})
	}

	if agentFilters.PhoneNumber != nil {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"businessInfo.phoneNumber": *agentFilters.PhoneNumber}})
	}

	if agentFilters.HeadOfficeAddress != nil {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"businessInfo.headOfficeAddress": *agentFilters.HeadOfficeAddress}})
	}

	if agentFilters.OfficialEmail != nil {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"businessInfo.officialEmail": *agentFilters.OfficialEmail}})
	}

	if agentFilters.AuthorizedCountryCodes != nil {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"businessInfo.authorizedCountryCodes": bson.M{"$in": *agentFilters.AuthorizedCountryCodes}}})
	}

	if agentFilters.WebsiteLink != nil {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"businessInfo.websiteLink": *agentFilters.WebsiteLink}})
	}

	// Bank account info filters
	if agentFilters.BankName != nil {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"bankAccountInfo.bankName": *agentFilters.BankName}})
	}

	if agentFilters.AccountNumber != nil {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"bankAccountInfo.accountNumber": *agentFilters.AccountNumber}})
	}

	if agentFilters.AccountHolderName != nil {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"bankAccountInfo.accountHolderName": *agentFilters.AccountHolderName}})
	}

	if agentFilters.Iban != nil {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"bankAccountInfo.iban": *agentFilters.Iban}})
	}

	if agentFilters.SwiftCode != nil {
		pipeline = append(pipeline, bson.M{"$match": bson.M{"bankAccountInfo.swiftCode": *agentFilters.SwiftCode}})
	}

	if addExtraOperations {
		pipeline = base_filters.BuildPopulatePipeline(agentFilters.PopFields, enums.AGENT, pipeline)

		pipeline = base_filters.BuildSortingPipeline(baseFilters.Sorts, pipeline)
	}

	return pipeline
}
