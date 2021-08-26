package doctor

import (
	"context"

	"github.com/matijapetrovic/clinichub/clinic-service/internal/entity"
	"github.com/matijapetrovic/clinichub/clinic-service/pkg/dbcontext"
	"github.com/matijapetrovic/clinichub/clinic-service/pkg/log"
)

type Repository interface {
	Save(ctx context.Context, clinic entity.Doctor) (entity.Doctor, error)
	GetAll(ctx context.Context) ([]entity.Doctor, error)
	GetByClinic(ctx context.Context, clinicId string) ([]entity.Doctor, error)
}

type repository struct {
	db     *dbcontext.DB
	logger log.Logger
}

func NewRepository(db *dbcontext.DB, logger log.Logger) Repository {
	return repository{db, logger}
}

func (r repository) GetByClinic(ctx context.Context, clinicId string) ([]entity.Doctor, error) {
	return make([]entity.Doctor, 0), nil
}

func (r repository) GetAll(ctx context.Context) ([]entity.Doctor, error) {
	return make([]entity.Doctor, 0), nil
}

func (r repository) Save(ctx context.Context, doctor entity.Doctor) (entity.Doctor, error) {
	return doctor, nil
}
