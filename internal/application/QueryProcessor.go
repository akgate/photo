package application

import (
	"context"

	"github.com/akgate/photo/internal/domain/Contract"
	"github.com/akgate/photo/internal/domain/Entity"
	"github.com/google/uuid"
)

type QueryProcessor struct {
	storage Contract.PhotoReadStorage
}

func NewQueryProcessor(storage Contract.PhotoReadStorage) *QueryProcessor {
	return &QueryProcessor{storage: storage}
}

func (c *QueryProcessor) GetById(ctx context.Context, id uuid.UUID) (*Entity.Photo, error) {
	photo, err := c.storage.GetById(ctx, id)
	return photo, err
}
