package model

import (
	"time"

	"github.com/aicacia/pickup/app/repository"
)

type PlaceSearchST struct {
	LimitAndOffsetSearchST
	Types         []string `query:"types"`
	LocationTypes []string `query:"locations"`
	Latitude      float64  `query:"latitude" validate:"required"`
	Longitude     float64  `query:"longitude" validate:"required"`
	MaxDistance   float64  `query:"max_distance" validate:"required"`
} // @name PlaceSearch

type PlaceFeatureST struct {
	Id           int32     `json:"id" validate:"required"`
	PlaceId      int32     `json:"place_id" validate:"required"`
	Type         string    `json:"type" validate:"required"`
	LocationType string    `json:"location_type" validate:"required"`
	Latitude     float64   `json:"latitude" validate:"required"`
	Longitude    float64   `json:"longitude" validate:"required"`
	UpdatedAt    time.Time `json:"updated_at" validate:"required" format:"date-time"`
	CreatedAt    time.Time `json:"created_at" validate:"required" format:"date-time"`
} // @name PlaceFeature

type PlaceST struct {
	Id                 int32            `json:"id" validate:"required"`
	CreatorId          int32            `json:"creator_id" validate:"required"`
	Name               string           `json:"name" validate:"required"`
	StreetAddressLine1 string           `json:"street_address_line1" validate:"required"`
	StreetAddressLine2 string           `json:"street_address_line2" validate:"required"`
	Locality           string           `json:"locality" validate:"required"`
	Region             string           `json:"region" validate:"required"`
	PostalCode         string           `json:"postal_code" validate:"required"`
	Country            string           `json:"country" validate:"required"`
	Latitude           float64          `json:"latitude" validate:"required"`
	Longitude          float64          `json:"longitude" validate:"required"`
	Features           []PlaceFeatureST `json:"features" validate:"required"`
	UpdatedAt          time.Time        `json:"updated_at" validate:"required" format:"date-time"`
	CreatedAt          time.Time        `json:"created_at" validate:"required" format:"date-time"`
} // @name Place

func PlaceFeatureFromPlaceFeatureRow(placeFeatureRow repository.PlaceFeatureRowST, latitude, longitude float64) PlaceFeatureST {
	return PlaceFeatureST{
		Id:           placeFeatureRow.Id,
		PlaceId:      placeFeatureRow.PlaceId,
		Type:         placeFeatureRow.Type,
		LocationType: placeFeatureRow.LocationType,
		Latitude:     latitude + placeFeatureRow.LatitudeOffset.Float64,
		Longitude:    longitude + placeFeatureRow.LongitudeOffset.Float64,
		UpdatedAt:    placeFeatureRow.UpdatedAt,
		CreatedAt:    placeFeatureRow.CreatedAt,
	}
}

func PlaceFromPlaceRow(placeRow repository.PlaceRowST, featureRows []repository.PlaceFeatureRowST) PlaceST {
	features := make([]PlaceFeatureST, 0, len(featureRows))
	for _, featureRow := range featureRows {
		features = append(features, PlaceFeatureFromPlaceFeatureRow(featureRow, placeRow.Latitude, placeRow.Longitude))
	}
	return PlaceST{
		Id:                 placeRow.Id,
		CreatorId:          placeRow.CreatorId,
		Name:               placeRow.Name,
		StreetAddressLine1: placeRow.StreetAddressLine1,
		StreetAddressLine2: placeRow.StreetAddressLine2,
		Locality:           placeRow.Locality,
		Region:             placeRow.Region,
		PostalCode:         placeRow.PostalCode,
		Country:            placeRow.Country,
		Latitude:           placeRow.Latitude,
		Longitude:          placeRow.Longitude,
		Features:           features,
		UpdatedAt:          placeRow.UpdatedAt,
		CreatedAt:          placeRow.CreatedAt,
	}
}
