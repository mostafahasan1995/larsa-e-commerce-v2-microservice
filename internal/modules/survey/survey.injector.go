package survey_init

import (
	survey_handlers "larsa-e-commerce-microservice-v2/internal/modules/survey/handlers"
	survey_repos "larsa-e-commerce-microservice-v2/internal/modules/survey/repos"
	survey_services "larsa-e-commerce-microservice-v2/internal/modules/survey/services"

	"github.com/go-chi/chi/v5"
	"github.com/samber/do"
)

func Init(injector *do.Injector, r *chi.Mux) *do.Injector {
	// repos
	do.Provide(injector, survey_repos.NewWaiterRepository)
	do.Provide(injector, survey_repos.NewQuestionRepository)
	do.Provide(injector, survey_repos.NewFeedbackRepository)

	// services
	do.Provide(injector, survey_services.NewWaiterService)
	do.Provide(injector, survey_services.NewQuestionService)
	do.Provide(injector, survey_services.NewFeedbackService)
	do.Provide(injector, survey_services.NewSurveyEmailTemplateService)

	// web services
	do.Provide(injector, survey_services.NewWebWaiterService)
	do.Provide(injector, survey_services.NewWebQuestionService)
	do.Provide(injector, survey_services.NewWebFeedbackService)

	// admin handlers
	survey_handlers.NewWaiterHandler(injector, r)
	survey_handlers.NewQuestionHandler(injector, r)
	survey_handlers.NewFeedbackHandler(injector, r)

	// web handlers
	survey_handlers.NewWebWaiterHandler(injector, r)
	survey_handlers.NewWebQuestionHandler(injector, r)
	survey_handlers.NewWebFeedbackHandler(injector, r)

	return injector
}
