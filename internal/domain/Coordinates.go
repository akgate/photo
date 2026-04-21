package domain

import "fmt"

type Coordinates struct {
	lat float64
	lng float64
}

func NewCoordinates(lat float64, lng float64) (Coordinates, error) {
	if lat < -90 || lat > 90 {
		return Coordinates{}, fmt.Errorf("invalid latitude")
	}
	if lng < -180 || lng > 180 {
		return Coordinates{}, fmt.Errorf("invalid longitude")
	}
	return Coordinates{lat: lat, lng: lng}, nil
}

func (c *Coordinates) Lat() float64 {
	return c.lat
}

func (c *Coordinates) Lng() float64 {
	return c.lng
}
