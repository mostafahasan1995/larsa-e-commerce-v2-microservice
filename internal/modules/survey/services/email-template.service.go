package survey_services

import (
	"bytes"
	"fmt"
	"html/template"
	coupon_models "larsa-e-commerce-microservice-v2/internal/modules/coupon/models"
	product_enums "larsa-e-commerce-microservice-v2/internal/modules/product/enums"
	survey_dtos "larsa-e-commerce-microservice-v2/internal/modules/survey/dtos"
	"path/filepath"

	"github.com/samber/do"
)

type SurveyEmailTemplateService interface {
	GenerateSurveyCouponEmail(dto *survey_dtos.CreateFeedbackDTO, coupon *coupon_models.CouponResponse, language string) (string, error)
}

type surveyEmailTemplateService struct {
	templatesPath string
}

func NewSurveyEmailTemplateService(i *do.Injector) (SurveyEmailTemplateService, error) {
	return &surveyEmailTemplateService{
		templatesPath: "templates/email",
	}, nil
}

func (e *surveyEmailTemplateService) GenerateSurveyCouponEmail(dto *survey_dtos.CreateFeedbackDTO, coupon *coupon_models.CouponResponse, language string) (string, error) {
	templateFile := "coupon-email-en.html"
	if language == "ar" {
		templateFile = "coupon-email-ar.html"
	}
	templatePath := filepath.Join(e.templatesPath, templateFile)

	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to parse coupon template: %v", err)
	}

	discountText := e.formatDiscountText(coupon, language)

	data := struct {
		ClientName   string
		Rating       float64
		DiscountText string
		CouponCode   string
	}{
		ClientName:   dto.ClientName,
		Rating:       dto.Points,
		DiscountText: discountText,
		CouponCode:   coupon.Code,
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return "", fmt.Errorf("failed to execute coupon template: %v", err)
	}

	return buf.String(), nil
}

func (e *surveyEmailTemplateService) formatDiscountText(coupon *coupon_models.CouponResponse, language string) string {
	discountValue := coupon.Discount.Value
	discountType := coupon.Discount.Type

	if discountType == product_enums.PERCENTAGE {
		if language == "ar" {
			return fmt.Sprintf("خصم %.0f%%", discountValue)
		}
		return fmt.Sprintf("%.0f%% OFF", discountValue)
	}

	// For fixed amount discounts
	currencySymbol := "$"
	if coupon.Discount.CurrencyCode != nil {
		switch *coupon.Discount.CurrencyCode {
		case "USD":
			currencySymbol = "$"
		case "EUR":
			currencySymbol = "€"
		default:
			currencySymbol = *coupon.Discount.CurrencyCode + " "
		}
	}

	if language == "ar" {
		return fmt.Sprintf("خصم %s%.0f", currencySymbol, discountValue)
	}
	return fmt.Sprintf("%s%.0f OFF", currencySymbol, discountValue)
}
