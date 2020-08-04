CREATE EXTENSION pgcrypto;

CREATE TYPE role AS ENUM ('super_admin', 'admin', 'sub_admin', 'sale_person');
CREATE TYPE type AS ENUM ('mobile', 'home', 'work');


CREATE TABLE merchants (
  id BIGSERIAL NOT NULL PRIMARY KEY,
  number_of_product INTEGER,
  role_type role NOT NULL DEFAULT 'admin',
  email VARCHAR(300) NOT NULL UNIQUE,
  phone VARCHAR(15) NOT NULL UNIQUE,
  phone_type type NOT NULL DEFAULT 'mobile',
  user_id BIGINT NOT NULL UNIQUE,
  number_of_outlet INTEGER,
  business_name VARCHAR(100) NOT NULL,
  is_suspended BOOLEAN DEFAULT FALSE,
  is_email_verified BOOLEAN DEFAULT FALSE,
  salt VARCHAR(20) NOT NULL,
  passwd VARCHAR(300),
  last_login TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
  deleted_at TIMESTAMP WITH TIME ZONE
)