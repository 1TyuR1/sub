###Subscriptions API###
Коротко: REST‑сервис учёта подписок (CRUDL) и расчёт суммы за период, чистая архитектура (Domain → UseCase → Adapters → Drivers).

Стек: Go 1.24, Chi, PostgreSQL (pgx), Docker Compose, Swagger.

###Запуск:###

cp .env.example .env

docker compose up -d --build

Проверка: curl http://localhost:8080/healthz

Документация: http://localhost:8080/swagger/index.html
