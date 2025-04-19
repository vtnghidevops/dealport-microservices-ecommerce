#!/bin/bash

set -e
set -u

function create_user_and_database() {
    local database=$1
    echo "Creating database '$database'"
    psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" <<-EOSQL
        CREATE DATABASE $database;
        GRANT ALL PRIVILEGES ON DATABASE $database TO $POSTGRES_USER;
EOSQL
}

# Create both users and products databases
echo "Creating multiple databases..."
create_user_and_database "users"
create_user_and_database "products"

# Run SQL initialization scripts for products database
echo "Initializing products database schema..."
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname="products" -f /docker-entrypoint-initdb.d/init-db.sql

echo "PostgreSQL initialization completed" 