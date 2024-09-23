CREATE TYPE ORDER_STATUS AS ENUM ('created', 'completed', 'cancelled');

CREATE TABLE IF NOT EXISTS orders
(
    id          BIGINT PRIMARY KEY,
    category_id SMALLINT     NOT NULL,
    job_id      BIGINT       NOT NULL REFERENCES river_job (id) ON DELETE SET NULL,
    email       VARCHAR(255) NOT NULL,
    status      ORDER_STATUS DEFAULT 'created',
    created_at  TIMESTAMP    DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP    DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX orders_unique ON orders (email, status) WHERE NOT (status = 'cancelled');

CREATE TABLE IF NOT EXISTS tickets
(
    id          SERIAL PRIMARY KEY,
    order_id    BIGINT   REFERENCES orders (id) ON DELETE SET NULL,
    category_id SMALLINT NOT NULL,
    row         SMALLINT NOT NULL,
    "column"    SMALLINT NOT NULL,
    updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX tickets_category_id ON tickets (category_id);

ALTER SYSTEM SET max_connections = '7000';

-- show max conn
SHOW max_connections;

UPDATE tickets
SET order_id = NULL
WHERE true;

TRUNCATE orders CASCADE;

TRUNCATE river_job;

SELECT COUNT(*)
FROM tickets
WHERE order_id IS NOT NULL;

SELECT COUNT(*)
FROM orders
WHERE status = 'completed';