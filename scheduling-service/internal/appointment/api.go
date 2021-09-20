package appointment

import (
	"net/http"

	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/matijapetrovic/clinichub/scheduling-service/internal/auth"
	"github.com/matijapetrovic/clinichub/scheduling-service/internal/errors"
	"github.com/matijapetrovic/clinichub/scheduling-service/pkg/log"
)

func RegisterHandlers(r *routing.RouteGroup, service Service, authHandler routing.Handler, logger log.Logger) {
	res := resource{service, logger}

	r.Use(authHandler)

	r.Get("/appointments", res.query)
	r.Get("/doctors/<id>/appointments", res.getDoctorAppointments)
	r.Post("/appointments", res.schedule)
	r.Get("/clinics/<id>/profit", res.getProfit)
}

type resource struct {
	service Service
	logger  log.Logger
}

func (r resource) getProfit(c *routing.Context) error {
	startDate := c.Request.URL.Query().Get("startDate")
	endDate := c.Request.URL.Query().Get("endDate")

	profit, err := r.service.GetClinicProfit(c.Request.Context(), GetClinicReportRequest{
		ClinicId:  c.Param("id"),
		StartDate: startDate,
		EndDate:   endDate,
	})

	if err != nil {
		return err
	}

	return c.Write(profit)
}

func (r resource) query(c *routing.Context) error {
	user := auth.CurrentUser(c.Request.Context())
	startDate := c.Request.URL.Query().Get("startDate")
	endDate := c.Request.URL.Query().Get("endDate")
	appointments, err := r.service.GetPatientAppointments(c.Request, GetPatientAppointmentsRequest{
		PatientId: user.GetID(),
		StartDate: startDate,
		EndDate:   endDate,
	})
	if err != nil {
		return err
	}

	return c.Write(appointments)
}

func (r resource) getDoctorAppointments(c *routing.Context) error {
	appointments, err := r.service.GetDoctorAppointments(c.Request.Context(), GetDoctorAppointmentsRequest{
		Date:     c.Request.URL.Query().Get("date"),
		DoctorId: c.Param("id"),
	})

	if err != nil {
		return err
	}

	// if err != nil {
	// 	if err.Error() == "bad request" {
	// 		return c.WriteWithStatus(err.Error(), http.StatusBadRequest)
	// 	}
	// 	return err
	// }

	return c.Write(appointments)
}

func (r resource) schedule(c *routing.Context) error {
	var request ScheduleAppointmentRequest
	if err := c.Read(&request); err != nil {
		r.logger.With(c.Request.Context()).Info(err)
		return errors.BadRequest("")
	}
	appointment, err := r.service.ScheduleAppointment(c.Request, request)
	if err != nil {
		if err.Error() == "conflict" {
			return c.WriteWithStatus(err.Error(), http.StatusConflict)
		} else if err.Error() == "unexpected" {
			return c.WriteWithStatus(err.Error(), http.StatusInternalServerError)
		}
		return err
	}

	return c.WriteWithStatus(appointment, http.StatusCreated)
}
