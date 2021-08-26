package clinic

import (
	"context"

	"github.com/matijapetrovic/clinichub/clinic-service/internal/entity"
	"github.com/matijapetrovic/clinichub/clinic-service/pkg/dbcontext"
	"github.com/matijapetrovic/clinichub/clinic-service/pkg/log"
)

type Repository interface {
	Save(ctx context.Context, clinic entity.Clinic) (entity.Clinic, error)
	GetAll(ctx context.Context) ([]entity.Clinic, error)
	GetById(ctx context.Context, id string) (entity.Clinic, error)
}

type repository struct {
	db     *dbcontext.DB
	logger log.Logger
}

func NewRepository(db *dbcontext.DB, logger log.Logger) Repository {
	return repository{db, logger}
}

func (r repository) GetById(ctx context.Context, id string) (entity.Clinic, error) {
	return entity.Clinic{}, nil
}

func (r repository) GetAll(ctx context.Context) ([]entity.Clinic, error) {
	return make([]entity.Clinic, 0), nil
}

func (r repository) Save(ctx context.Context, clinic entity.Clinic) (entity.Clinic, error) {
	return clinic, nil
}
