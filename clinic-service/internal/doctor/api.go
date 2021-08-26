package doctor

import (
	"net/http"

	routing "github.com/go-ozzo/ozzo-routing/v2"
	"github.com/matijapetrovic/clinichub/clinic-service/internal/errors"
	"github.com/matijapetrovic/clinichub/clinic-service/pkg/log"
)

func RegisterHandlers(r *routing.RouteGroup, service Service, authHandler routing.Handler, logger log.Logger) {
	res := resource{service, logger}

	// r.Use(authHandler)

	r.Get("/doctors", res.getAll)

	r.Post("/doctors", res.create)
}

type resource struct {
	service Service
	logger  log.Logger
}

func (r resource) getAll(c *routing.Context) error {
	doctors, err := r.service.GetAll(c.Request.Context())
	if err != nil {
		return err
	}

	return c.Write(doctors)
}

func (r resource) create(c *routing.Context) error {
	var request CreateDoctorRequest
	if err := c.Read(&request); err != nil {
		r.logger.With(c.Request.Context()).Info(err)
		return errors.BadRequest("")
	}
	doctor, err := r.service.Create((c.Request.Context()), request)
	if err != nil {
		return err
	}

	return c.WriteWithStatus(doctor, http.StatusCreated)
}
