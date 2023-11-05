-- Tạo bảng "users"
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255),
    email VARCHAR(255),
    phone_number VARCHAR(20)
);

-- Tạo bảng "bank_accounts" với khóa ngoại user_id
CREATE TABLE bank_accounts (
    id SERIAL PRIMARY KEY,
    user_id INT,
    opening_date DATE,
    current_balance BIGINT,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

-- Tạo bảng "transactions" với khóa ngoại account_id
CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    type VARCHAR(50),
    amount BIGINT,
    date DATE,
    account_id INT,
    FOREIGN KEY (account_id) REFERENCES bank_accounts(id)
);
