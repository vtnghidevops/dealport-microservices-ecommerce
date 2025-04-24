#!/bin/bash

set -e
set -u

echo "PostgreSQL initialization script starting..."

# Create the 'products' database if it doesn't exist
echo "Creating products database..."
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" <<-EOSQL
  SELECT 'Creating database products' AS "Info";
  DROP DATABASE IF EXISTS products;
  CREATE DATABASE products;
  GRANT ALL PRIVILEGES ON DATABASE products TO $POSTGRES_USER;
EOSQL

echo "Products database created successfully"

# Run initialization script for products database
echo "Initializing products database schema and data..."
set +e  # Allow errors to continue for diagnostic purposes
psql -v ON_ERROR_STOP=0 --username "$POSTGRES_USER" --dbname="products" -f /docker-entrypoint-initdb.d/init-db.sql
RESULT=$?
set -e

if [ $RESULT -ne 0 ]; then
  echo "WARNING: There were errors during products database initialization"
  echo "Please check the PostgreSQL logs for details"
else 
  echo "Products database initialization completed successfully"
fi

# Verify data was inserted
echo "Verifying product data..."
psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname="products" <<-EOSQL
  SELECT COUNT(*) AS "Total products" FROM products;
  SELECT COUNT(*) AS "Total categories" FROM categories;
EOSQL

echo "PostgreSQL initialization completed" 