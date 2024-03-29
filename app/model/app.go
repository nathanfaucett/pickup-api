package model

type HealthST struct {
	DB bool `json:"db" validate:"required"`
} // @name Health

func (health HealthST) IsHealthy() bool {
	return health.DB
}

type VersionST struct {
	Version string `json:"version" validate:"required"`
	Build   string `json:"build" validate:"required"`
} // @name Version

type LimitAndOffsetSearchST struct {
	Limit  int `query:"limit" validate:"required"`
	Offset int `query:"offset" validate:"required"`
} // @name LimitAndOffsetSearch
