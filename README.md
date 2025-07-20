# Wallet App

REST API-сервис для работы с кошельками: пополнение, списание и получение баланса.

---

## Стек

- Golang
- PostgreSQL
- Docker + Docker Compose

---

## Переменные окружения `config.env`

```env
DB_HOST=postgres
DB_PORT=5432
DB_USER=db_user
DB_PASS=db_pass
DB_NAME=wallet_db
DB_SSLMODE=disable
```

---

## Сборка и запуск

```commandline
docker compose --env-file config.env -f docker-compose.yaml up --build
```