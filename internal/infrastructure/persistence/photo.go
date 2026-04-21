package persistence

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/akgate/photo/internal/domain/Entity"
	"github.com/akgate/platform/pkg/db"
	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
)

type PhotoRepository struct {
	db db.DB
}

func NewPhotoRepository(db db.DB) *PhotoRepository {
	return &PhotoRepository{db: db}
}

func (r *PhotoRepository) GetById(ctx context.Context, id uuid.UUID) (photo *Entity.Photo, err error) {
	query := db.Query{Name: "getPhotoById", QueryRaw: "SELECT id, location_id, coordinates FROM photo WHERE id = $1"}

	row := r.db.QueryRowContext(ctx, query, id)

	if row == nil {
		return photo, errors.New("photo not found")
	}

	var (
		locationID     uuid.UUID
		rawCoordinates string
	)

	err = row.Scan(&id, &locationID, &rawCoordinates)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return photo, nil
		}

		return photo, err
	}

	coordinates, err := parseCoordinates(rawCoordinates)
	if err != nil {
		return photo, err
	}

	photo = Entity.NewPhoto(Entity.LocationId(locationID), coordinates, id)

	return photo, err
}

func (r *PhotoRepository) CreateMany(ctx context.Context, photos []*Entity.Photo) error {
	q := db.Query{
		Name:     "createPhoto",
		QueryRaw: "INSERT INTO photo (id, location_id, coordinates) VALUES ($1, $2, point($3, $4))",
	}

	for _, photo := range photos {
		coordinates := photo.Coordinates()

		_, err := r.db.ExecContext(
			ctx,
			q,
			photo.ID(),
			uuid.UUID(photo.LocationId()),
			coordinates.Lat(),
			coordinates.Lng(),
		)
		if err != nil {
			return fmt.Errorf("insert photo %s: %w", photo.ID(), err)
		}
	}

	return nil
}

func (r *PhotoRepository) Update(ctx context.Context, photo *Entity.Photo) error {
	q := db.Query{
		Name:     "updatePhoto",
		QueryRaw: "UPDATE photo SET location_id = $2, coordinates = point($3, $4) WHERE id = $1",
	}

	coordinates := photo.Coordinates()

	tag, err := r.db.ExecContext(
		ctx,
		q,
		photo.ID(),
		uuid.UUID(photo.LocationId()),
		coordinates.Lat(),
		coordinates.Lng(),
	)
	if err != nil {
		return fmt.Errorf("update photo %s: %w", photo.ID(), err)
	}

	if tag.RowsAffected() == 0 {
		return errors.New("photo not found")
	}

	return nil
}

func (r *PhotoRepository) Delete(ctx context.Context, id uuid.UUID) error {
	q := db.Query{
		Name:     "deletePhoto",
		QueryRaw: "DELETE FROM photo WHERE id = $1",
	}

	tag, err := r.db.ExecContext(ctx, q, id)
	if err != nil {
		return fmt.Errorf("delete photo %s: %w", id, err)
	}

	if tag.RowsAffected() == 0 {
		return errors.New("photo not found")
	}

	return nil
}

func parseCoordinates(raw string) (Entity.Coordinates, error) {
	raw = strings.TrimSpace(raw)
	raw = strings.TrimPrefix(raw, "(")
	raw = strings.TrimSuffix(raw, ")")

	parts := strings.Split(raw, ",")
	if len(parts) != 2 {
		return Entity.Coordinates{}, fmt.Errorf("invalid coordinates format: %q", raw)
	}

	lat, err := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
	if err != nil {
		return Entity.Coordinates{}, fmt.Errorf("parse latitude: %w", err)
	}

	lng, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
	if err != nil {
		return Entity.Coordinates{}, fmt.Errorf("parse longitude: %w", err)
	}

	coordinates, err := Entity.NewCoordinates(lat, lng)
	if err != nil {
		return Entity.Coordinates{}, fmt.Errorf("build coordinates: %w", err)
	}

	return coordinates, nil
}
