package clinic

import (
	"context"

	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/matijapetrovic/clinichub/clinic-service/internal/entity"
	"github.com/matijapetrovic/clinichub/clinic-service/pkg/dbcontext"
	"github.com/matijapetrovic/clinichub/clinic-service/pkg/log"
)

type Repository interface {
	Create(ctx context.Context, clinic entity.Clinic) error
	Update(ctx context.Context, clinic entity.Clinic) error
	GetAll(ctx context.Context) ([]entity.Clinic, error)
	GetById(ctx context.Context, id string) (entity.Clinic, error)
	GetAppointmentTypePrices(ctx context.Context, clinicId string) ([]entity.AppointmentTypePrice, error)
	AddAppointmentTypePrice(ctx context.Context, appointmentTypePrice entity.AppointmentTypePrice) error
	UpdateAppointmentTypePrice(ctx context.Context, appointmentTypePrice entity.AppointmentTypePrice) error
}

type repository struct {
	db     *dbcontext.DB
	logger log.Logger
}

func NewRepository(db *dbcontext.DB, logger log.Logger) Repository {
	return repository{db, logger}
}

func (r repository) GetById(ctx context.Context, id string) (entity.Clinic, error) {
	var clinic entity.Clinic
	err := r.db.With(ctx).Select().Model(id, &clinic)
	return clinic, err
}

func (r repository) GetAll(ctx context.Context) ([]entity.Clinic, error) {
	var clinics []entity.Clinic
	err := r.db.With(ctx).Select().All(&clinics)
	return clinics, err
}

func (r repository) Create(ctx context.Context, clinic entity.Clinic) error {
	return r.db.With(ctx).Model(&clinic).Exclude("AppointmentPrices").Insert()
}

func (r repository) Update(ctx context.Context, clinic entity.Clinic) error {
	return r.db.With(ctx).Model(&clinic).Exclude("AppointmentPrices").Update()
}

func (r repository) GetAppointmentTypePrices(ctx context.Context, clinicId string) ([]entity.AppointmentTypePrice, error) {
	var appointmentTypePrices []entity.AppointmentTypePrice
	err := r.db.With(ctx).
		Select().
		Where(dbx.HashExp{"clinic_id": clinicId}).
		All(&appointmentTypePrices)
	return appointmentTypePrices, err
}

func (r repository) AddAppointmentTypePrice(ctx context.Context, appointmentTypePrice entity.AppointmentTypePrice) error {
	return r.db.With(ctx).Model(&appointmentTypePrice).Insert()
}

func (r repository) UpdateAppointmentTypePrice(ctx context.Context, appointmentTypePrice entity.AppointmentTypePrice) error {
	return r.db.With(ctx).Model(&appointmentTypePrice).Update()
}
