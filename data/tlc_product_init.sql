--Tea a la carte DB Provisioning
-- Add Customer and Admin view

-- PGOPTIONS="-c 'custom.postgres_password=${POSTGRES_PASSWORD}'";
-- ALTER USER tealacarte WITH PASSWORD (SELECT current_setting('custom.postgres_password'));

CREATE SCHEMA IF NOT EXISTS tealacarte;
SET search_path TO tealacarte;
REVOKE ALL ON SCHEMA tealacarte FROM public;
REVOKE ALL ON DATABASE tealacarte FROM public;
GRANT ALL PRIVILEGES ON DATABASE tealacarte TO tealacarte;

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
    SKU         sku UNIQUE,
    path        TEXT,
    created_on  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE VIEW frontend as 
SELECT product_id, name, price, path, SKU
FROM products;

CREATE OR REPLACE TRIGGER Set_timestamp
BEFORE UPDATE ON tealacarte.products
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_modify_timestamp();
