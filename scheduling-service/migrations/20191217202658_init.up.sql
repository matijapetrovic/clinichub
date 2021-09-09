CREATE TABLE appointment (
  id VARCHAR(255) NOT NULL,
  clinic_id VARCHAR(255) NOT NULL,
  doctor_id VARCHAR(255) NOT NULL,
  patient_id VARCHAR(255) NOT NULL,
  appointment_type_id VARCHAR(255) NOT NULL,
  price INT NOT NULL,
  time DATETIME NOT NULL,
  PRIMARY KEY (`id`)
);