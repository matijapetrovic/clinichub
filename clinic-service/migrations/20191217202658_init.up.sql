CREATE TABLE clinic (
  id  VARCHAR(255) NOT NULL,
  name VARCHAR(255) NOT NULL,
  description VARCHAR(255) NOT NULL,
  address_line VARCHAR(255) NOT NULL,
  city VARCHAR(255) NOT NULL,
  country VARCHAR(255) NOT NULL,

  PRIMARY KEY (`id`)
);

CREATE TABLE appointment_type (
  id VARCHAR(255) NOT NULL,
  name VARCHAR(255) NOT NULL,

  PRIMARY KEY (`id`)
);

CREATE TABLE appointment_type_price (
  clinic_id VARCHAR(255) NOT NULL,
  appointment_type_id VARCHAR(255) NOT NULL,
  price INT NOT NULL,

  PRIMARY KEY (`clinic_id`, `appointment_type_id`),
  FOREIGN KEY (`clinic_id`) REFERENCES clinic(`id`),
  FOREIGN KEY (`appointment_type_id`) REFERENCES appointment_type(`id`)
);

CREATE TABLE doctor (
  id VARCHAR(255) NOT NULL,
  first_name VARCHAR(255) NOT NULL,
  last_name VARCHAR(255) NOT NULL,
  work_start_hour INT NOT NULL,
  work_start_minute INT NOT NULL,
  work_end_hour INT NOT NULL,
  work_end_minute INT NOT NULL,
  clinic_id VARCHAR(255) NOT NULL,
  specialization_id VARCHAR(255) NOT NULL,

  PRIMARY KEY (`id`),
  FOREIGN KEY (`clinic_id`) REFERENCES clinic(`id`),
  FOREIGN KEY (`specialization_id`) REFERENCES appointment_type(`id`)
);