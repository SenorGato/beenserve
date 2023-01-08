docker run --name postgres-db -e POSTGRES_PASSWORD=docker -p 5432:5432 -d postgres
docker cp teaalacarte.sql a85bca48a585:/
docker exec -it a85bca48a585 psql -U postgres -W postgres -a -f teaalacarte.sql

REVOKE ALL ON DATABASE foo FROM publc;
GRANT role to role;
GRANT CONNNECT ON DATABASE foo TO role;
REVOKE ALL ON SCHEMA public FROM public;

docker build -t tealacarte:1.0

ToDo:
Fix secrets properly
Password for tealacarte not getting set properly
