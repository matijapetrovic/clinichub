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
	GetPaged(ctx context.Context, offset int, limit int) ([]entity.Clinic, error)
	Count(ctx context.Context) (int, error)
	GetById(ctx context.Context, id string) (entity.Clinic, error)
	GetIdsByHasPrice(ctx context.Context, appointmentTypeId string) ([]string, error)
	GetByIdList(ctx context.Context, clinicIds []string) ([]entity.Clinic, error)

	GetAppointmentTypePrices(ctx context.Context, clinicId string) ([]entity.AppointmentTypePrice, error)
	GetAppointmentTypePrice(ctx context.Context, clinicId string, appointmentTypeId string) (entity.AppointmentTypePrice, error)
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

func (r repository) GetByIdList(ctx context.Context, clinicIds []string) ([]entity.Clinic, error) {
	var clinics []entity.Clinic
	b := make([]interface{}, len(clinicIds))
	for i := range clinicIds {
		b[i] = clinicIds[i]
	}
	err := r.db.With(ctx).Select().Where(dbx.In("id", b...)).All(&clinics)
	return clinics, err
}

func (r repository) GetIdsByHasPrice(ctx context.Context, appointmentTypeId string) ([]string, error) {
	var prices []entity.AppointmentTypePrice
	err := r.db.With(ctx).
		Select().
		Distinct(true).
		Where(dbx.HashExp{"appointment_type_id": appointmentTypeId}).
		All(&prices)

	var clinicIds []string
	for _, price := range prices {
		clinicIds = append(clinicIds, price.ClinicId)
	}
	return clinicIds, err
}

func (r repository) GetById(ctx context.Context, id string) (entity.Clinic, error) {
	var clinic entity.Clinic
	err := r.db.With(ctx).Select().Model(id, &clinic)
	return clinic, err
}

func (r repository) Count(ctx context.Context) (int, error) {
	var count int
	err := r.db.With(ctx).Select("COUNT(*)").From("clinic").Row(&count)
	return count, err
}

func (r repository) GetPaged(ctx context.Context, offset int, limit int) ([]entity.Clinic, error) {
	var clinics []entity.Clinic
	err := r.db.With(ctx).
		Select().
		OrderBy("name").
		Offset(int64(offset)).
		Limit(int64(limit)).
		All(&clinics)
	return clinics, err
}

func (r repository) Create(ctx context.Context, clinic entity.Clinic) error {
	return r.db.With(ctx).Model(&clinic).Exclude("AppointmentPrices").Insert()
}

func (r repository) Update(ctx context.Context, clinic entity.Clinic) error {
	return r.db.With(ctx).Model(&clinic).Exclude("AppointmentPrices").Update()
}

func (r repository) GetAppointmentTypePrice(ctx context.Context, clinicId string, appointmentTypeId string) (entity.AppointmentTypePrice, error) {
	var appointmentTypePrice entity.AppointmentTypePrice
	err := r.db.With(ctx).
		Select().
		Where(dbx.HashExp{"clinic_id": clinicId, "appointment_type_id": appointmentTypeId}).
		One(&appointmentTypePrice)
	return appointmentTypePrice, err
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
