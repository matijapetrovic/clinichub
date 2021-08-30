package clinic

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	appointment_type "github.com/matijapetrovic/clinichub/clinic-service/internal/appointment-type"
	"github.com/matijapetrovic/clinichub/clinic-service/internal/entity"
	"github.com/matijapetrovic/clinichub/clinic-service/pkg/log"
)

type Service interface {
	GetById(ctx context.Context, id string) (entity.Clinic, error)
	GetAll(ctx context.Context) ([]entity.Clinic, error)
	Create(ctx context.Context, req CreateClinicRequest) (entity.Clinic, error)
	Update(ctx context.Context, clinicId string, req UpdateClinicRequest) (entity.Clinic, error)
	AddAppointmentTypePrice(ctx context.Context, req AddAppointmentTypePriceRequest) error
	GetAppointmentTypePrices(ctx context.Context, clinicId string) ([]entity.AppointmentTypePrice, error)
}

type CreateClinicRequest struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Address     entity.Address `json:"address"`
}

func (m CreateClinicRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Name, validation.Required, validation.Length(1, 50)),
		validation.Field(&m.Description, validation.Required, validation.Length(1, 256)),
	)
}

type UpdateClinicRequest struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Address     entity.Address `json:"address"`
}

func (m UpdateClinicRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Name, validation.Required, validation.Length(1, 50)),
		validation.Field(&m.Description, validation.Required, validation.Length(1, 256)),
	)
}

type AddAppointmentTypePriceRequest struct {
	ClinicId          string `json:"clinicId"`
	AppointmentTypeId string `json:"appointmentTypeId"`
	Price             uint   `json:"price"`
}

func (m AddAppointmentTypePriceRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.ClinicId, validation.Required, validation.Length(36, 36)),
		validation.Field(&m.AppointmentTypeId, validation.Required, validation.Length(36, 36)),
		validation.Field(&m.Price, validation.Required, validation.Min(1)),
	)
}

type service struct {
	repo                Repository
	appointmentTypeRepo appointment_type.Repository
	logger              log.Logger
}

func NewService(repo Repository, appointmentTypeRepo appointment_type.Repository, logger log.Logger) Service {
	return service{repo, appointmentTypeRepo, logger}
}

func (s service) GetById(ctx context.Context, id string) (entity.Clinic, error) {
	clinic, err := s.repo.GetById(ctx, id)
	if err != nil {
		return entity.Clinic{}, err
	}

	return clinic, nil
}

func (s service) Create(ctx context.Context, req CreateClinicRequest) (entity.Clinic, error) {
	if err := req.Validate(); err != nil {
		return entity.Clinic{}, err
	}

	id := entity.GenerateID()
	err := s.repo.Create(ctx, entity.Clinic{
		Id:                id,
		Name:              req.Name,
		Description:       req.Description,
		Address:           req.Address,
		AppointmentPrices: make(map[string]uint),
	})

	if err != nil {
		return entity.Clinic{}, err
	}

	return s.repo.GetById(ctx, id)
}

func (s service) Update(ctx context.Context, clinicId string, req UpdateClinicRequest) (entity.Clinic, error) {
	if err := req.Validate(); err != nil {
		return entity.Clinic{}, err
	}

	clinic, err := s.repo.GetById(ctx, clinicId)
	if err != nil {
		return entity.Clinic{}, err
	}

	clinic.Name = req.Name
	clinic.Description = req.Description
	clinic.Address = req.Address

	err = s.repo.Update(ctx, clinic)
	if err != nil {
		return entity.Clinic{}, err
	}

	return clinic, nil
}

func (s service) GetAll(ctx context.Context) ([]entity.Clinic, error) {
	clinics, err := s.repo.GetAll(ctx)

	if err != nil {
		return nil, err
	}

	return clinics, nil
}

func (s service) AddAppointmentTypePrice(ctx context.Context, req AddAppointmentTypePriceRequest) error {
	clinic, err := s.repo.GetById(ctx, req.ClinicId)
	if err != nil {
		return err
	}

	appointmentType, err := s.appointmentTypeRepo.GetById(ctx, req.AppointmentTypeId)
	if err != nil {
		return err
	}

	clinic.AppointmentPrices[appointmentType.Id] = req.Price

	return nil
}

func (s service) GetAppointmentTypePrices(ctx context.Context, clinicId string) ([]entity.AppointmentTypePrice, error) {
	appointmentTypePrices, err := s.repo.GetAppointmentTypePrices(ctx, clinicId)
	if err != nil {
		return nil, err
	}

	for i, price := range appointmentTypePrices {
		appointmentType, err := s.appointmentTypeRepo.GetById(ctx, price.AppointmentTypeId)
		if err != nil {
			return nil, err
		}
		price.AppointmentType = appointmentType
		appointmentTypePrices[i] = price
	}

	return appointmentTypePrices, nil
}
