package Contract

import (
	"context"

	"github.com/akgate/photo/internal/domain/Entity"
	"github.com/google/uuid"
)

type PhotoReadStorage interface {
	GetById(ctx context.Context, id uuid.UUID) (*Entity.Photo, error)
}
