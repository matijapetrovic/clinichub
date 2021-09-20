CREATE TABLE doctor_rating
(
    id         VARCHAR PRIMARY KEY,
    rating     DECIMAL NOT NULL,
    patient_id VARCHAR NOT NULL,
    doctor_id  VARCHAR NOT NULL
);

CREATE TABLE clinic_rating
(
    id         VARCHAR PRIMARY KEY,
    rating     DECIMAL NOT NULL,
    patient_id VARCHAR NOT NULL,
    clinic_id  VARCHAR NOT NULL
);
