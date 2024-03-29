package controller

import (
	"log"

	"github.com/aicacia/pickup/app/model"
	"github.com/aicacia/pickup/app/repository"
	"github.com/gofiber/fiber/v2"
)

// GetPlaces
//
//	@Summary		Get places
//	@Tags			place
//	@Accept			json
//	@Produce		json
//	@Param          search          query    model.PlaceSearchST true "place query params"
//	@Success		200	{array}	    model.PlaceST
//	@Failure		400	{object}	model.ErrorST
//	@Failure		401	{object}	model.ErrorST
//	@Failure		403	{object}	model.ErrorST
//	@Failure		500	{object}	model.ErrorST
//	@Router			/places [get]
//	@Security		Authorization
func GetPlaces(c *fiber.Ctx) error {
	var search model.PlaceSearchST
	if err := c.QueryParser(&search); err != nil {
		log.Printf("error parsing search: %v\n", err)
		return model.NewError(400).AddError("search", "invalid", "query").Send(c)
	}
	placeRows, err := repository.GetPlaces(search.Latitude, search.Longitude, search.MaxDistance, search.Types, search.LocationTypes, search.Limit, search.Offset)
	if err != nil {
		log.Printf("error getting places from db: %v\n", err)
		return model.NewError(500).AddError("database", "internal", "application").Send(c)
	}
	placeFeaturesRows, err := repository.GetPlaceFeatures(search.Latitude, search.Longitude, search.MaxDistance, search.Types, search.LocationTypes, search.Limit, search.Offset)
	if err != nil {
		log.Printf("error getting places features from db: %v\n", err)
		return model.NewError(500).AddError("database", "internal", "application").Send(c)
	}
	placeFeaturesRowsByPlaceId := make(map[int32][]repository.PlaceFeatureRowST, len(placeFeaturesRows))
	for _, placeFeaturesRow := range placeFeaturesRows {
		placeFeaturesRowsByPlaceId[placeFeaturesRow.PlaceId] = append(placeFeaturesRowsByPlaceId[placeFeaturesRow.PlaceId], placeFeaturesRow)
	}
	places := make([]model.PlaceST, 0, len(placeRows))
	for _, placeRow := range placeRows {
		places = append(places, model.PlaceFromPlaceRow(placeRow, placeFeaturesRowsByPlaceId[placeRow.Id]))
	}
	return c.JSON(places)
}
