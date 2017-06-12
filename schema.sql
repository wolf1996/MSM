DROP EXTENSION IF EXISTS CITEXT CASCADE ;

DROP TABLE IF EXISTS USERS;

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

SELECT id, e_mail, pass_hash "
		"FROM USERS WHERE e_mail = 'ml@gmail.com' ;