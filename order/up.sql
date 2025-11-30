CREATE TABLE IF NOT EXISTS orders (
    id VARCHAR(27) PRIMARY KEY,
    account_id VARCHAR(27) NOT NULL,
    total_price NUMERIC(19, 2) NOT NULL,
    payment_status INT NOT NULL DEFAULT 0,
    payment_id VARCHAR(127),
    `status` INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL
);

CREATE TABLE IF NOT EXISTS order_products (
    order_id VARCHAR(27) NOT NULL,
    product_id CHAR(27),
    quantity INT NOT NULL,
    PRIMARY KEY (order_id, product_id),
    FOREIGN KEY (order_id) REFERENCES orders(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS product_sellers (
    product_id VARCHAR(27) UNIQUE NOT NULL,
    seller_id VARCHAR(27) NOT NULL,
    PRIMARY KEY (product_id, seller_id)
);