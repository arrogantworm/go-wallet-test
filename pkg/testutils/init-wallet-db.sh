#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "user" --dbname "database" <<-EOSQL
    CREATE TABLE IF NOT EXISTS wallets (id UUID PRIMARY KEY, balance NUMERIC NOT NULL DEFAULT 0.0 CHECK (balance >= 0));
EOSQL