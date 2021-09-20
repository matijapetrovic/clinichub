package appointment

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/matijapetrovic/clinichub/scheduling-service/internal/auth"
	"github.com/matijapetrovic/clinichub/scheduling-service/internal/entity"
	"github.com/matijapetrovic/clinichub/scheduling-service/pkg/httpclient"
	"github.com/matijapetrovic/clinichub/scheduling-service/pkg/log"
)

type Service interface {
	ScheduleAppointment(request *http.Request, req ScheduleAppointmentRequest) (entity.Appointment, error)
	GetDoctorAppointments(ctx context.Context, req GetDoctorAppointmentsRequest) ([]entity.Appointment, error)
	GetPatientAppointments(request *http.Request, req GetPatientAppointmentsRequest) ([]entity.Appointment, error)
	GetClinicProfit(ctx context.Context, req GetClinicReportRequest) (int, error)
}

type GetClinicReportRequest struct {
	ClinicId  string `json:"clinicId"`
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
}

func (m GetClinicReportRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.ClinicId, validation.Required, validation.Length(36, 36)),
	)
}

type ScheduleAppointmentRequest struct {
	DoctorId string    `json:"doctorId"`
	Time     time.Time `json:"time"`
}

func (m ScheduleAppointmentRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.DoctorId, validation.Length(36, 36)),
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

type GetPatientAppointmentsRequest struct {
	PatientId string `json:"patientId"`
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
}

func (m GetPatientAppointmentsRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.PatientId, validation.Required, validation.Length(36, 36)),
		validation.Field(&m.StartDate, validation.Required, validation.Date("2021-09-04")),
		validation.Field(&m.EndDate, validation.Required, validation.Date("2021-09-04")),
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

type Clinic struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type AppointmentType struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func (s service) GetClinicProfit(ctx context.Context, req GetClinicReportRequest) (int, error) {
	startDate, err := parseDate(req.StartDate)
	if err != nil {
		return -1, err
	}
	endDate, err := parseDate(req.EndDate)
	if err != nil {
		return -1, err
	}

	profit, err := s.repo.GetClinicProfit(ctx, req.ClinicId, startDate, endDate)
	if err != nil {
		return -1, err
	}

	return profit, nil
}

func (s service) ScheduleAppointment(request *http.Request, req ScheduleAppointmentRequest) (entity.Appointment, error) {
	ctx := request.Context()
	if err := req.Validate(); err != nil {
		return entity.Appointment{}, err
	}

	doctor, err := getDoctor(request.Header.Get("Authorization"), req.DoctorId)
	if err != nil {
		return entity.Appointment{}, err
	}
	_, err = s.repo.GetByDoctorIdAndTime(ctx, req.DoctorId, req.Time)
	if err == nil {
		return entity.Appointment{}, errors.New("conflict")
	} else if err != sql.ErrNoRows {
		return entity.Appointment{}, err
	}
	user := auth.CurrentUser(ctx)
	id := entity.GenerateID()
	err = s.repo.Create(ctx, entity.Appointment{
		Id:                id,
		DoctorId:          req.DoctorId,
		ClinicId:          doctor.ClinicId,
		AppointmentTypeId: doctor.AppointmentType.Id,
		PatientId:         user.GetID(),
		Price:             int(doctor.AppointmentTypePrice),
		Time:              req.Time,
	})

	if err != nil {
		return entity.Appointment{}, err
	}

	return s.repo.GetById(ctx, id)
}

func getDoctor(token string, doctorId string) (Doctor, error) {
	url, err := url.Parse("http://localhost:8081/v1/doctors/" + doctorId)
	if err != nil {
		return Doctor{}, err
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
		token,
		nil,
	)

	res, err := client.Endpoint()(context.Background(), struct{}{})
	if err != nil {
		return Doctor{}, err
	}
	doctor, ok := res.(Doctor)
	if !ok {
		return Doctor{}, errors.New("unexpected error")
	}

	return doctor, nil
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

func (s service) GetPatientAppointments(request *http.Request, req GetPatientAppointmentsRequest) ([]entity.Appointment, error) {
	ctx := request.Context()
	startDate, err := parseDate(req.StartDate)
	if err != nil {
		return nil, err
	}
	endDate, err := parseDate(req.EndDate)
	if err != nil {
		return nil, err
	}

	appointments, err := s.repo.GetByPatientIdAndDate(ctx, req.PatientId, startDate, endDate)
	if err != nil {
		return nil, err
	}

	for idx, appointment := range appointments {
		doctor, err := getDoctor(request.Header.Get("Authorization"), appointment.DoctorId)
		if err != nil {
			return nil, err
		}
		appointment.DoctorFullName = doctor.FirstName + " " + doctor.LastName

		appointments[idx] = appointment
	}

	return appointments, nil
}

func parseDate(date string) (time.Time, error) {
	if date == "" {
		return time.Time{}, nil
	}
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
