with import <nixpkgs> {};
stdenv.mkDerivation {
  name = "beenserve-test";
  nativeBuildInputs = [];
  buildInputs = [];
  shellHook = ''
    export POSTGRES_USER="tealacarte";
    export POSTGRES_PASSWORD="test123";
    export DATABASE_PORT="5432"
    export DATABASE_HOST="localhost"
    export DATABASE_URL="postgres://$POSTGRES_USER:$POSTGRES_PASSWORD@$DATABASE_HOST:$DATABASE_PORT/$POSTGRES_USER";
    export PRODUCT_DATABASE_INIT="tlc_product_init.sql"
    export PRODUCT_DATA_INIT="test_tlc_product_data.sql"

    export WEB_SERVER_PORT="9090"
    export STRIPE_KEY="sk_test_51MNgItJUna26uIQEc7yGt2dYnwLjWOrpRSEsnITSK87j3Ff0BB5N7aKs1eOKYwmwEaRNIAnUD7Wz7IWLstq3ovku00vLwGPfEW";
    export STRIPE_PUBLISHABLE_KEY="pk_test_51MNgItJUna26uIQEAZhCYdVAvrc0pM7qtJvKP6oe8lEbgcefGL9hEhLeoOZYaxklq0ih6enZVnwMw8DzO2VY5Tmj00njjYCafM"
    export TLC_VERSION="test";
    '';
}
