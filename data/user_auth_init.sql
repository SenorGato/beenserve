--User Database Schema

CREATE SCHEMA IF NOT EXISTS userauth;
SET search_path TO userauth;
REVOKE ALL ON SCHEMA userauth FROM public;
REVOKE ALL ON DATABASE userauth FROM public;
GRANT ALL PRIVILEGES ON DATABASE userauth TO userauth;

DO $$
BEGIN
    CREATE EXTENSION citext;
    CREATE DOMAIN email AS citext 
    CHECK ( value ~ '^[a-zA-Z0-9.!#$%&''*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$' );
    EXCEPTION WHEN duplicate_object THEN 
        RAISE NOTICE 'domain email already exists';
END;
$$;

CREATE OR REPLACE FUNCTION trigger_set_modify_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TABLE IF NOT EXISTS userauth.users(
    email       email PRIMARY KEY,
    name        TEXT,
    pass_hash varchar,
    api_hash varchar,
    test_api_hash varchar,
    created_on  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE OR REPLACE TRIGGER Set_timestamp
BEFORE UPDATE ON userauth.users
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_modify_timestamp();
