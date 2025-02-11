CREATE TYPE productTitle AS
    enum ('t-shirt', 'cup', 'book', 'pen', 'powerbank',
    'hoody', 'umbrella', 'socks', 'wallet', 'pink-hoody');

CREATE TABLE "products" (
    "id" serial PRIMARY KEY,
    "title" productTitle UNIQUE NOT NULL,
    "price" integer NOT NULL CHECK (price >= 0)
);

INSERT INTO "products" (title, price)
VALUES
    ('t-shirt'::productTitle, 80),
    ('cup'::productTitle, 20),
    ('book'::productTitle, 50),
    ('pen'::productTitle, 10),
    ('powerbank'::productTitle, 200),
    ('hoody'::productTitle, 300),
    ('umbrella'::productTitle, 200),
    ('socks'::productTitle, 10),
    ('wallet'::productTitle, 50),
    ('pink-hoody'::productTitle, 500);

CREATE TABLE "employees" (
    "id" serial PRIMARY KEY,
    "username" varchar NOT NULL,
    "password" varchar NOT NULL,
    "balance" int NOT NULL CHECK (balance >= 0) DEFAULT (1000)
);

-- индекс для ускорения поиска по имени пользователя сотрудника
CREATE INDEX employees_username_idx
ON employees(username);

CREATE TABLE "coin_transactions" (
    "id" serial PRIMARY KEY,
    "from_employee_id" int REFERENCES employees(id),
    "to_employee_id" int REFERENCES employees(id),
    "amount" int CHECK (amount > 0),
    "date" timestamp NOT NULL
);

-- индекс для ускорения поиска по имени пользователя сотрудника который передает монеты
CREATE INDEX coin_transactions_from_employee_id_idx
ON coin_transactions(from_employee_id);

-- индекс для ускорения поиска по имени пользователя сотрудника которому передаются монеты
CREATE INDEX coin_transactions_to_employee_id_idx
ON coin_transactions(to_employee_id);

CREATE TABLE "purchases" (
     "id" serial PRIMARY KEY,
     "employee_id" int REFERENCES employees(id),
     "product_id" int REFERENCES products(id),
     "date" timestamp NOT NULL
);

-- индекс для ускорения поиска по айди пользователя сотрудника при выводе информации о купленном мерче
CREATE INDEX purchases_employee_id_idx
ON purchases(employee_id);
