package appointment

import (
	"net/http"

	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/matijapetrovic/clinichub/scheduling-service/internal/errors"
	"github.com/matijapetrovic/clinichub/scheduling-service/pkg/log"
)

func RegisterHandlers(r *routing.RouteGroup, service Service, authHandler routing.Handler, logger log.Logger) {
	res := resource{service, logger}

	// r.Use(authHandler)

	r.Get("/appointments", res.query)

	r.Post("/appointments", res.schedule)

}

type resource struct {
	service Service
	logger  log.Logger
}

func (r resource) query(c *routing.Context) error {
	clinics, err := r.service.GetDoctorAppointments(c.Request.Context(), GetDoctorAppointmentsRequest{
		Date:     c.Request.URL.Query().Get("date"),
		DoctorId: c.Request.URL.Query().Get("doctorId"),
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

	return c.Write(clinics)
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
