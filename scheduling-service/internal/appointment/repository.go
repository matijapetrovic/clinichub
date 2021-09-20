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
	GetByPatientIdAndDate(ctx context.Context, patientId string, startDate time.Time, endDate time.Time) ([]entity.Appointment, error)
	GetClinicProfit(ctx context.Context, clinicId string, startDate time.Time, endDate time.Time) (int, error)
}

type repository struct {
	db     *dbcontext.DB
	logger log.Logger
}

func NewRepository(db *dbcontext.DB, logger log.Logger) Repository {
	return repository{db, logger}
}

func (r repository) GetClinicProfit(ctx context.Context, clinicId string, startDate time.Time, endDate time.Time) (int, error) {
	var profit int
	dbExp := dbx.NewExp("clinic_id={:clinicId}", dbx.Params{"clinicId": clinicId})
	dbExp = dbx.And(dbExp, dbx.NewExp("time>={:startDate}", dbx.Params{"startDate": startDate}))
	dbExp = dbx.And(dbExp, dbx.NewExp("time<={:endDate}", dbx.Params{"endDate": endDate}))

	err := r.db.With(ctx).Select("SUM(price) AS profit").From("appointment").Where(dbExp).Row(&profit)
	return profit, err
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

func (r repository) GetByPatientIdAndDate(ctx context.Context, patientId string, startDate time.Time, endDate time.Time) ([]entity.Appointment, error) {
	var appointments []entity.Appointment

	dbExp := dbx.NewExp("patient_id={:patientId}", dbx.Params{"patientId": patientId})
	if !startDate.IsZero() {
		dbExp = dbx.And(dbExp, dbx.NewExp("time>={:startDate}", dbx.Params{"startDate": startDate}))
	}
	if !endDate.IsZero() {
		dbExp = dbx.And(dbExp, dbx.NewExp("time<={:endDate}", dbx.Params{"endDate": endDate}))
	}

	err := r.db.With(ctx).
		Select().
		Where(dbExp).
		All(&appointments)

	return appointments, err
}

func (r repository) GetDoctorAppointments(ctx context.Context, doctorId string, dateStart time.Time, dateEnd time.Time) ([]entity.Appointment, error) {
	var appointments []entity.Appointment
	err := r.db.With(ctx).
		Select().
		Where(dbx.And(dbx.HashExp{"doctor_id": doctorId}, dbx.Between("time", dateStart, dateEnd))).
		All(&appointments)
	return appointments, err
}
