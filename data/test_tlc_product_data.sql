DO $$
    BEGIN
        INSERT INTO tealacarte.products (name, price, SKU, path)
        VALUES ('golang', 10, '0123456791abc', '/client/images/gopher.jpg'),
               ('psql', 100, '0123456792abc', '/client/images/psql.jpg'),
               ('docker', 1000, '0123456793abc', '/client/images/Moby-logo.png');
        EXCEPTION WHEN undefined_table THEN
            RAISE NOTICE 'tealacarte.products undef.  Please init';
    END;
$$;
