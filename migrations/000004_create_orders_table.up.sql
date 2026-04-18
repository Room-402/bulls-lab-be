CREATE TABLE IF NOT EXISTS orders (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    stock_ticker VARCHAR(100) NOT NULL,
    order_type VARCHAR(100) NOT NULL,
    order_category VARCHAR(100) NOT NULL,
    product_type VARCHAR(100) NOT NULL,
    quantity INTEGER NOT NULL,
    execution_type VARCHAR(100) NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    trigger_price DECIMAL(10, 2) NULL,
    order_status VARCHAR(100) NOT NULL,
    active BOOLEAN DEFAULT true NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_orders_user_id ON orders(user_id);
