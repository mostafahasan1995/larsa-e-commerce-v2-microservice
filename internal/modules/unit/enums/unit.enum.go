package unit_enums

type UnitType string

const (
	BiggerThanReference  UnitType = "bigger_than_reference"
	Reference            UnitType = "reference"
	SmallerThanReference UnitType = "smaller_than_reference"
)

type UnitCategory string

const (
	LENGTH       UnitCategory = "length"
	UNIT         UnitCategory = "unit"
	WEIGHT       UnitCategory = "weight"
	VOLUME       UnitCategory = "volume"
	WORKING_TIME UnitCategory = "working_time"
	AREA         UnitCategory = "area"
)
