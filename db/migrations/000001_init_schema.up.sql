CREATE TABLE "products" (
    "id" serial PRIMARY KEY,
    "title" varchar UNIQUE NOT NULL,
    "price" integer NOT NULL CHECK (price >= 0)
);

CREATE TABLE "employees" (
--     "id" serial PRIMARY KEY,
    "username" varchar PRIMARY KEY,
    "password_hash" bytea NOT NULL,
    "balance" int NOT NULL CHECK (balance >= 0) DEFAULT (1000)
);

-- -- индекс для ускорения поиска по имени пользователя сотрудника
-- CREATE INDEX employees_username_idx
-- ON employees(username);

CREATE TABLE "coin_transactions" (
    "id" serial PRIMARY KEY,
    "from_employee" varchar REFERENCES employees(username),
    "to_employee" varchar REFERENCES employees(username),
    "amount" int CHECK (amount > 0),
    "date" timestamp NOT NULL
);

-- индекс для ускорения поиска по имени пользователя сотрудника который передает монеты
CREATE INDEX coin_transactions_from_employee_idx
ON coin_transactions(from_employee);

-- индекс для ускорения поиска по имени пользователя сотрудника которому передаются монеты
CREATE INDEX coin_transactions_to_employee_idx
ON coin_transactions(to_employee);

CREATE TABLE "purchases" (
    "employee" varchar REFERENCES employees(username),
    "product_id" int REFERENCES products(id),
    "amount" int NOT NULL,

    PRIMARY KEY (employee, product_id)
);

-- индекс для ускорения поиска по имени пользователя сотрудника при выводе информации о купленном мерче
CREATE INDEX purchases_employee_idx
ON purchases(employee);
