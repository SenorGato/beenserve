with import <nixpkgs> {};
stdenv.mkDerivation {
  name = "beenserve-test";
  nativeBuildInputs = [];
  buildInputs = [];
  shellHook = ''
    export POSTGRES_LOGIN="tealacarte-test";
    export POSTGRES_PASSWORD="test123";
    export DATABASE_URL="postgres://$POSTGRES_LOGIN:$POSTGRES_PASSWORD@localhost:5432/tealacarte";
    export STRIPE_KEY="sk_test_51MNgItJUna26uIQEc7yGt2dYnwLjWOrpRSEsnITSK87j3Ff0BB5N7aKs1eOKYwmwEaRNIAnUD7Wz7IWLstq3ovku00vLwGPfEW";
    export TLC_VERSION="test";
    export PRODUCT_DATABASE_INIT="tlc_product_init.sql"
    export PRODUCT_DATA_INIT="test_tlc_product_data.sql"
    '';
}
