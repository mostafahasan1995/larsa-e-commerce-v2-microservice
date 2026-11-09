package config_enums

type CategoryVisibility string

const (
	CATEGORIES_LIST     CategoryVisibility = "categories_list"
	SEARCH_RESULT       CategoryVisibility = "search_result"
	CATEGORIES_CAROUSEL CategoryVisibility = "categories_carousel"
	INVISIBLE           CategoryVisibility = "invisible"
)

type StatusName string

const (
	APPROVED StatusName = "approved"
	REJECTED StatusName = "rejected"
	PENDING  StatusName = "pending"
)
