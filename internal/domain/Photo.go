package domain

import "github.com/google/uuid"

type LocationId uuid.UUID

type Photo struct {
	id          uuid.UUID
	locationId  LocationId
	coordinates Coordinates
}

func NewPhotoEntity(locationId LocationId, coordinates Coordinates, id *uuid.UUID) Photo {
	return Photo{
		id:          *id,
		locationId:  locationId,
		coordinates: coordinates,
	}
}

func (e *Photo) ID() uuid.UUID {
	return e.id
}

func (e *Photo) Coordinates() Coordinates {
	return e.coordinates
}

func (e *Photo) LocationId() LocationId {
	return e.locationId
}
