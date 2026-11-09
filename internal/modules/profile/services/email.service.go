package profile_services

import (
	"bytes"
	"context"
	"html/template"
	"path/filepath"

	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/types"
)

// EmailTemplateData represents the data structure for email templates
type EmailTemplateData struct {
	FirstName string
	Email     string
	Password  string
}

// sendWelcomeEmail sends a welcome email to newly created or existing users
func (s *profileService) sendWelcomeEmail(ctx context.Context, email, firstName, password, language string) {
	subjectTexts := types.Localizable[string]{
		enums.EN: "Welcome to Our Community - Account Created!",
		enums.AR: "مرحبًا بك في مجتمعنا - تم إنشاء الحساب!",
	}

	languageKey := enums.LanguageKey(language)
	subject := types.LocalizeValue(subjectTexts, languageKey)

	htmlBody := s.generateWelcomeEmailHTML(firstName, email, password, language)
	msg := &types.Message{
		Email:       email,
		Subject:     subject,
		MessageHtml: htmlBody,
	}

	s.emailgw.SendEmail(ctx, msg)
}

// loadEmailTemplate loads and renders an email template
func (s *profileService) loadEmailTemplate(templatePath string, data EmailTemplateData) (string, error) {
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// generateWelcomeEmailHTML creates HTML for welcome email using templates
func (s *profileService) generateWelcomeEmailHTML(firstName, email, password, language string) string {
	templateName := "welcome-email-en.html"
	if language == "ar" {
		templateName = "welcome-email-ar.html"
	}

	templatePath := filepath.Join("templates/email", templateName)
	data := EmailTemplateData{
		FirstName: firstName,
		Email:     email,
		Password:  password,
	}

	html, _ := s.loadEmailTemplate(templatePath, data)

	return html
}
