CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,          -- Email пользователя
    password VARCHAR(255) NOT NULL,              -- Зашифрованный пароль
    balance INTEGER NOT NULL DEFAULT 1000,       -- Баланс пользователя (coins)
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Время создания
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP -- Время последнего обновления
);

-- Создание таблицы merch
CREATE TABLE IF NOT EXISTS merch (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    price INTEGER NOT NULL
);

-- Заполнение таблицы merch
INSERT INTO merch (name, price) VALUES
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
ON CONFLICT (name) DO NOTHING;

-- Создание таблицы coin_transfers
CREATE TABLE IF NOT EXISTS coin_transfers (
    id SERIAL PRIMARY KEY,
    from_user_id INTEGER NOT NULL,
    to_user_id INTEGER NOT NULL,
    amount INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Создание таблицы purchases
CREATE TABLE IF NOT EXISTS purchases (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    item_name VARCHAR(255) NOT NULL,
    price INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);