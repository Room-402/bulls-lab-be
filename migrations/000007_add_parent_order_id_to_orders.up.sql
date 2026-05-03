ALTER TABLE orders ADD COLUMN parent_order_id INTEGER REFERENCES orders(id);
