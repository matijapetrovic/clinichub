package appointment

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/matijapetrovic/clinichub/scheduling-service/internal/entity"
	"github.com/matijapetrovic/clinichub/scheduling-service/pkg/httpclient"
	"github.com/matijapetrovic/clinichub/scheduling-service/pkg/log"
)

type Service interface {
	ScheduleAppointment(request *http.Request, req ScheduleAppointmentRequest) (entity.Appointment, error)
	GetDoctorAppointments(ctx context.Context, req GetDoctorAppointmentsRequest) ([]entity.Appointment, error)
}

type ScheduleAppointmentRequest struct {
	DoctorId  string    `json:"doctorId"`
	PatientId string    `json:"patientId"`
	Time      time.Time `json:"time"`
}

func (m ScheduleAppointmentRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.DoctorId, validation.Length(36, 36)),
		validation.Field(&m.PatientId, validation.Length(36, 36)),
		validation.Field(&m.Time, validation.Min(time.Now())),
	)
}

type GetDoctorAppointmentsRequest struct {
	DoctorId string `json:"doctorId"`
	Date     string `json:"date"`
}

func (m GetDoctorAppointmentsRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.DoctorId, validation.Required, validation.Length(36, 36)),
		//		validation.Field(&m.Date, validation.Required, validation.Date("2021-09-04"), validation.Min(time.Now())),
	)
}

type service struct {
	repo   Repository
	logger log.Logger
}

func NewService(repo Repository, logger log.Logger) Service {
	return service{repo, logger}
}

type Doctor struct {
	Id                   string `json:"id"`
	ClinicId             string `json:"clinicId"`
	FirstName            string `json:"firstName"`
	LastName             string `json:"lastName"`
	WorkStart            string `json:"workStart"`
	WorkEnd              string `json:"workEnd"`
	AppointmentType      `json:"specialization"`
	AppointmentTypePrice uint     `json:"specializationPrice"`
	AvailableHours       []string `json:"availableHours"`
}

type AppointmentType struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func (s service) ScheduleAppointment(request *http.Request, req ScheduleAppointmentRequest) (entity.Appointment, error) {
	ctx := request.Context()
	if err := req.Validate(); err != nil {
		return entity.Appointment{}, err
	}

	url, err := url.Parse("http://localhost:8081/v1/doctors/" + req.DoctorId)
	if err != nil {
		return entity.Appointment{}, err
	}

	client := httpclient.NewJsonClient(
		"GET",
		url,
		func(ctx context.Context, r *http.Response) (interface{}, error) {
			var doctor Doctor
			err := json.NewDecoder(r.Body).Decode(&doctor)
			if err != nil {
				return nil, err
			}
			return doctor, nil
		},
		request.Header.Get("Authorization"),
		nil,
	)

	res, err := client.Endpoint()(context.Background(), struct{}{})
	if err != nil {
		return entity.Appointment{}, err
	}

	doctor, ok := res.(Doctor)
	if !ok {
		return entity.Appointment{}, errors.New("unexpected error")
	}

	fmt.Printf("%v", doctor)

	_, err = s.repo.GetByDoctorIdAndTime(ctx, req.DoctorId, req.Time)
	if err == nil {
		return entity.Appointment{}, errors.New("conflict")
	} else if err != sql.ErrNoRows {
		return entity.Appointment{}, err
	}

	id := entity.GenerateID()
	err = s.repo.Create(ctx, entity.Appointment{
		Id:                id,
		DoctorId:          req.DoctorId,
		ClinicId:          doctor.ClinicId,
		AppointmentTypeId: doctor.AppointmentType.Id,
		PatientId:         req.PatientId,
		Price:             int(doctor.AppointmentTypePrice),
		Time:              req.Time,
	})

	if err != nil {
		return entity.Appointment{}, err
	}

	return s.repo.GetById(ctx, id)
}

func (s service) GetDoctorAppointments(ctx context.Context, req GetDoctorAppointmentsRequest) ([]entity.Appointment, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	date, err := parseDate(req.Date)
	if err != nil {
		return nil, err
	}

	appointments, err := s.repo.GetDoctorAppointments(ctx, req.DoctorId, date, date.Add(24*time.Hour))
	if err != nil {
		return nil, err
	}

	return appointments, nil
}

func parseDate(date string) (time.Time, error) {
	split := strings.Split(date, "-")
	year, err := strconv.Atoi(split[0])
	if err != nil {
		return time.Time{}, err
	}

	month, err := strconv.Atoi(split[1])
	if err != nil {
		return time.Time{}, err
	}

	day, err := strconv.Atoi(split[2])
	if err != nil {
		return time.Time{}, err
	}

	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC), nil
}
