package entity

type Clinic struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Address     `json:"address"`
}

type Address struct {
	AddressLine string `json:"addressLine"`
	City        string `json:"city"`
	Country     string `json:"country"`
}

type AppointmentTypePrice struct {
	ClinicId          string `json:"clinicId" db:"pk"`
	AppointmentTypeId string `json:"-" db:"pk"`
	AppointmentType   `json:"appointmentType" db:"-"`
	Price             uint `json:"price"`
}
