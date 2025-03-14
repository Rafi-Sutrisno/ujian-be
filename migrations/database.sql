CREATE DATABASE golang_template;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
  id          UUID PRIMARY KEY DEFAULT uuid_generate_v4() NOT NULL,
  name        VARCHAR(100) NOT NULL,
  email       VARCHAR(100) NULL,
  password    VARCHAR(100) NOT NULL,
  role        VARCHAR(100) NOT NULL,
  noid        VARCHAR(100) NOT NULL,
  created_at  TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  updated_at  TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

DROP DATABASE golang_template;