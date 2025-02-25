# Avito shop. Магазин мерча

## Описание проекта

Внутренний магазин мерча для сотрудников Авито, в котором сотрудники могут приобретать товары за монеты. Каждому новому сотруднику выдается 1000 монет, которые можно использовать для покупки товаров. Кроме того монеты можно передавать другим сотрудникам в знак благодарности или в подарок.

Магазин поддерживает следующие действия:
- **Аутентификация и авторизация.**
Аутентификация и авторизация сотрудника по его имени и паролю. При первом входе происходит регистрация пользователя - создается аккаунт с указанным паролем, при последующих входах потребуется указать этот пароль.

- **Покупка товара.**
Сотрудник может купить товар, который находится в продаже. Сотрудник не может купить товар, если у него не хватает средств. 
Список доступных к покупке товаров
  | Название     | Цена |
  |--------------|------|
  | t-shirt      | 80   |
  | cup          | 20   |
  | book         | 50   |
  | pen          | 10   |
  | powerbank    | 200  |
  | hoody        | 300  |
  | umbrella     | 200  |
  | socks        | 10   |
  | wallet       | 50   |
  | pink-hoody   | 500  |

- **Отправка средств.**
Сотрудник может передать часть имеющихся средств другому сотруднику. Сотрудник не может отправить больше средств, чем у него есть. Можно отправить средства только зарегистрированному в системе коллеге.  

- **Получение информации о пользователе.**
Сотрудник может получить информацию своих о монетах, инвентаре и истории транзакций.

## Технологический стек

- **Язык:** Go 1.23.3
- **База данных:** Postgres 13
- **WEB-фреймворк:** Gin 1.10.0
- **Аутентификация:** JWT (секрет передается через переменные окружения)
- **Документация:** swaggo/swag 1.16.4
- **Деплой:** Docker Compose

## Архитектура проекта

- **`cmd/main.go`** – точка входа для запуска сервера.
- **`db/`** 
  - `structs.go` - ДТО
  - `/migrations` – миграции базы данных.
- **`docs/`** - документация API
- **`web-server`**
  - **`/handlers`** - слой клиента: контроллеры, их инициализация, и middleware
  - **`/repository`** - слой данных
  - **`/servers`** - слой бизнес-логики

## Инструкция по запуску
1. Настройка переменных окружения: при необходимости поменять переменные окружения в `.env` и `config/local.yaml`
Например: 
```
POSTGRES_PASSWORD=мой_пароль
JWT_SECRET=мой_секрет
```
2. Запустить сервис с помощью команды
```
make build && make run
```
3. (опционально) Для облегчения пользования использовать swagger для общения с сервером - `http://localhost:8080/swagger/index.html`

## Вопросы и проблемы
- В описании API тип AuthRequest содержит только имя и пароль. Из-за этого реализовал регистрацию идентификацию только по имени, что исключает возможность существования нескольких сотрудников с одинаковыми именами
- Т.к. имя пользователя является его уникальным идентификатором, то сделал соответсвующий параметр первичным ключом. В ином случае, первичным ключом был бы id (serial или uuid)
- Для разделения миграций и логики вынес вставку товаров в отдельный shell скрипт `create_products.sh`, который вызывается после билда и запуска проекта
- Так как мы храним минимальную информацию о транзакциях, то таблица `purchases` фактически является инвентарем. При расширении логики сервиса, я бы разделил таблицу на инвентарь и транзакции. 
