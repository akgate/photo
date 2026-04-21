package Contract

import (
	"context"

	"github.com/akgate/photo/internal/domain/Entity"
	"github.com/google/uuid"
)

type PhotoWriteStorage interface {
	CreateMany(ctx context.Context, photos []*Entity.Photo) error
	Update(ctx context.Context, photo *Entity.Photo) error
	Delete(ctx context.Context, id uuid.UUID) error
}
