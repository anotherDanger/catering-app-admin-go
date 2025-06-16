CREATE TABLE products (
    id VARCHAR(6) PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    price INT NOT NULL,
    stock INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    modified_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE users (
    id CHAR(36) PRIMARY KEY,
    username VARCHAR(100) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL
);

CREATE TABLE admin (
    id CHAR(36) PRIMARY KEY,
    username VARCHAR(100) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL
);

CREATE TABLE orders (
    id CHAR(36) PRIMARY KEY,
    product_id VARCHAR(6) NOT NULL,
    product_name VARCHAR(100) NOT NULL,
    username VARCHAR(100) NOT NULL,
    quantity INT NOT NULL,
    total DOUBLE NOT NULL,
    status VARCHAR(20) DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (product_id) REFERENCES products(id),
    FOREIGN KEY (username) REFERENCES users(username)
);

CREATE INDEX idx_orders_status ON orders(status);
CREATE INDEX idx_products_price ON products(price);
CREATE INDEX idx_users_username ON users(username);

INSERT INTO products (id, name, description, price, stock) VALUES
('PRD001', 'Produk A', 'Deskripsi produk A', 100000, 50);

INSERT INTO users (id, username, password) VALUES
('f1e2d3c4-b5a6-7890-abcd-ef9876543210', 'user1', 'hashed_password_1');

INSERT INTO admin (id, username, password) VALUES
('e7b8a9d4-3f5a-4c82-b7e2-2c3f49b0e9c1', 'admin', 'hashed_admin_password');

INSERT INTO orders (id, product_id, product_name, username, quantity, total, status) VALUES
('user1_20250616120000', 'PRD001', 'Produk A', 'user1', 2, 200000, 'pending');
