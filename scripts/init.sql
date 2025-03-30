CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    customer_name VARCHAR(255),
    product_name VARCHAR(255),
    quantity INT,
    status VARCHAR(50),
    created_at TIMESTAMP DEFAULT NOW()
);