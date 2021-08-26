package appointment_type

import (
	"context"

	"github.com/matijapetrovic/clinichub/clinic-service/internal/entity"
	"github.com/matijapetrovic/clinichub/clinic-service/pkg/dbcontext"
	"github.com/matijapetrovic/clinichub/clinic-service/pkg/log"
)

type Repository interface {
	GetById(ctx context.Context, id string) (entity.AppointmentType, error)
}

type repository struct {
	db     *dbcontext.DB
	logger log.Logger
}

func NewRepository(db *dbcontext.DB, logger log.Logger) Repository {
	return repository{db, logger}
}

func (r repository) GetById(ctx context.Context, id string) (entity.AppointmentType, error) {
	return entity.AppointmentType{}, nil
}
