package entity

type DoctorRating struct {
	ID        string  `json:"id"`
	Rating    float32 `json:"rating"`
	PatientId string  `json:"patientId"`
	DoctorId  string  `json:"clinicId"`
}

type ClinicRating struct {
	ID        string  `json:"id"`
	Rating    float32 `json:"rating"`
	PatientId string  `json:"patientId"`
	ClinicId  string  `json:"clinicId"`
}

type AverageRating struct {
	Rating float32 `json:"rating"`
	Count  int     `json:"count"`
}
