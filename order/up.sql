CREATE TABLE IF NOT EXISTS orders (
    id CHAR(27) PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    account_id CHAR(27) NOT NULL,
    total_price NUMERIC(19, 2) NOT NULL
);

CREATE TABLE IF NOT EXISTS order_products (
    order_id CHAR(27) REFERENCES orders (id) ON DELETE CASCADE,
    product_id CHAR(27),
    quantity INT NOT NULL,
    PRIMARY KEY (order_id, product_id)
);