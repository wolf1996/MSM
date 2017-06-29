DROP EXTENSION IF EXISTS CITEXT CASCADE ;

DROP TABLE IF EXISTS USERS CASCADE ;
DROP TABLE IF EXISTS CONTROLLERS CASCADE ;
DROP TABLE IF EXISTS SENSOR CASCADE ;
DROP TABLE IF EXISTS DATA CASCADE ;

CREATE EXTENSION IF NOT EXISTS CITEXT WITH SCHEMA public;

CREATE TABLE IF NOT EXISTS USERS(
  id SERIAL PRIMARY KEY,
  family_name VARCHAR(256),
  name VARCHAR(256),
  second_name VARCHAR(256) DEFAULT  NULL ,
  date_receiving DATE,
  issued_by TEXT,
  division_number VARCHAR(50),
  registration_addres TEXT,
  mailing_addres TEXT,
  home_phone VARCHAR(20),
  mobile_phone VARCHAR(20),
  citizenship VARCHAR(256),
  e_mail VARCHAR(50),
  pass_hash VARCHAR(256)
);

CREATE TABLE IF NOT EXISTS CONTROLLERS(
  id SERIAL PRIMARY KEY,
  name VARCHAR(256),
  user_id INT REFERENCES USERS(id),
  addres TEXT,
  activation_date DATE,
  status INT,
  mac MACADDR,
  deactivation_date DATE,
  controller_type INT
);

CREATE TABLE IF NOT EXISTS SENSOR(
  id SERIAL PRIMARY KEY,
  name VARCHAR(256),
  controller_id INT REFERENCES CONTROLLERS(id),
  activation_date DATE,
  status INT,
  deactivation_date DATE,
  sensor_type INT,
  company  VARCHAR(256)
);

CREATE TABLE IF NOT EXISTS DATA(
  sensor_id INT REFERENCES SENSOR(id),
  date DATE,
  value BIGINT,
  hs UUID
);


CREATE TABLE IF NOT EXISTS TAX (
  id   SERIAL PRIMARY KEY,
  name VARCHAR(256),
  TAX  FLOAT
);

INSERT INTO  TAX VALUES (
  1,
  'Флэш',
  2.5
);

-- INSERT INTO DATA VALUES (
--     1,
--     '1961-06-16',
--     1000,
--     'ff62d9d0cc926a7516e408b4ad1a0537'
-- );


INSERT INTO USERS VALUES (
    1,
  'Иванов',
  'Иван',
  'Иванович',
  '1961-06-16',
  '1961-06-16',
  '11111',
  'Улица Пушкина, Дом Колотушкина',
  'Улица Пушкина, Дом Колотушкина',
  '111 555',
  '8 800 555 35 35',
  'Албания',
  'ml@gmail.com',
  '123456'
);

INSERT INTO CONTROLLERS VALUES (
  1,
  'test_controller',
  1,
  'Улица Пушкина, Дом Колотушкина',
  '2001-01-01',
  1,
  '6B-45-CD-97-48-48',
  NULL,
  1
);


INSERT INTO SENSOR VALUES (
  1,
  'test_sensor',
  1,
  '2001-01-01',
  1,
  NULL,
  1,
  'GASPROM'
);

SELECT
SENSOR.id, SENSOR.Name, SENSOR.controller_id, SENSOR.activation_date, SENSOR.status, SENSOR.deactivation_date, SENSOR.sensor_type, SENSOR.company
FROM SENSOR INNER JOIN CONTROLLERS ON CONTROLLERS.id = SENSOR.controller_id
WHERE controller_id = 1 AND user_id = 1;

SELECT
DATA.sensor_id, DATA.date, DATA.value, DATA.hs
FROM DATA INNER JOIN SENSOR ON DATA.sensor_id = SENSOR.id
INNER JOIN CONTROLLERS ON SENSOR.controller_id = CONTROLLERS.id
WHERE sensor_id = 1 AND user_id = 1 AND date > '1961-06-15'
LIMIT 100;

GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO iot_api_user;
