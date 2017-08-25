DROP EXTENSION IF EXISTS CITEXT CASCADE;

DROP TABLE IF EXISTS USERS CASCADE;
DROP TABLE IF EXISTS CONTROLLERS CASCADE;
DROP TABLE IF EXISTS SENSOR CASCADE;
DROP TABLE IF EXISTS DATA CASCADE;
DROP TABLE IF EXISTS TAX CASCADE;

CREATE EXTENSION IF NOT EXISTS CITEXT WITH SCHEMA public;

CREATE TABLE IF NOT EXISTS USERS (
  id                   SERIAL PRIMARY KEY,
  family_name          VARCHAR(256)       NOT NULL,
  name                 VARCHAR(256)       NOT NULL,
  second_name          VARCHAR(256) DEFAULT NULL,
  date_receiving       DATE         DEFAULT NOW(),
  issued_by            TEXT         DEFAULT NULL,
  division_number      VARCHAR(50)  DEFAULT NULL,
  registration_address TEXT         DEFAULT NULL,
  mailing_address      TEXT         DEFAULT NULL,
  home_phone           VARCHAR(20)  DEFAULT NULL,
  mobile_phone         VARCHAR(20)  DEFAULT NULL,
  citizenship          VARCHAR(256) DEFAULT NULL,
  email                VARCHAR(50) UNIQUE NOT NULL,
  pass_hash            VARCHAR(256)       NOT NULL
);

CREATE TABLE IF NOT EXISTS CONTROLLERS (
  id                SERIAL PRIMARY KEY,
  name              VARCHAR(256)              NOT NULL UNIQUE,
  user_id           INT REFERENCES USERS (id) NOT NULL,
  address           TEXT                      NOT NULL,
  activation_date   DATE DEFAULT NOW(),
  status            INT  DEFAULT NULL,
  mac               MACADDR                   NOT NULL,
  deactivation_date DATE DEFAULT NOW(),
  controller_type   INT  DEFAULT NULL
);


CREATE TABLE IF NOT EXISTS TAX (
  id   SERIAL PRIMARY KEY,
  name VARCHAR(256) NOT NULL,
  TAX  FLOAT        NOT NULL
);


CREATE TABLE IF NOT EXISTS SENSOR (
  id                SERIAL PRIMARY KEY,
  name              VARCHAR(256)                    NOT NULL UNIQUE,
  controller_id     INT REFERENCES CONTROLLERS (id) NOT NULL,
  activation_date   DATE                    DEFAULT NOW(),
  status            INT                     DEFAULT NULL,
  deactivation_date DATE                    DEFAULT NOW(),
  sensor_type       INT                     DEFAULT NULL,
  company           VARCHAR(256)            DEFAULT NULL,
  tax               INT REFERENCES TAX (id) DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS DATA (
  sensor_id INT REFERENCES SENSOR (id) NOT NULL,
  date      DATE   DEFAULT NOW(),
  value     BIGINT DEFAULT NULL,
  hs        UUID   DEFAULT NULL
);

INSERT INTO TAX VALUES (
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
  'GASPROM',
  1
);

SELECT
  SENSOR.id,
  SENSOR.Name,
  SENSOR.controller_id,
  SENSOR.activation_date,
  SENSOR.status,
  SENSOR.deactivation_date,
  SENSOR.sensor_type,
  SENSOR.company
FROM SENSOR
  INNER JOIN CONTROLLERS ON CONTROLLERS.id = SENSOR.controller_id
WHERE controller_id = 1 AND user_id = 1;

SELECT
  DATA.sensor_id,
  DATA.date,
  DATA.value,
  DATA.hs
FROM DATA
  INNER JOIN SENSOR ON DATA.sensor_id = SENSOR.id
  INNER JOIN CONTROLLERS ON SENSOR.controller_id = CONTROLLERS.id
WHERE sensor_id = 1 AND user_id = 1 AND date > '1961-06-15'
LIMIT 100;

SELECT max(DATA.value) - min(DATA.value)
FROM DATA
  INNER JOIN SENSOR ON DATA.sensor_id = SENSOR.id
  INNER JOIN CONTROLLERS ON SENSOR.controller_id = CONTROLLERS.id
WHERE sensor_id = 1 AND user_id = 1 AND date >= '01-01-2017'
      AND date < '01-02-2017';

WITH mnths AS (SELECT
                 extract(MONTH FROM date)            AS mnth,
                 (max(DATA.value) - min(DATA.value)) AS value
               FROM DATA
                 INNER JOIN SENSOR ON DATA.sensor_id = SENSOR.id
                 INNER JOIN CONTROLLERS ON SENSOR.controller_id = CONTROLLERS.id
               WHERE sensor_id = 1 AND user_id = 1 AND date < '12-12-2017'
                     AND date > '12-12-2016'
               GROUP BY 1
)
SELECT avg(value)
FROM mnths;

SELECT *
FROM DATA;


GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO iot_api_user;
