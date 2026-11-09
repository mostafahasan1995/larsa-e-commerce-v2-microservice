package survey_utils

import (
	"context"
	"fmt"
	"strings"

	"larsa-e-commerce-microservice-v2/internal/common/helpers"
	excel_services "larsa-e-commerce-microservice-v2/internal/common/services/excel"
	survey_models "larsa-e-commerce-microservice-v2/internal/modules/survey/models"
)

// GetFeedbackExportColumns returns the column definitions for feedback export
func GetFeedbackExportColumns() []excel_services.ExcelColumn {
	return []excel_services.ExcelColumn{
		excel_services.NewRequiredStringColumn("clientName", "Client Name", excel_services.IntPtr(120)),
		excel_services.NewRequiredStringColumn("email", "Email", excel_services.IntPtr(200)),
		excel_services.NewOptionalStringColumn("phone", "Phone", "", excel_services.IntPtr(150)),
		excel_services.NewOptionalStringColumn("birthday", "Birthday", "", excel_services.IntPtr(100)),
		excel_services.NewFloatColumn("points", "Points", false, excel_services.Float64Ptr(0), nil),
		excel_services.NewFloatColumn("waiterRate", "Waiter Rate", false, excel_services.Float64Ptr(0), nil),
		excel_services.NewOptionalStringColumn("questionAnswers", "Question & Answers", "", excel_services.IntPtr(500)),
		excel_services.NewOptionalStringColumn("answersSummary", "Answers Summary", "", excel_services.IntPtr(300)),
		excel_services.NewOptionalStringColumn("notes", "Notes", "", excel_services.IntPtr(300)),
		excel_services.NewOptionalStringColumn("waiterName", "Waiter Name", "", excel_services.IntPtr(150)),
		excel_services.NewOptionalStringColumn("branchName", "Branch Name", "", excel_services.IntPtr(150)),
		excel_services.NewDateTimeColumn("createdAt", "Created Date", true, "2006-01-02 15:04:05"),
		excel_services.NewDateTimeColumn("updatedAt", "Updated Date", true, "2006-01-02 15:04:05"),
	}
}

// ConvertFeedbacksToExcelData converts feedback responses to Excel export data
func ConvertFeedbacksToExcelData(ctx context.Context, feedbacks []survey_models.FeedbackResponse) ([]map[string]interface{}, *helpers.ApiError) {
	excelData := make([]map[string]interface{}, len(feedbacks))

	for i, feedback := range feedbacks {
		data := map[string]interface{}{
			"clientName": feedback.ClientName,
			"email":      feedback.Email,
			"phone":      feedback.Phone,
			"birthday":   feedback.Birthday,
			"points":     feedback.Points,
			"waiterRate": feedback.WaiterRate,
			"notes":      feedback.Notes,
			"createdAt":  feedback.CreatedAt.Format("2006-01-02 15:04:05"),
			"updatedAt":  feedback.UpdatedAt.Format("2006-01-02 15:04:05"),
		}

		// Process answers with question context
		questionAnswers, answersSummary := formatAnswersForExport(ctx, feedback)
		data["questionAnswers"] = questionAnswers
		data["answersSummary"] = answersSummary

		// Add waiter name if available (handle localized field)
		if feedback.Waiter != nil {
			waiterName := feedback.Waiter.Name.GetValue(ctx)
			data["waiterName"] = waiterName
		} else {
			data["waiterName"] = ""
		}

		// Add branch name if available (handle localized field)
		if feedback.Branch != nil {
			branchName := feedback.Branch.Name.GetValue(ctx)
			data["branchName"] = branchName
		} else {
			data["branchName"] = ""
		}

		excelData[i] = data
	}

	return excelData, nil
}

// formatAnswersForExport formats answers for Excel export with question context
func formatAnswersForExport(ctx context.Context, feedback survey_models.FeedbackResponse) (string, string) {
	if len(feedback.Answers) == 0 {
		return "", ""
	}

	var questionAnswerPairs []string
	var answerTexts []string

	// Create a map of question ID to question for quick lookup
	questionMap := make(map[string]survey_models.QuestionResponse)
	if feedback.Questions != nil {
		for _, question := range feedback.Questions {
			questionMap[question.Id.Hex()] = question
		}
	}

	for _, answer := range feedback.Answers {
		answerTexts = append(answerTexts, answer.Answer)

		// If we have question data, create question-answer pairs
		if question, exists := questionMap[answer.QuestionId.Hex()]; exists {
			questionText := question.Question.GetValue(ctx)
			if questionText != "" {
				pair := fmt.Sprintf("Q: %s | A: %s", questionText, answer.Answer)
				questionAnswerPairs = append(questionAnswerPairs, pair)
			} else {
				pair := fmt.Sprintf("ID: %s | A: %s", answer.QuestionId.Hex()[:8], answer.Answer)
				questionAnswerPairs = append(questionAnswerPairs, pair)
			}
		} else {
			// Fallback to question ID if question data not available
			pair := fmt.Sprintf("ID: %s | A: %s", answer.QuestionId.Hex()[:8], answer.Answer)
			questionAnswerPairs = append(questionAnswerPairs, pair)
		}
	}

	questionAnswers := strings.Join(questionAnswerPairs, " || ")
	answersSummary := strings.Join(answerTexts, ", ")

	return questionAnswers, answersSummary
}
