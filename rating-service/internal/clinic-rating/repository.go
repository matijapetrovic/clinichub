package doctor_rating

import (
	"context"

	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/matijapetrovic/clinichub/rating-service/pkg/log"

	"github.com/matijapetrovic/clinichub/rating-service/internal/entity"
	"github.com/matijapetrovic/clinichub/rating-service/pkg/dbcontext"
)

type Repository interface {
	GetById(ctx context.Context, id string) (entity.ClinicRating, error)
	GetRating(ctx context.Context, patientId string, clinicId string) (entity.ClinicRating, error)
	GetClinicRating(ctx context.Context, clinicId string) (entity.AverageRating, error)
	RateClinic(ctx context.Context, rating entity.ClinicRating) error
}

type repository struct {
	db     *dbcontext.DB
	logger log.Logger
}

func NewRepository(db *dbcontext.DB, logger log.Logger) Repository {
	return repository{db, logger}
}

func (r repository) GetRating(ctx context.Context, patientId string, clinicId string) (entity.ClinicRating, error) {
	var clinicRating entity.ClinicRating
	err := r.db.With(ctx).Select().Where(dbx.HashExp{"patient_id": patientId, "clinic_id": clinicId}).One(&clinicRating)
	return clinicRating, err
}

func (r repository) GetClinicRating(ctx context.Context, clinicId string) (entity.AverageRating, error) {
	var count int
	var rating float32
	err := r.db.With(ctx).Select("COUNT(rating) AS count", "AVG(rating) AS rating").Where(dbx.HashExp{"clinic_id": clinicId}).Row(&count, &rating)

	return entity.AverageRating{Count: count, Rating: rating}, err
}

func (r repository) RateClinic(ctx context.Context, rating entity.ClinicRating) error {
	return r.db.With(ctx).Model(&rating).Insert()
}

func (r repository) GetById(ctx context.Context, id string) (entity.ClinicRating, error) {
	var rating entity.ClinicRating
	err := r.db.With(ctx).Select().Model(id, &rating)
	return rating, err
}
