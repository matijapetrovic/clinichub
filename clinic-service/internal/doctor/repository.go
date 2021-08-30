package doctor

import (
	"context"

	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/matijapetrovic/clinichub/clinic-service/internal/entity"
	"github.com/matijapetrovic/clinichub/clinic-service/pkg/dbcontext"
	"github.com/matijapetrovic/clinichub/clinic-service/pkg/log"
)

type Repository interface {
	Create(ctx context.Context, doctor entity.Doctor) error
	Update(ctx context.Context, doctor entity.Doctor) error
	GetById(ctx context.Context, id string) (entity.Doctor, error)
	GetAll(ctx context.Context) ([]entity.Doctor, error)
	GetByClinicId(ctx context.Context, clinicId string) ([]entity.Doctor, error)
}

type repository struct {
	db     *dbcontext.DB
	logger log.Logger
}

func NewRepository(db *dbcontext.DB, logger log.Logger) Repository {
	return repository{db, logger}
}

func (r repository) GetByClinicId(ctx context.Context, clinicId string) ([]entity.Doctor, error) {
	var doctors []entity.Doctor
	err := r.db.With(ctx).
		Select().
		Where(dbx.HashExp{"clinic_id": clinicId}).
		All(&doctors)
	return doctors, err
}

func (r repository) GetById(ctx context.Context, id string) (entity.Doctor, error) {
	var doctor entity.Doctor
	err := r.db.With(ctx).
		Select().
		Model(id, &doctor)
	return doctor, err
}

func (r repository) GetAll(ctx context.Context) ([]entity.Doctor, error) {
	var doctors []entity.Doctor
	err := r.db.With(ctx).Select().All(&doctors)
	return doctors, err
}

func (r repository) Create(ctx context.Context, doctor entity.Doctor) error {
	return r.db.With(ctx).Model(&doctor).Insert()
}

func (r repository) Update(ctx context.Context, doctor entity.Doctor) error {
	return r.db.With(ctx).Model(&doctor).Update()
}
