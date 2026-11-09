package api

import (
	// "context"
	"fmt"
	"larsa-e-commerce-microservice-v2/internal/common/configuration"
	"larsa-e-commerce-microservice-v2/internal/common/cron"
	"larsa-e-commerce-microservice-v2/internal/common/enums"
	"larsa-e-commerce-microservice-v2/internal/common/logger"
	middlewares "larsa-e-commerce-microservice-v2/internal/common/middleware"
	"larsa-e-commerce-microservice-v2/internal/common/nats"
	excel_services "larsa-e-commerce-microservice-v2/internal/common/services/excel"
	"larsa-e-commerce-microservice-v2/internal/common/types"
	"larsa-e-commerce-microservice-v2/internal/gateway"
	config_services "larsa-e-commerce-microservice-v2/internal/modules/config/services"

	addressInjector "larsa-e-commerce-microservice-v2/internal/modules/address"
	appointmentInjector "larsa-e-commerce-microservice-v2/internal/modules/appointment"
	callCenterInjector "larsa-e-commerce-microservice-v2/internal/modules/call-center"
	cartInjector "larsa-e-commerce-microservice-v2/internal/modules/carts"
	configInjector "larsa-e-commerce-microservice-v2/internal/modules/config"
	contactInjector "larsa-e-commerce-microservice-v2/internal/modules/contact"
	couponInjector "larsa-e-commerce-microservice-v2/internal/modules/coupon"
	featureInjector "larsa-e-commerce-microservice-v2/internal/modules/feature"
	loyaltyInjector "larsa-e-commerce-microservice-v2/internal/modules/loyalty"
	loyaltySettingsInjector "larsa-e-commerce-microservice-v2/internal/modules/loyalty-settings"
	orderInjector "larsa-e-commerce-microservice-v2/internal/modules/order"
	paymentInjector "larsa-e-commerce-microservice-v2/internal/modules/payment"
	planInjector "larsa-e-commerce-microservice-v2/internal/modules/plan"
	productInjector "larsa-e-commerce-microservice-v2/internal/modules/product"
	profileInjector "larsa-e-commerce-microservice-v2/internal/modules/profile"
	promotionInjector "larsa-e-commerce-microservice-v2/internal/modules/promotion"
	quotationInjector "larsa-e-commerce-microservice-v2/internal/modules/quotation"
	recipeInjector "larsa-e-commerce-microservice-v2/internal/modules/recipe"
	seedInjector "larsa-e-commerce-microservice-v2/internal/modules/seed"
	settingsGeneralInjector "larsa-e-commerce-microservice-v2/internal/modules/settings/general"
	settingsInventoryInjector "larsa-e-commerce-microservice-v2/internal/modules/settings/inventory"
	settingsSalesInjector "larsa-e-commerce-microservice-v2/internal/modules/settings/sales"
	statsInjector "larsa-e-commerce-microservice-v2/internal/modules/stats"
	surveyInjector "larsa-e-commerce-microservice-v2/internal/modules/survey"
	ticketInjector "larsa-e-commerce-microservice-v2/internal/modules/ticket"
	unitInjector "larsa-e-commerce-microservice-v2/internal/modules/unit"
	unservedDeliveryRequestInjector "larsa-e-commerce-microservice-v2/internal/modules/unserved-delivery-request"
	watchlistInjector "larsa-e-commerce-microservice-v2/internal/modules/watchlist"
	"log"
	"net/http"
	"reflect"
	"strings"

	"git.larsa.io/mahdawi/microservices-commons.git/common"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-playground/validator/v10"
	"github.com/samber/do"
)

func Start() error {

	Config := configuration.EnvironmentVariables

	fmt.Println("start server on port:", Config.Port)
	logger.Init(Config.Mode) // Initialize zerolog with mode
	injector := do.New()
	conn, err := configuration.InitDB()
	if err != nil {
		return err
	}
	configuration.InitRedis()
	configuration.InitNats()
	enums.InitCapabilities()

	r := chi.NewRouter()
	validateInstance := validator.New()
	validateInstance.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}
		return name
	})
	// Register custom validators
	types.RegisterValidators(validateInstance)

	// Replace default logger with Zerolog middleware
	// r.Use(middleware.Logger)
	r.Use(middlewares.ZerologRequestLogger)

	r.Use(middleware.Recoverer) // Add panic recovery
	//r.Use(middlewares.RequestTiming)
	r.Use(middleware.Compress(5))
	r.Use(middlewares.AdminCache()) // Add Redis caching for admin endpoints AFTER compression
	r.Use(middleware.SetHeader("Content-Type", "application/json"))

	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "x-client", "x-access-token", "x-service", "x-expire", "x-service-token", "x-user-id", "Content-Disposition", "x-Access-Token", "x-Refresh-Token"},
		ExposedHeaders:   []string{"Link", "Content-Disposition", "X-Access-Token", "X-Refresh-Token"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	r.Use(common.ClientSubscriptionInstance.CheckSubscription)

	r.Use(middlewares.AddLangToContext)
	r.Use(middlewares.AddQueryKeysToContext)
	r.Use(middlewares.AddInjectorToContext(injector))

	//database
	do.ProvideValue(injector, conn)
	do.Provide(injector, configuration.NewWithTxn)

	// validator
	do.ProvideValue(injector, validateInstance)
	//chat.Init(injector, r)
	//message.Init(injector, r)
	//whatsapp.Init(injector, r)
	//leads

	//gateways
	do.Provide(injector, gateway.NewUsersGw)
	do.Provide(injector, gateway.NewRolesGw)
	do.Provide(injector, gateway.NewPaymentGw)
	do.Provide(injector, gateway.NewOptionsGw)
	do.Provide(injector, gateway.NewEmailGw)

	// Common services
	do.Provide(injector, excel_services.NewExcelService)

	// Modules
	appointmentInjector.Init(injector, r)
	configInjector.Init(injector, r)
	featureInjector.Init(injector, r)
	profileInjector.Init(injector, r)
	loyaltySettingsInjector.RegisterLoyaltySettingsModule(injector, r)
	loyaltyInjector.Init(injector, r)
	planInjector.Init(injector, r)
	unitInjector.Init(injector, r)
	settingsSalesInjector.Init(injector, r)
	settingsInventoryInjector.Init(injector, r)
	settingsGeneralInjector.Init(injector, r)
	unservedDeliveryRequestInjector.Init(injector, r)
	productInjector.Init(injector, r)
	callCenterInjector.Init(injector, r) // Call-center after product (depends on BranchWebService, WebProductService)
	orderInjector.Init(injector, r)

	// Dont change order to above DEPENDS ON ORDER MODULE
	statsInjector.Init(injector, r)
	seedInjector.Init(injector, r)
	contactInjector.Init(injector, r)
	paymentInjector.Init(injector, r)
	quotationInjector.Init(injector, r)
	addressInjector.Init(injector, r)

	// Initialize cart and coupon modules together to resolve circular dependency
	// First, provide all services from both modules
	cartInjector.ProvideServices(injector)
	couponInjector.ProvideServices(injector)

	// Then, initialize handlers (which can now safely invoke dependencies)
	cartInjector.InitHandlers(injector, r)
	couponInjector.InitHandlers(injector, r)

	watchlistInjector.Init(injector, r)
	surveyInjector.Init(injector, r)
	ticketInjector.Init(injector, r)
	recipeInjector.Init(injector, r)
	promotionInjector.Init(injector, r)
	// seedInjector.Init(injector, r)

	// Initialize cron services
	cron.Init(injector)

	nats.InitListeners(injector)

	/*
	 * Upsert settings and indexes
	 */
	settingService := do.MustInvoke[config_services.SettingService](injector)
	settingService.Upsert()

	fmt.Println("start server")
	fmt.Println(Config.Port)

	if err := http.ListenAndServe(":"+Config.Port, r); err != nil {
		log.Fatal(err)
	}

	return nil
}
