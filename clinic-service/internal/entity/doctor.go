package entity

type Doctor struct {
	Id             string          `json:"doctorId"`
	Clinic         Clinic          `json:"clinic"`
	FirstName      string          `json:"firstName"`
	LastName       string          `json:"lastName"`
	WorkStart      Time            `json:"workStart"`
	WorkEnd        Time            `json:"workEnd"`
	Specialization AppointmentType `json:"specialization"`
}

type Time struct {
	Hour   uint `json:"hour"`
	Minute uint `json:"minute"`
}
