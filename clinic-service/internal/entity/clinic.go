package entity

type Clinic struct {
	Id                string `json:"clinicId"`
	Name              string `json:"name"`
	Description       string `json:"description"`
	Address           `json:"address"`
	AppointmentPrices map[string]uint `json:"appointmentPrices" db:"-"`
}

type Address struct {
	AddressLine string `json:"addressLine"`
	City        string `json:"city"`
	Country     string `json:"country"`
}

type AppointmentTypePrice struct {
	ClinicId          string `json:"clinicId"`
	AppointmentTypeId string `json:"-" db:"appointment_type_id"`
	AppointmentType   `json:"appointmentType"`
	Price             uint `json:"price"`
}
