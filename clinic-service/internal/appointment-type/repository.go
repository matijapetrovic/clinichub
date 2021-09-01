package appointment_type

import (
	"context"

	"github.com/matijapetrovic/clinichub/clinic-service/internal/entity"
	"github.com/matijapetrovic/clinichub/clinic-service/pkg/dbcontext"
	"github.com/matijapetrovic/clinichub/clinic-service/pkg/log"
)

type Repository interface {
	GetById(ctx context.Context, id string) (entity.AppointmentType, error)
	GetAll(ctx context.Context) ([]entity.AppointmentType, error)
	Create(ctx context.Context, appointmentType entity.AppointmentType) error
	Update(ctx context.Context, appointmentType entity.AppointmentType) error
}

type repository struct {
	db     *dbcontext.DB
	logger log.Logger
}

func NewRepository(db *dbcontext.DB, logger log.Logger) Repository {
	return repository{db, logger}
}

func (r repository) GetAll(ctx context.Context) ([]entity.AppointmentType, error) {
	var appointmentTypes []entity.AppointmentType
	err := r.db.With(ctx).Select().All(&appointmentTypes)
	return appointmentTypes, err
}

func (r repository) GetById(ctx context.Context, id string) (entity.AppointmentType, error) {
	var appointmentType entity.AppointmentType
	err := r.db.With(ctx).Select().Model(id, &appointmentType)
	return appointmentType, err
}

func (r repository) Create(ctx context.Context, appointmentType entity.AppointmentType) error {
	return r.db.With(ctx).Model(&appointmentType).Insert()
}

func (r repository) Update(ctx context.Context, appointmentType entity.AppointmentType) error {
	return r.db.With(ctx).Model(&appointmentType).Update()
}
