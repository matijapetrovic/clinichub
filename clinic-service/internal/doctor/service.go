package doctor

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	appointment_type "github.com/matijapetrovic/clinichub/clinic-service/internal/appointment-type"
	"github.com/matijapetrovic/clinichub/clinic-service/internal/clinic"
	"github.com/matijapetrovic/clinichub/clinic-service/internal/entity"
	"github.com/matijapetrovic/clinichub/clinic-service/pkg/log"
)

type Service interface {
	GetAll(ctx context.Context) ([]entity.Doctor, error)
	Create(ctx context.Context, req CreateDoctorRequest) (entity.Doctor, error)
}

type CreateDoctorRequest struct {
	FirstName        string      `json:"firstName"`
	LastName         string      `json:"lastName"`
	WorkStart        entity.Time `json:"workStart"`
	WorkEnd          entity.Time `json:"workEnd"`
	SpecializationId string      `json:"specializationId"`
	ClinicId         string      `json:"clinicId"`
}

func (m CreateDoctorRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.FirstName, validation.Required, validation.Length(1, 50)),
		validation.Field(&m.LastName, validation.Required, validation.Length(1, 50)),
		validation.Field(&m.SpecializationId, validation.Required, validation.Length(36, 36)),
		validation.Field(&m.ClinicId, validation.Required, validation.Length(36, 36)),
	)
}

type service struct {
	repo                Repository
	clinicRepo          clinic.Repository
	appointmentTypeRepo appointment_type.Repository
	logger              log.Logger
}

func NewService(repo Repository, clinicRepo clinic.Repository, appointmentTypeRepo appointment_type.Repository, logger log.Logger) Service {
	return service{repo, clinicRepo, appointmentTypeRepo, logger}
}

func (s service) Create(ctx context.Context, req CreateDoctorRequest) (entity.Doctor, error) {
	if err := req.Validate(); err != nil {
		return entity.Doctor{}, err
	}

	clinic, err := s.clinicRepo.GetById(ctx, req.ClinicId)
	if err != nil {
		return entity.Doctor{}, err
	}

	appointmentType, err := s.appointmentTypeRepo.GetById(ctx, req.SpecializationId)
	if err != nil {
		return entity.Doctor{}, err
	}

	id := entity.GenerateID()
	doctor, err := s.repo.Save(ctx, entity.Doctor{
		Id:             id,
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		WorkStart:      req.WorkStart,
		WorkEnd:        req.WorkEnd,
		Specialization: appointmentType,
		Clinic:         clinic,
	})
	if err != nil {
		return entity.Doctor{}, err
	}

	clinic.Doctors = append(clinic.Doctors, doctor)
	clinic, err = s.clinicRepo.Save(ctx, clinic)
	if err != nil {
		return entity.Doctor{}, err
	}

	return doctor, nil
}

func (s service) GetAll(ctx context.Context) ([]entity.Doctor, error) {
	clinics, err := s.repo.GetAll(ctx)

	if err != nil {
		return nil, err
	}

	return clinics, nil
}
