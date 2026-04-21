package application

import (
	"context"
	"fmt"

	"github.com/akgate/photo/internal/domain/Contract"
	"github.com/akgate/photo/internal/domain/Entity"
	"github.com/akgate/platform/pkg/db"
	"github.com/google/uuid"
)

type CreatePhotoCommand struct {
	LocationID uuid.UUID
	Lat        float64
	Lng        float64
}

type CommandProcessor struct {
	storage   Contract.PhotoWriteStorage
	txManager db.TxManager
}

func NewCommandProcessor(storage Contract.PhotoWriteStorage, txManager db.TxManager) *CommandProcessor {
	return &CommandProcessor{
		storage:   storage,
		txManager: txManager,
	}
}

func (c *CommandProcessor) CreateMany(ctx context.Context, cmds []CreatePhotoCommand) ([]*Entity.Photo, error) {
	created := make([]*Entity.Photo, 0, len(cmds))

	for _, cmd := range cmds {
		coordinates, err := Entity.NewCoordinates(cmd.Lat, cmd.Lng)
		if err != nil {
			return nil, fmt.Errorf("build coordinates: %w", err)
		}

		photo := Entity.NewPhoto(
			Entity.LocationId(cmd.LocationID),
			coordinates,
			uuid.New(),
		)

		created = append(created, photo)
	}

	err := c.txManager.ReadCommited(ctx, func(ctx context.Context) error {
		return c.storage.CreateMany(ctx, created)
	})
	if err != nil {
		return nil, fmt.Errorf("bulk create photos: %w", err)
	}

	return created, nil
}

func (c *CommandProcessor) Delete(ctx context.Context, id uuid.UUID) error {
	return c.txManager.ReadCommited(ctx, func(ctx context.Context) error {
		return c.storage.Delete(ctx, id)
	})
}

func (c *CommandProcessor) Update(ctx context.Context, photo *Entity.Photo) error {
	return c.txManager.ReadCommited(ctx, func(ctx context.Context) error {
		return c.storage.Update(ctx, photo)
	})
}
