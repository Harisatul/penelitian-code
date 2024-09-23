CREATE TABLE IF NOT EXISTS tickets
(
    id          SERIAL PRIMARY KEY,
    order_id    BIGINT,
    category_id SMALLINT NOT NULL,
    version     INT      NOT NULL
);

CREATE INDEX tickets_order_id ON tickets (order_id);
CREATE INDEX tickets_category_id ON tickets (category_id);
CREATE INDEX tickets_version ON tickets (version);