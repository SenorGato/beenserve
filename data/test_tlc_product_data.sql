DO $$
    BEGIN
        INSERT INTO tealacarte.products (name, price, SKU, path)
        VALUES ('test4', 19, '0123456791abc', '/home'),
               ('test5', 18, '0123456792abc', '/home'),
               ('test6', 17, '0123456793abc', '/home'),
               ('test7', 16, '0123456794abc', '/home'),
               ('test8', 12, '0123456795abc', '/home');
        EXCEPTION WHEN undefined_table THEN
            RAISE NOTICE 'tealacarte.products undef.  Please init';
    END;
$$;
