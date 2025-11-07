# !/bin/bash
set -e

# Expected environment variables:
# DB_USER, DB_PASS, DB_NAME
# TEST_DB_USER, TEST_DB_PASS, TEST_DB_NAME
# POSTGRES_USER, POSTGRES_DB (from Postgres image)

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    CREATE USER ${DB_USER} WITH PASSWORD '${DB_PASS}';
    CREATE DATABASE ${DB_NAME};
    GRANT ALL PRIVILEGES ON DATABASE ${DB_NAME} TO ${DB_USER};

    CREATE USER ${TEST_DB_USER} WITH PASSWORD '${TEST_DB_PASS}';
    CREATE DATABASE ${TEST_DB_NAME};
    GRANT ALL PRIVILEGES ON DATABASE ${TEST_DB_NAME} TO ${TEST_DB_USER};
EOSQL

# Apply schema to each DB
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$DB_NAME" -f /docker-entrypoint-initdb.d/schema/init.sql
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$TEST_DB_NAME" -f /docker-entrypoint-initdb.d/schema/init.sql
