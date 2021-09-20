package doctor_rating

import (
	"context"

	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/matijapetrovic/clinichub/rating-service/pkg/log"

	"github.com/matijapetrovic/clinichub/rating-service/internal/entity"
	"github.com/matijapetrovic/clinichub/rating-service/pkg/dbcontext"
)

type Repository interface {
	GetById(ctx context.Context, id string) (entity.DoctorRating, error)
	GetRating(ctx context.Context, patientId string, doctorId string) (entity.DoctorRating, error)
	GetDoctorRating(ctx context.Context, doctorId string) (entity.AverageRating, error)
	RateDoctor(ctx context.Context, rating entity.DoctorRating) error
}

type repository struct {
	db     *dbcontext.DB
	logger log.Logger
}

func NewRepository(db *dbcontext.DB, logger log.Logger) Repository {
	return repository{db, logger}
}

func (r repository) GetRating(ctx context.Context, patientId string, doctorID string) (entity.DoctorRating, error) {
	var doctorRating entity.DoctorRating
	err := r.db.With(ctx).Select().Where(dbx.HashExp{"patient_id": patientId, "doctor_id": doctorID}).One(&doctorRating)
	return doctorRating, err
}

func (r repository) GetDoctorRating(ctx context.Context, doctorId string) (entity.AverageRating, error) {
	var count int
	var rating float32
	err := r.db.With(ctx).Select("COUNT(rating) AS count", "AVG(rating) AS rating").Where(dbx.HashExp{"doctor_id": doctorId}).Row(&count, &rating)

	return entity.AverageRating{Count: count, Rating: rating}, err
}

func (r repository) RateDoctor(ctx context.Context, rating entity.DoctorRating) error {
	return r.db.With(ctx).Model(&rating).Insert()
}

func (r repository) GetById(ctx context.Context, id string) (entity.DoctorRating, error) {
	var rating entity.DoctorRating
	err := r.db.With(ctx).Select().Model(id, &rating)
	return rating, err
}
