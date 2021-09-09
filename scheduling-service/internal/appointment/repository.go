package appointment

import (
	"context"
	"time"

	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/matijapetrovic/clinichub/scheduling-service/internal/entity"
	"github.com/matijapetrovic/clinichub/scheduling-service/pkg/dbcontext"
	"github.com/matijapetrovic/clinichub/scheduling-service/pkg/log"
)

type Repository interface {
	Create(ctx context.Context, appointment entity.Appointment) error
	GetById(ctx context.Context, id string) (entity.Appointment, error)
	GetDoctorAppointments(ctx context.Context, doctorId string, dateStart time.Time, dateEnd time.Time) ([]entity.Appointment, error)
	GetByDoctorIdAndTime(ctx context.Context, doctorId string, time time.Time) (entity.Appointment, error)
}

type repository struct {
	db     *dbcontext.DB
	logger log.Logger
}

func NewRepository(db *dbcontext.DB, logger log.Logger) Repository {
	return repository{db, logger}
}

func (r repository) Create(ctx context.Context, appointment entity.Appointment) error {
	return r.db.With(ctx).Model(&appointment).Insert()
}

func (r repository) GetById(ctx context.Context, id string) (entity.Appointment, error) {
	var appointment entity.Appointment
	err := r.db.With(ctx).Select().Model(id, &appointment)
	return appointment, err
}

func (r repository) GetByDoctorIdAndTime(ctx context.Context, doctorId string, time time.Time) (entity.Appointment, error) {
	var appointment entity.Appointment
	err := r.db.With(ctx).
		Select().
		Where(dbx.HashExp{"doctor_id": doctorId, "time": time}).
		One(&appointment)

	return appointment, err
}

func (r repository) GetDoctorAppointments(ctx context.Context, doctorId string, dateStart time.Time, dateEnd time.Time) ([]entity.Appointment, error) {
	var appointments []entity.Appointment
	err := r.db.With(ctx).
		Select().
		Where(dbx.And(dbx.HashExp{"doctor_id": doctorId}, dbx.Between("time", dateStart, dateEnd))).
		All(&appointments)
	return appointments, err
}
