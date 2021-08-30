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
	GetByClinicId(ctx context.Context, clinicId string) ([]entity.Doctor, error)
	Create(ctx context.Context, req CreateDoctorRequest) (entity.Doctor, error)
	Update(ctx context.Context, doctorId string, req UpdateDoctorRequest) (entity.Doctor, error)
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
		// dodaj validaciju za work
	)
}

type UpdateDoctorRequest struct {
	FirstName string      `json:"firstName"`
	LastName  string      `json:"lastName"`
	WorkStart entity.Time `json:"workStart"`
	WorkEnd   entity.Time `json:"workEnd"`
}

func (m UpdateDoctorRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.FirstName, validation.Required, validation.Length(1, 50)),
		validation.Field(&m.LastName, validation.Required, validation.Length(1, 50)),
		// dodaj validaciju za work
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

	_, err = s.appointmentTypeRepo.GetById(ctx, req.SpecializationId)
	if err != nil {
		return entity.Doctor{}, err
	}

	id := entity.GenerateID()
	err = s.repo.Create(ctx, entity.Doctor{
		Id:               id,
		FirstName:        req.FirstName,
		LastName:         req.LastName,
		WorkStart:        req.WorkStart.ToString(),
		WorkEnd:          req.WorkEnd.ToString(),
		SpecializationId: req.SpecializationId,
		ClinicId:         clinic.Id,
	})
	if err != nil {
		return entity.Doctor{}, err
	}

	return s.repo.GetById(ctx, id)
}

func (s service) Update(ctx context.Context, doctorId string, req UpdateDoctorRequest) (entity.Doctor, error) {
	if err := req.Validate(); err != nil {
		return entity.Doctor{}, err
	}

	doctor, err := s.repo.GetById(ctx, doctorId)
	if err != nil {
		return entity.Doctor{}, err
	}

	doctor.FirstName = req.FirstName
	doctor.LastName = req.LastName
	doctor.WorkStart = req.WorkStart.ToString()
	doctor.WorkEnd = req.WorkEnd.ToString()

	err = s.repo.Update(ctx, doctor)
	if err != nil {
		return entity.Doctor{}, err
	}

	specialization, err := s.appointmentTypeRepo.GetById(ctx, doctor.SpecializationId)
	if err != nil {
		return entity.Doctor{}, err
	}
	doctor.AppointmentType = specialization
	return doctor, nil
}

func (s service) GetByClinicId(ctx context.Context, clinicId string) ([]entity.Doctor, error) {
	doctors, err := s.repo.GetByClinicId(ctx, clinicId)
	if err != nil {
		return nil, err
	}

	for i, doctor := range doctors {
		specialization, err := s.appointmentTypeRepo.GetById(ctx, doctor.SpecializationId)
		if err != nil {
			return nil, err
		}
		doctor.AppointmentType = specialization
		doctors[i] = doctor
	}

	return doctors, nil
}

func (s service) GetAll(ctx context.Context) ([]entity.Doctor, error) {
	doctors, err := s.repo.GetAll(ctx)

	if err != nil {
		return nil, err
	}

	return doctors, nil
}
