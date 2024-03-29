package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
)

type PlaceRowST struct {
	Id                 int32     `db:"id"`
	CreatorId          int32     `db:"creator_id"`
	Name               string    `db:"name"`
	StreetAddressLine1 string    `db:"street_address_line1"`
	StreetAddressLine2 string    `db:"street_address_line2"`
	Locality           string    `db:"locality"`
	Region             string    `db:"region"`
	PostalCode         string    `db:"postal_code"`
	Country            string    `db:"country"`
	Latitude           float64   `db:"latitude"`
	Longitude          float64   `db:"longitude"`
	UpdatedAt          time.Time `db:"updated_at"`
	CreatedAt          time.Time `db:"created_at"`
}

type PlaceFeatureRowST struct {
	Id              int32           `db:"id"`
	PlaceId         int32           `db:"place_id"`
	Type            string          `db:"type"`
	LocationType    string          `db:"location_type"`
	LatitudeOffset  sql.NullFloat64 `db:"latitude_offset"`
	LongitudeOffset sql.NullFloat64 `db:"longitude_offset"`
	UpdatedAt       time.Time       `db:"updated_at"`
	CreatedAt       time.Time       `db:"created_at"`
}

func GetPlaces(latitude, longitude, maxDistance float64, types []string, locations []string, limit, offset int) ([]PlaceRowST, error) {
	typesWhere := ""
	if len(types) != 0 {
		typesWhere = fmt.Sprintf("AND pf.type = ANY(%s)", strings.Join(types, ","))
	}
	locationsWhere := ""
	if len(locations) != 0 {
		locationsWhere = fmt.Sprintf("AND pf.location_type = ANY(%s)", strings.Join(locations, ","))
	}
	return All[PlaceRowST](fmt.Sprintf(`SELECT p.* 
		FROM places p 
		LEFT JOIN place_features pf ON pf.place_id = p.id
		WHERE
			(earth_box(ll_to_earth($1, $2), $3) @> ll_to_earth(p.latitude, p.longitude)
			AND earth_distance(ll_to_earth($1, $2), ll_to_earth(p.latitude, p.longitude)) < $3)
			%s
			%s
			LIMIT $4 OFFSET $5;`, typesWhere, locationsWhere),
		latitude,
		longitude,
		maxDistance,
		limit, offset)
}

func GetPlaceFeatures(latitude, longitude, maxDistance float64, types []string, locations []string, limit, offset int) ([]PlaceFeatureRowST, error) {
	typesWhere := ""
	if len(types) != 0 {
		typesWhere = fmt.Sprintf("AND pf.type = ANY(%s)", strings.Join(types, ","))
	}
	locationsWhere := ""
	if len(locations) != 0 {
		locationsWhere = fmt.Sprintf("AND pf.location_type = ANY(%s)", strings.Join(locations, ","))
	}
	return All[PlaceFeatureRowST](fmt.Sprintf(`SELECT pf.* 
		FROM place_features pf 
		WHERE pf.place_id IN (SELECT p.id
			FROM places p
			LEFT JOIN place_features pf ON pf.place_id = p.id
			WHERE
			(earth_box(ll_to_earth($1, $2), $3) @> ll_to_earth(p.latitude, p.longitude)
			AND earth_distance(ll_to_earth($1, $2), ll_to_earth(p.latitude, p.longitude)) < $3)
			%s
			%s
			LIMIT $4 OFFSET $5
		);`, typesWhere, locationsWhere),
		latitude,
		longitude,
		maxDistance,
		limit,
		offset)
}

func GetPlaceById(id int32) (*PlaceRowST, error) {
	return GetOptional[PlaceRowST](`SELECT p.*
		FROM plcaes p
		WHERE p.id = $1
		LIMIT 1;`,
		id)
}

func GetPlaceFeaturesByPlaceId(placeId int32) ([]PlaceFeatureRowST, error) {
	return All[PlaceFeatureRowST](`SELECT p.*
		FROM place_features p
		WHERE p.place_id = $1;`,
		placeId)
}
