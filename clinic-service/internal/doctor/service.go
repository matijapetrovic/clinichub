package doctor

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"sort"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	appointment_type "github.com/matijapetrovic/clinichub/clinic-service/internal/appointment-type"
	"github.com/matijapetrovic/clinichub/clinic-service/internal/clinic"
	"github.com/matijapetrovic/clinichub/clinic-service/internal/entity"
	"github.com/matijapetrovic/clinichub/clinic-service/pkg/httpclient"
	"github.com/matijapetrovic/clinichub/clinic-service/pkg/log"
)

type Service interface {
	GetById(ctx context.Context, id string) (entity.Doctor, error)
	GetAll(ctx context.Context) ([]entity.Doctor, error)
	GetByClinicId(request *http.Request, clinicId string, req GetByClinicIdRequest) ([]entity.Doctor, error)
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

type GetByClinicIdRequest struct {
	AppointmentTypeId string `json:"appointmentTypeId"`
	Date              string `json:"date"`
}

func (m GetByClinicIdRequest) Validate() error {
	return validation.ValidateStruct(&m,

		validation.Field(&m.AppointmentTypeId, validation.Length(36, 36)),
		validation.Field(&m.Date, validation.Required, validation.Date("30-12-1970"), validation.Min(time.Now())),
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

func (s service) GetById(ctx context.Context, id string) (entity.Doctor, error) {
	doctor, err := s.repo.GetById(ctx, id)
	if err != nil {
		return entity.Doctor{}, err
	}

	appointmentPrice, err := s.clinicRepo.GetAppointmentTypePrice(ctx, doctor.ClinicId, doctor.SpecializationId)
	if err != nil {
		return entity.Doctor{}, err
	}

	appointmentType, err := s.appointmentTypeRepo.GetById(ctx, doctor.SpecializationId)
	if err != nil {
		return entity.Doctor{}, err
	}

	doctor.AppointmentType = appointmentType
	doctor.AppointmentTypePrice = appointmentPrice.Price

	return doctor, nil
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

type Appointment struct {
	Id                string    `json:"id"`
	ClinicId          string    `json:"clinicId"`
	DoctorId          string    `json:"doctorId"`
	PatientId         string    `json:"patientId"`
	AppointmentTypeId string    `json:"appointmentTypeId"`
	Price             uint      `json:"price"`
	Time              time.Time `json:"time"`
}

func (s service) GetByClinicId(request *http.Request, clinicId string, req GetByClinicIdRequest) ([]entity.Doctor, error) {
	if req.AppointmentTypeId == "" {
		doctors, err := s.repo.GetByClinicId(request.Context(), clinicId)
		if err != nil {
			return nil, err
		}

		for i, doctor := range doctors {
			specialization, err := s.appointmentTypeRepo.GetById(request.Context(), doctor.SpecializationId)
			if err != nil {
				return nil, err
			}
			doctor.AppointmentType = specialization
			doctors[i] = doctor
		}

		return doctors, nil
	}
	doctors, err := s.repo.GetByClinicIdAndSpecializationId(request.Context(), clinicId, req.AppointmentTypeId)
	if err != nil {
		return nil, err
	}

	appointmentPrice, err := s.clinicRepo.GetAppointmentTypePrice(request.Context(), clinicId, req.AppointmentTypeId)
	if err != nil {
		return nil, err
	}

	for i, doctor := range doctors {
		appointments, err := getDoctorAppointments(doctor.Id, req.Date, request.Header.Get("Authorization"))
		if err != nil {
			return nil, err
		}

		workStart, _ := entity.ParseTime(doctor.WorkStart)
		workEnd, _ := entity.ParseTime(doctor.WorkEnd)
		workingHours := entity.GetHours(workStart, workEnd)

		for _, appointment := range appointments {
			hour, _, _ := appointment.Time.Clock()
			delete(workingHours, uint(hour))
		}

		sortedWorkingHours := make([]string, 0, len(workingHours))
		for k, _ := range workingHours {
			time := entity.Time{Hour: k, Minute: 0}.ToString()
			sortedWorkingHours = append(sortedWorkingHours, time)
		}

		sort.Strings(sortedWorkingHours)
		doctor.AvailableHours = sortedWorkingHours
		doctor.AppointmentTypePrice = appointmentPrice.Price

		rating, err := getDoctorRating(doctor.Id, request.Header.Get("Authorization"))
		if err != nil {
			return nil, err
		}
		doctor.Rating = rating
		doctors[i] = doctor
	}

	return doctors, nil
}

func getDoctorAppointments(doctorId string, date string, token string) ([]Appointment, error) {
	url, err := url.Parse("http://localhost:8083/v1/doctors/" + doctorId + "/appointments")
	if err != nil {
		return nil, err
	}

	queryParamMap := map[string]string{
		"date": date,
	}

	client := httpclient.NewJsonClient(
		"GET",
		url,
		func(ctx context.Context, r *http.Response) (interface{}, error) {
			var appointment []Appointment
			err := json.NewDecoder(r.Body).Decode(&appointment)
			if err != nil {
				return nil, err
			}
			return appointment, nil
		},
		token,
		httpclient.QueryParamBeforeFunc(queryParamMap),
	)

	res, err := client.Endpoint()(context.Background(), struct{}{})
	if err != nil {
		return nil, err
	}
	appointments, ok := res.([]Appointment)
	if !ok {
		return nil, errors.New("unexpected error")
	}

	return appointments, nil
}

func getDoctorRating(doctorId string, token string) (entity.Rating, error) {
	url, err := url.Parse("http://localhost:8082/v1/doctors/" + doctorId + "/average-rating")
	if err != nil {
		return entity.Rating{}, err
	}

	client := httpclient.NewJsonClient(
		"GET",
		url,
		func(ctx context.Context, r *http.Response) (interface{}, error) {
			var rating entity.Rating
			err := json.NewDecoder(r.Body).Decode(&rating)
			if err != nil {
				return nil, err
			}
			return rating, nil
		},
		token,
		nil,
	)

	res, err := client.Endpoint()(context.Background(), struct{}{})
	if err != nil {
		return entity.Rating{}, err
	}
	rating, ok := res.(entity.Rating)
	if !ok {
		return entity.Rating{}, errors.New("unexpected error")
	}

	return rating, nil
}

func (s service) GetAll(ctx context.Context) ([]entity.Doctor, error) {
	doctors, err := s.repo.GetAll(ctx)

	if err != nil {
		return nil, err
	}

	return doctors, nil
}
