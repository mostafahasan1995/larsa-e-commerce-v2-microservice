package enums

type CollectionName string

const (
	PRODUCT            CollectionName = "products"
	PRICE              CollectionName = "prices"
	PRODUCTION_LINE    CollectionName = "productionLines"
	CATEGORY           CollectionName = "categories"
	CAPABILITY         CollectionName = "capabilities"
	UNIT               CollectionName = "units"
	BRANCH             CollectionName = "branches"
	SEED               CollectionName = "seeds"
	ATTRIBUTE          CollectionName = "attributes"
	GROUP              CollectionName = "groups"
	DISCOUNT           CollectionName = "discounts"
	ADDON              CollectionName = "addons"
	SHIPPING_COMPANY   CollectionName = "shippingCompanies"
	BROKER             CollectionName = "brokers"
	AGENT              CollectionName = "agents"
	CUSTOMER           CollectionName = "customers"
	PROFILE            CollectionName = "profiles"
	FEATURE            CollectionName = "features"
	FEATURE_GROUP      CollectionName = "feature_groups"
	PLAN               CollectionName = "plans"
	CURRENCY           CollectionName = "currencies"
	BANNER             CollectionName = "banners"
	COUPON             CollectionName = "coupons"
	USER_COUPON        CollectionName = "user_coupons"
	PAYMENT            CollectionName = "payments"
	INVOICE            CollectionName = "invoices"
	TICKET             CollectionName = "tickets"
	TICKET_TYPE        CollectionName = "ticket_types"
	PAYOUT             CollectionName = "payouts"
	DOCUMENT           CollectionName = "documents"
	QUOTATION          CollectionName = "quotations"
	QUOTATION_SETTINGS CollectionName = "quotation_settings"
	SHIPPING           CollectionName = "shippings"
	QUOTATION_SUMMARY  CollectionName = "quotation_summaries"
	TRANSACTION        CollectionName = "transactions"
	APPOINTMENT        CollectionName = "appointments"
	ORDER              CollectionName = "orders"
	SETTING            CollectionName = "settings"
	MONITERING_LOG     CollectionName = "monitering_logs"
	BRANCH_PRODUCT     CollectionName = "branch_products"
	ZONE               CollectionName = "zones"
	COUNTRY            CollectionName = "countries"
	SEED_COUNTRY       CollectionName = "seed_countries"
	SEED_CITY          CollectionName = "seed_cities"
	SEED_AREA          CollectionName = "seed_areas"
	AREA               CollectionName = "areas"
	CITY               CollectionName = "cities"
	CART_ITEM          CollectionName = "cart_items"
	CHECKOUT           CollectionName = "checkouts"
	CART               CollectionName = "carts"
	WATCH_LIST         CollectionName = "watch_lists"
	SALES_SETTINGS     CollectionName = "sales-settings"
	INVENTORY_SETTINGS CollectionName = "inventory-settings"
	GENERAL_SETTINGS   CollectionName = "general-settings"
	ADDRESS            CollectionName = "addresses"
	QUESTION           CollectionName = "questions"
	WAITER             CollectionName = "waiters"
	FEEDBACK           CollectionName = "feedbacks"
	REVIEW             CollectionName = "reviews"
	RECIPE             CollectionName = "recipes"
	PROMOTION          CollectionName = "promotions"
	LOYALTY_POINTS     CollectionName = "loyalty_points"
	LOYALTY_SETTINGS   CollectionName = "loyalty_settings"
	// B2B
	PLAN_FEATURES CollectionName = "plan_features"
	STATS         CollectionName = "stats"
	UNSERVED_DELIVERY_REQUEST CollectionName = "unserved_delivery_requests"
)

const (
	PRODUCT_MODEL_NAME            = "ecommerceV2Products"
	DISCOUNT_MODEL_NAME           = "ecommerceV2Discounts"
	PRICE_MODEL_NAME              = "ecommerceV2Prices"
	PRODUCTION_LINE_MODEL_NAME    = "ecommerceV2ProductionLines"
	CATEGORY_MODEL_NAME           = "ecommerceV2Categories"
	COUNTRY_MODEL_NAME            = "ecommerceV2Countries"
	UNIT_MODEL_NAME               = "ecommerceV2Units"
	ATTRIBUTE_MODEL_NAME          = "ecommerceV2Attributes"
	SETTING_MODEL_NAME            = "ecommerceV2Settings"
	GROUP_MODEL_NAME              = "ecommerceV2Groups"
	REVIEW_MODEL_NAME             = "ecommerceV2Reviews"
	SHIPPING_COMPANY_MODEL_NAME   = "ecommerceV2ShippingCompanies"
	BROKER_MODEL_NAME             = "ecommerceV2Brokers"
	AGENT_MODEL_NAME              = "ecommerceV2Agents"
	CUSTOMER_MODEL_NAME           = "ecommerceV2Customers"
	PROFILE_MODEL_NAME            = "ecommerceV2Profiles"
	BRANCH_MODEL_NAME             = "ecommerceV2Branches"
	SEED_MODEL_NAME               = "ecommerceV2Seeds"
	FEATURE_MODEL_NAME            = "ecommerceV2Features"
	FEATURE_GROUP_MODEL_NAME      = "ecommerceV2FeatureGroups"
	PLAN_MODEL_NAME               = "ecommerceV2Plans"
	CURRENCY_MODEL_NAME           = "ecommerceV2Currencies"
	PAYMENT_MODEL_NAME            = "ecommerceV2Payments"
	INVOICE_MODEL_NAME            = "ecommerceV2Invoices"
	TRANSACTION_MODEL_NAME        = "ecommerceV2Transactions"
	PAYOUT_MODEL_NAME             = "ecommerceV2Payouts"
	DOCUMENT_MODEL_NAME           = "ecommerceV2Documents"
	QUOTATION_MODEL_NAME          = "ecommerceV2Quotations"
	QUOTATION_SETTINGS_MODEL_NAME = "ecommerceV2QuotationSettings"
	SHIPPING_MODEL_NAME           = "ecommerceV2Shippings"
	QUOTATION_SUMMARY_MODEL_NAME  = "ecommerceV2QuotationSummaries"
	APPOINTMENT_MODEL_NAME        = "ecommerceV2Appointments"
	ORDER_MODEL_NAME              = "ecommerceV2Orders"
	MONITERING_LOG_MODEL_NAME     = "ecommerceV2MoniteringLogs"
	ADDON_MODEL_NAME              = "ecommerceV2Addons"
	BRANCH_PRODUCT_MODEL_NAME     = "ecommerceV2BranchesProducts"
	ZONE_MODEL_NAME               = "ecommerceV2Zones"
	AREA_MODEL_NAME               = "ecommerceV2Areas"
	CITY_MODEL_NAME               = "ecommerceV2Cities"
	SEED_COUNTRY_MODEL_NAME       = "ecommerceV2SeedCountries"
	SEED_CITY_MODEL_NAME          = "ecommerceV2SeedCities"
	SEED_AREA_MODEL_NAME          = "ecommerceV2SeedAreas"
	CART_ITEM_MODEL_NAME          = "ecommerceV2CartItems"
	CART_MODEL_NAME               = "ecommerceV2Carts"
	CHECKOUT_MODEL_NAME           = "ecommerceV2Checkouts"
	WATCHLIST_MODEL_NAME          = "ecommerceV2WatchLists"
	SALES_SETTINGS_MODEL_NAME     = "ecommerceV2SalesSettings"
	INVENTORY_SETTINGS_MODEL_NAME = "ecommerceV2InventorySettings"
	GENERAL_SETTINGS_MODEL_NAME   = "ecommerceV2GeneralSettings"
	ADDRESS_MODEL_NAME            = "ecommerceV2Addresses"
	COUPON_MODEL_NAME             = "ecommerceV2Coupons"
	USER_COUPON_MODEL_NAME        = "ecommerceV2UserCoupons"
	QUESTION_MODEL_NAME           = "ecommerceV2Questions"
	WAITER_MODEL_NAME             = "ecommerceV2Waiters"
	FEEDBACK_MODEL_NAME           = "ecommerceV2Feedbacks"
	TICKET_TYPE_MODEL_NAME        = "ecommerceV2TicketTypes"
	RECIPE_MODEL_NAME             = "ecommerceV2Recipes"
	PROMOTION_MODEL_NAME          = "ecommerceV2Promotions"
	PLAN_FEATURES_MODEL_NAME      = "ecommerceV2PlanFeatures"
	BANNER_MODEL_NAME                       = "ecommerceV2Banners"
	UNSERVED_DELIVERY_REQUEST_MODEL_NAME    = "ecommerceV2UnservedDeliveryRequests"

)

func (i CollectionName) GetCollectionName() string {
	switch i {
	case BANNER:
		return BANNER_MODEL_NAME
	case PRODUCT:
		return PRODUCT_MODEL_NAME
	case DISCOUNT:
		return DISCOUNT_MODEL_NAME
	case REVIEW:
		return REVIEW_MODEL_NAME
	case PRICE:
		return PRICE_MODEL_NAME
	case PRODUCTION_LINE:
		return PRODUCTION_LINE_MODEL_NAME
	case COUPON:
		return COUPON_MODEL_NAME
	case USER_COUPON:
		return USER_COUPON_MODEL_NAME
	case CATEGORY:
		return CATEGORY_MODEL_NAME
	case COUNTRY:
		return COUNTRY_MODEL_NAME
	case UNIT:
		return UNIT_MODEL_NAME
	case ATTRIBUTE:
		return ATTRIBUTE_MODEL_NAME
	case SETTING:
		return SETTING_MODEL_NAME
	case SHIPPING_COMPANY:
		return SHIPPING_COMPANY_MODEL_NAME
	case BROKER:
		return BROKER_MODEL_NAME
	case BRANCH:
		return BRANCH_MODEL_NAME
	case AGENT:
		return AGENT_MODEL_NAME
	case CUSTOMER:
		return CUSTOMER_MODEL_NAME
	case PROFILE:
		return PROFILE_MODEL_NAME
	case FEATURE:
		return FEATURE_MODEL_NAME
	case FEATURE_GROUP:
		return FEATURE_GROUP_MODEL_NAME
	case PLAN:
		return PLAN_MODEL_NAME
	case CURRENCY:
		return CURRENCY_MODEL_NAME
	case PAYMENT:
		return PAYMENT_MODEL_NAME
	case INVOICE:
		return INVOICE_MODEL_NAME
	case TRANSACTION:
		return TRANSACTION_MODEL_NAME
	case PAYOUT:
		return PAYOUT_MODEL_NAME
	case DOCUMENT:
		return DOCUMENT_MODEL_NAME
	case QUOTATION:
		return QUOTATION_MODEL_NAME
	case QUOTATION_SETTINGS:
		return QUOTATION_SETTINGS_MODEL_NAME
	case SHIPPING:
		return SHIPPING_MODEL_NAME
	case QUOTATION_SUMMARY:
		return QUOTATION_SUMMARY_MODEL_NAME
	case APPOINTMENT:
		return APPOINTMENT_MODEL_NAME
	case ORDER:
		return ORDER_MODEL_NAME
	case GROUP:
		return GROUP_MODEL_NAME
	case MONITERING_LOG:
		return MONITERING_LOG_MODEL_NAME
	case ADDON:
		return ADDON_MODEL_NAME
	case BRANCH_PRODUCT:
		return BRANCH_PRODUCT_MODEL_NAME
	case ZONE:
		return ZONE_MODEL_NAME
	case AREA:
		return AREA_MODEL_NAME
	case CITY:
		return CITY_MODEL_NAME
	case SEED_COUNTRY:
		return SEED_COUNTRY_MODEL_NAME
	case SEED_CITY:
		return SEED_CITY_MODEL_NAME
	case SEED_AREA:
		return SEED_AREA_MODEL_NAME
	case CART_ITEM:
		return CART_ITEM_MODEL_NAME
	case CHECKOUT:
		return CHECKOUT_MODEL_NAME
	case CART:
		return CART_MODEL_NAME
	case WATCH_LIST:
		return WATCHLIST_MODEL_NAME
	case SALES_SETTINGS:
		return SALES_SETTINGS_MODEL_NAME
	case INVENTORY_SETTINGS:
		return INVENTORY_SETTINGS_MODEL_NAME
	case GENERAL_SETTINGS:
		return GENERAL_SETTINGS_MODEL_NAME
	case ADDRESS:
		return ADDRESS_MODEL_NAME
	case QUESTION:
		return QUESTION_MODEL_NAME
	case WAITER:
		return WAITER_MODEL_NAME
	case FEEDBACK:
		return FEEDBACK_MODEL_NAME
	case TICKET_TYPE:
		return TICKET_TYPE_MODEL_NAME
	case RECIPE:
		return RECIPE_MODEL_NAME
	case PROMOTION:
		return PROMOTION_MODEL_NAME
	case PLAN_FEATURES:
		return PLAN_FEATURES_MODEL_NAME
	case UNSERVED_DELIVERY_REQUEST:
		return UNSERVED_DELIVERY_REQUEST_MODEL_NAME
	default:
		return "ecommerceV2Default"
	}
}

func GetCollectionNameList() []string {
	return []string{
		DISCOUNT_MODEL_NAME,
		BANNER_MODEL_NAME,
		PRODUCT_MODEL_NAME,
		REVIEW_MODEL_NAME,
		PRICE_MODEL_NAME,
		PRODUCTION_LINE_MODEL_NAME,
		CATEGORY_MODEL_NAME,
		COUNTRY_MODEL_NAME,
		UNIT_MODEL_NAME,
		ATTRIBUTE_MODEL_NAME,
		SETTING_MODEL_NAME,
		SHIPPING_COMPANY_MODEL_NAME,
		BROKER_MODEL_NAME,
		AGENT_MODEL_NAME,
		CUSTOMER_MODEL_NAME,
		PROFILE_MODEL_NAME,
		BRANCH_MODEL_NAME,
		SEED_MODEL_NAME,
		FEATURE_MODEL_NAME,
		FEATURE_GROUP_MODEL_NAME,
		PLAN_MODEL_NAME,
		CURRENCY_MODEL_NAME,
		PAYMENT_MODEL_NAME,
		INVOICE_MODEL_NAME,
		TRANSACTION_MODEL_NAME,
		PAYOUT_MODEL_NAME,
		DOCUMENT_MODEL_NAME,
		QUOTATION_MODEL_NAME,
		QUOTATION_SETTINGS_MODEL_NAME,
		SHIPPING_MODEL_NAME,
		QUOTATION_SUMMARY_MODEL_NAME,
		APPOINTMENT_MODEL_NAME,
		ORDER_MODEL_NAME,
		GROUP_MODEL_NAME,
		SETTING_MODEL_NAME,
		MONITERING_LOG_MODEL_NAME,
		BRANCH_PRODUCT_MODEL_NAME,
		ZONE_MODEL_NAME,
		AREA_MODEL_NAME,
		CITY_MODEL_NAME,
		SEED_COUNTRY_MODEL_NAME,
		SEED_CITY_MODEL_NAME,
		SEED_AREA_MODEL_NAME,
		CART_ITEM_MODEL_NAME,
		CHECKOUT_MODEL_NAME,
		CART_MODEL_NAME,
		WATCHLIST_MODEL_NAME,
		SALES_SETTINGS_MODEL_NAME,
		INVENTORY_SETTINGS_MODEL_NAME,
		GENERAL_SETTINGS_MODEL_NAME,
		ADDRESS_MODEL_NAME,
		COUPON_MODEL_NAME,
		USER_COUPON_MODEL_NAME,
		TICKET_TYPE_MODEL_NAME,
		PROMOTION_MODEL_NAME,
		PLAN_FEATURES_MODEL_NAME,
		UNSERVED_DELIVERY_REQUEST_MODEL_NAME,
	}
}

type Route string

const (
	// Model routes
	ROUTE_CATEGORIES          Route = "/categories"
	ROUTE_DISCOUNTS           Route = "/discounts"
	ROUTE_COUNTRIES           Route = "/countries"
	ROUTE_BANNERS             Route = "/banners"
	ROUTE_PLANS               Route = "/plans"
	ROUTE_PRICES              Route = "/prices"
	ROUTE_SEEDS               Route = "/seeds"
	ROUTE_UNITS               Route = "/units"
	ROUTE_BRANCHES            Route = "/branches"
	ROUTE_BROKER_AGENTS       Route = "/agents"
	ROUTE_BROKERS             Route = "/brokers"
	ROUTE_TICKETS             Route = "/tickets"
	ROUTE_PAYOUTS             Route = "/payouts"
	ROUTE_INVOICES            Route = "/invoices"
	ROUTE_PAYMENTS            Route = "/payments"
	ROUTE_PRODUCTS            Route = "/products"
	ROUTE_BROKER_CUSTOMERS    Route = "/customers"
	ROUTE_SHIPPING_COMPANIES  Route = "/shipping-companies"
	ROUTE_ATTRIBUTES          Route = "/attributes"
	ROUTE_GROUPS              Route = "/groups"
	ROUTE_CAPABILITIES        Route = "/capabilities"
	ROUTE_PRODUCTION_LINES    Route = "/production-lines"
	ROUTE_FEATURES            Route = "/features"
	ROUTE_FEATURE_GROUPS      Route = "/feature-groups"
	ROUTE_PROFILES            Route = "/profiles"
	ROUTE_CURRENCIES          Route = "/currencies"
	ROUTE_DOCUMENTS           Route = "/documents"
	ROUTE_QUOTATIONS          Route = "/quotations"
	ROUTE_QUOTATION_SETTINGS  Route = "/quotation-settings"
	ROUTE_SHIPPINGS           Route = "/shippings"
	ROUTE_QUOTATION_SUMMARIES Route = "/quotation-summaries"
	ROUTE_TRANSACTIONS        Route = "/transactions"
	ROUTE_APPOINTMENTS        Route = "/appointments"
	ROUTE_ORDERS              Route = "/orders"
	ROUTE_ADDONS              Route = "/addons"
	ROUTE_ZONES               Route = "/zones"
	ROUTE_AREAS               Route = "/areas"
	ROUTE_CITIES              Route = "/cities"
	ROUTE_BRANCH_PRODUCTS     Route = "/branch-products"
	ROUTE_CART_ITEMS          Route = "/carts"
	ROUTE_CHECKOUTS           Route = "/checkouts"
	ROUTE_SETTINGS            Route = "/settings"
	ROUTE_WATCH_LISTS         Route = "/watch-lists"
	ROUTE_SETTING_GENERAL     Route = "/settings/general"
	ROUTE_SETTING_SALES       Route = "/settings/sales"
	ROUTE_SETTING_INVENTORY   Route = "/settings/inventory"
	ROUTE_ADDRESSES           Route = "/addresses"
	ROUTE_COUPONS             Route = "/coupons"
	ROUTE_QUESTIONS           Route = "/questions"
	ROUTE_WAITERS             Route = "/waiters"
	ROUTE_FEEDBACKS           Route = "/feedbacks"
	ROUTE_TICKET_TYPES        Route = "/ticket-types"
	ROUTE_REVIEWS             Route = "/reviews"
	ROUTE_RECIPES             Route = "/recipes"
	ROUTE_PROMOTIONS          Route = "/promotions"
	ROUTE_PLAN_FEATURES       Route = "/plan-features"
	ROUTE_LOYALTY_POINTS      Route = "/loyalty-points"
	ROUTE_LOYALTY_SETTINGS    Route = "/loyalty-settings"
	ROUTE_STATS               Route = "/stats"
	ROUTE_UNSERVED_DELIVERY_REQUESTS Route = "/unserved-delivery-requests"
)

// GetAdminRoute returns the full admin route path for a given route
func (i CollectionName) GetVendorRoute() string {
	VendorPrefix := "/vendor"

	buildRoute := func(route Route) string {
		return string(VendorPrefix) + string(route)
	}

	switch i {
	case QUOTATION:
		return buildRoute(ROUTE_QUOTATIONS)
	default:
		return "/"
	}
}

// GetAdminRoute returns the full admin route path for a given route
func (i CollectionName) GetAdminRoute() string {
	AdminPrefix := "/admin"

	buildRoute := func(route Route) string {
		return string(AdminPrefix) + string(route)
	}

	switch i {
	case BANNER:
		return buildRoute(ROUTE_BANNERS)
	case DISCOUNT:
		return buildRoute(ROUTE_DISCOUNTS)
	case CAPABILITY:
		return buildRoute(ROUTE_CAPABILITIES)
	case REVIEW:
		return buildRoute(ROUTE_REVIEWS)
	case PRODUCT:
		return buildRoute(ROUTE_PRODUCTS)
	case PRICE:
		return buildRoute(ROUTE_PRICES)
	case PRODUCTION_LINE:
		return buildRoute(ROUTE_PRODUCTION_LINES)
	case CATEGORY:
		return buildRoute(ROUTE_CATEGORIES)
	case COUNTRY:
		return buildRoute(ROUTE_COUNTRIES)
	case UNIT:
		return buildRoute(ROUTE_UNITS)
	case BRANCH:
		return buildRoute(ROUTE_BRANCHES)
	case SEED:
		return buildRoute(ROUTE_SEEDS)
	case ATTRIBUTE:
		return buildRoute(ROUTE_ATTRIBUTES)
	case GROUP:
		return buildRoute(ROUTE_GROUPS)
	case SHIPPING_COMPANY:
		return buildRoute(ROUTE_SHIPPING_COMPANIES)
	case BROKER:
		return buildRoute(ROUTE_BROKERS)
	case AGENT:
		return buildRoute(ROUTE_BROKER_AGENTS)
	case CUSTOMER:
		return buildRoute(ROUTE_BROKER_CUSTOMERS)
	case PROFILE:
		return buildRoute(ROUTE_PROFILES)
	case FEATURE:
		return buildRoute(ROUTE_FEATURES)
	case FEATURE_GROUP:
		return buildRoute(ROUTE_FEATURE_GROUPS)
	case PLAN:
		return buildRoute(ROUTE_PLANS)
	case CURRENCY:
		return buildRoute(ROUTE_CURRENCIES)
	case PAYMENT:
		return buildRoute(ROUTE_PAYMENTS)
	case INVOICE:
		return buildRoute(ROUTE_INVOICES)
	case TRANSACTION:
		return buildRoute(ROUTE_TRANSACTIONS)
	case PAYOUT:
		return buildRoute(ROUTE_PAYOUTS)
	case DOCUMENT:
		return buildRoute(ROUTE_DOCUMENTS)
	case QUOTATION:
		return buildRoute(ROUTE_QUOTATIONS)
	case QUOTATION_SETTINGS:
		return buildRoute(ROUTE_QUOTATION_SETTINGS)
	case SHIPPING:
		return buildRoute(ROUTE_SHIPPINGS)
	case QUOTATION_SUMMARY:
		return buildRoute(ROUTE_QUOTATION_SUMMARIES)
	case APPOINTMENT:
		return buildRoute(ROUTE_APPOINTMENTS)
	case ORDER:
		return buildRoute(ROUTE_ORDERS)
	case ADDON:
		return buildRoute(ROUTE_ADDONS)
	case ZONE:
		return buildRoute(ROUTE_ZONES)
	case AREA:
		return buildRoute(ROUTE_AREAS)
	case CITY:
		return buildRoute(ROUTE_CITIES)
	case BRANCH_PRODUCT:
		return buildRoute(ROUTE_BRANCH_PRODUCTS)
	case CART_ITEM:
		return buildRoute(ROUTE_CART_ITEMS)
	case CHECKOUT:
		return buildRoute(ROUTE_CHECKOUTS)
	case WATCH_LIST:
		return buildRoute(ROUTE_WATCH_LISTS)
	case SALES_SETTINGS:
		return buildRoute(ROUTE_SETTING_SALES)
	case INVENTORY_SETTINGS:
		return buildRoute(ROUTE_SETTING_INVENTORY)
	case GENERAL_SETTINGS:
		return buildRoute(ROUTE_SETTING_GENERAL)
	case ADDRESS:
		return buildRoute(ROUTE_ADDRESSES)
	case QUESTION:
		return buildRoute(ROUTE_QUESTIONS)
	case WAITER:
		return buildRoute(ROUTE_WAITERS)
	case FEEDBACK:
		return buildRoute(ROUTE_FEEDBACKS)
	case COUPON:
		return buildRoute(ROUTE_COUPONS)
	case TICKET:
		return buildRoute(ROUTE_TICKETS)
	case TICKET_TYPE:
		return buildRoute(ROUTE_TICKET_TYPES)
	case RECIPE:
		return buildRoute(ROUTE_RECIPES)
	case PROMOTION:
		return buildRoute(ROUTE_PROMOTIONS)
	case PLAN_FEATURES:
		return buildRoute(ROUTE_PLAN_FEATURES)
	case LOYALTY_SETTINGS:
		return buildRoute(ROUTE_LOYALTY_SETTINGS)
	case SETTING:
		return buildRoute(ROUTE_SETTINGS)
	case STATS:
		return buildRoute(ROUTE_STATS)
	case UNSERVED_DELIVERY_REQUEST:
		return buildRoute(ROUTE_UNSERVED_DELIVERY_REQUESTS)
	default:
		return "/"
	}
}

// GetWebRoute returns the full web route path for a given route
func (i CollectionName) GetWebRoute() string {
	buildRoute := func(route Route) string {
		return string(route)
	}

	switch i {
	case BANNER:
		return buildRoute(ROUTE_BANNERS)
	case DISCOUNT:
		return buildRoute(ROUTE_DISCOUNTS)
	case PRODUCT:
		return buildRoute(ROUTE_PRODUCTS)
	case REVIEW:
		return buildRoute(ROUTE_REVIEWS)
	case PRODUCTION_LINE:
		return buildRoute(ROUTE_PRODUCTION_LINES)
	case CATEGORY:
		return buildRoute(ROUTE_CATEGORIES)
	case CAPABILITY:
		return buildRoute(ROUTE_CAPABILITIES)
	case UNIT:
		return buildRoute(ROUTE_UNITS)
	case ATTRIBUTE:
		return buildRoute(ROUTE_ATTRIBUTES)
	case SHIPPING_COMPANY:
		return buildRoute(ROUTE_SHIPPING_COMPANIES)
	case BROKER:
		return buildRoute(ROUTE_BROKERS)
	case AGENT:
		return buildRoute(ROUTE_BROKER_AGENTS)
	case CUSTOMER:
		return buildRoute(ROUTE_BROKER_CUSTOMERS)
	case PROFILE:
		return buildRoute(ROUTE_PROFILES)
	case FEATURE:
		return buildRoute(ROUTE_FEATURES)
	case FEATURE_GROUP:
		return buildRoute(ROUTE_FEATURE_GROUPS)
	case PLAN:
		return buildRoute(ROUTE_PLANS)
	case CURRENCY:
		return buildRoute(ROUTE_CURRENCIES)
	case PAYMENT:
		return buildRoute(ROUTE_PAYMENTS)
	case INVOICE:
		return buildRoute(ROUTE_INVOICES)
	case TRANSACTION:
		return buildRoute(ROUTE_TRANSACTIONS)
	case PAYOUT:
		return buildRoute(ROUTE_PAYOUTS)
	case QUOTATION:
		return buildRoute(ROUTE_QUOTATIONS)
	case QUOTATION_SETTINGS:
		return buildRoute(ROUTE_QUOTATION_SETTINGS)
	case QUOTATION_SUMMARY:
		return buildRoute(ROUTE_QUOTATION_SUMMARIES)
	case APPOINTMENT:
		return buildRoute(ROUTE_APPOINTMENTS)
	case ORDER:
		return buildRoute(ROUTE_ORDERS)
	case ADDON:
		return buildRoute(ROUTE_ADDONS)
	case ZONE:
		return buildRoute(ROUTE_ZONES)
	case AREA:
		return buildRoute(ROUTE_AREAS)
	case CITY:
		return buildRoute(ROUTE_CITIES)
	case BRANCH:
		return buildRoute(ROUTE_BRANCHES)
	case CART_ITEM:
		return buildRoute(ROUTE_CART_ITEMS)
	case CHECKOUT:
		return buildRoute(ROUTE_CHECKOUTS)
	case WATCH_LIST:
		return buildRoute(ROUTE_WATCH_LISTS)
	case SALES_SETTINGS:
		return buildRoute(ROUTE_SETTING_SALES)
	case INVENTORY_SETTINGS:
		return buildRoute(ROUTE_SETTING_INVENTORY)
	case GENERAL_SETTINGS:
		return buildRoute(ROUTE_SETTING_GENERAL)
	case ADDRESS:
		return buildRoute(ROUTE_ADDRESSES)
	case COUPON:
		return buildRoute(ROUTE_COUPONS)
	case TICKET:
		return buildRoute(ROUTE_TICKETS)
	case QUESTION:
		return buildRoute(ROUTE_QUESTIONS)
	case WAITER:
		return buildRoute(ROUTE_WAITERS)
	case FEEDBACK:
		return buildRoute(ROUTE_FEEDBACKS)
	case TICKET_TYPE:
		return buildRoute(ROUTE_TICKET_TYPES)
	case RECIPE:
		return buildRoute(ROUTE_RECIPES)
	case PROMOTION:
		return buildRoute(ROUTE_PROMOTIONS)
	case PLAN_FEATURES:
		return buildRoute(ROUTE_PLAN_FEATURES)
	case LOYALTY_POINTS:
		return buildRoute(ROUTE_LOYALTY_POINTS)
	default:
		return "/"
	}
}

func (i CollectionName) GetDocumentName() string {
	switch i {
	case PRODUCT:
		return "Product"
	case DISCOUNT:
		return "Discount"
	case REVIEW:
		return "Review"
	case PRODUCTION_LINE:
		return "Production Line"
	case CATEGORY:
		return "Category"
	case UNIT:
		return "Unit"
	case ATTRIBUTE:
		return "Attribute"
	case CAPABILITY:
		return "Capability"
	case GROUP:
		return "Group"
	case SEED:
		return "Seed"
	case SHIPPING_COMPANY:
		return "Broker Shipping"
	case BROKER:
		return "Broker"
	case AGENT:
		return "Agent"
	case CUSTOMER:
		return "Customer"
	case PROFILE:
		return "Profile"
	case FEATURE:
		return "Feature"
	case FEATURE_GROUP:
		return "Feature Group"
	case PLAN:
		return "Plan"
	case CURRENCY:
		return "Currency"
	case PAYMENT:
		return "Payment"
	case INVOICE:
		return "Invoice"
	case TRANSACTION:
		return "Transaction"
	case PAYOUT:
		return "Payout"
	case DOCUMENT:
		return "Document"
	case APPOINTMENT:
		return "Appointment"
	case ORDER:
		return "Order"
	case SHIPPING:
		return "Shipping"
	case QUOTATION_SETTINGS:
		return "Quotation Settings"
	case QUOTATION_SUMMARY:
		return "Quotation Summary"
	case QUOTATION:
		return "Quotation"
	case COUNTRY:
		return "Country"
	case ADDON:
		return "Addon"
	case ZONE:
		return "Zone"
	case AREA:
		return "Area"
	case CITY:
		return "City"
	case CART_ITEM:
		return "Cart Item"
	case CHECKOUT:
		return "Checkout"
	case SEED_COUNTRY:
		return "Seed Country"
	case SEED_CITY:
		return "Seed City"
	case SEED_AREA:
		return "Seed Area"
	case BRANCH:
		return "Branch"
	case BRANCH_PRODUCT:
		return "Branch Product"
	case WATCH_LIST:
		return "Watchlist"
	case SALES_SETTINGS:
		return "Sales Settings"
	case INVENTORY_SETTINGS:
		return "Inventory Settings"
	case GENERAL_SETTINGS:
		return "General Settings"
	case ADDRESS:
		return "Address"
	case COUPON:
		return "Coupon"
	case USER_COUPON:
		return "User Coupon"
	case CART:
		return "Cart"
	case QUESTION:
		return "Question"
	case WAITER:
		return "Waiter"
	case FEEDBACK:
		return "Feedback"
	case TICKET_TYPE:
		return "Ticket Type"
	case RECIPE:
		return "Recipe"
	case PROMOTION:
		return "Promotion"
	case PLAN_FEATURES:
		return "Plan Feature"
	default:
		return "Document"
	}
}
