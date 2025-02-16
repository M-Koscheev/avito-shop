#!/bin/bash

source .env

# Подключение к базе данных и выполнение SQL-запросов
#psql -h "${POSTGRES_HOST}" -p 5432 -U "${POSTGRES_USER}" -d "${POSTGRES_DB}" PGPASSWORD "${POSTGRES_PASSWORD}" <<EOF
psql postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_HOST}:5432/${POSTGRES_DB}?sslmode=disable <<EOF
INSERT INTO "products" (title, price)
VALUES
    ('t-shirt', 80),
    ('cup', 20),
    ('book', 50),
    ('pen', 10),
    ('powerbank', 200),
    ('hoody', 300),
    ('umbrella', 200),
    ('socks', 10),
    ('wallet', 50),
    ('pink-hoody', 500)
ON CONFLICT DO NOTHING;
EOF

echo "create products script executed"