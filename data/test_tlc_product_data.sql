DO $$
    BEGIN
        INSERT INTO tealacarte.products (name, price, SKU)
        VALUES ('test4', 11, '0123456791abc'),
               ('test5', 12, '0123456792abc'),
               ('test6', 13, '0123456793abc'),
               ('test7', 14, '0123456794abc'),
               ('test8', 15, '0123456795abc');
        EXCEPTION WHEN undefined_table THEN
            RAISE NOTICE 'tealacarte.products undef.  Please init';
    END;
$$;
