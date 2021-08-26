package entity

type Clinic struct {
	Id                string          `json:"clinicId"`
	Name              string          `json:"name"`
	Description       string          `json:"description"`
	Address           Address         `json:"address"`
	Doctors           []Doctor        `json:"doctors"`
	AppointmentPrices map[string]uint `json:"appointmentPrices"`
}

type Address struct {
	AddressLine string `json:"addressLine"`
	City        string `json:"city"`
	Country     string `json:"country"`
}
