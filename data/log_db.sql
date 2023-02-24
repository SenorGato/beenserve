--LogDB

CREATE SCHEMA IF NOT EXISTS logdb;
SET search_path TO logdb;
REVOKE ALL ON SCHEMA logdb FROM public;
REVOKE ALL ON DATABASE logdb FROM public;
GRANT ALL PRIVILEGES ON DATABASE logdb TO logdb;

CREATE OR REPLACE FUNCTION trigger_set_modify_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TABLE IF NOT EXISTS logdb.logs(
    log_entry varchar,
    created_on  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE OR REPLACE TRIGGER Set_timestamp
BEFORE UPDATE ON userauth.users
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_modify_timestamp();
