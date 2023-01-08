--Tea a la carte DB Provisioning
-- Set permissions, and default schema
-- Add Customer and Admin view

CREATE SCHEMA IF NOT EXISTS tealacarte;
SET search_path TO tealacarte;
REVOKE ALL ON SCHEMA tealacarte FROM public;
REVOKE ALL ON DATABASE tealacarte FROM public;
GRANT ALL PRIVILEGES ON DATABASE tealacarte TO tealacarte;
ALTER USER tealacarte WITH PASSWORD 'smoke';

DO $$
BEGIN
    CREATE DOMAIN sku AS TEXT CHECK (VALUE ~* '^[A-Za-z0-9]{13}$');
    EXCEPTION WHEN duplicate_object THEN 
        RAISE NOTICE 'domain sku already exists';
END;
$$;

CREATE OR REPLACE FUNCTION trigger_set_modify_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TABLE IF NOT EXISTS tealacarte.products (
    product_id  SERIAL PRIMARY KEY,
    name        TEXT,
    price       DECIMAL,
    SKU         sku,
    created_on  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);


CREATE OR REPLACE TRIGGER Set_timestamp
BEFORE UPDATE ON tealacarte.products
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_modify_timestamp();
