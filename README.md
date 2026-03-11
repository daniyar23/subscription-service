# subscription-service
test service ( Effective Mobile )

# Subscription Service API

API для управления подписками пользователей. Позволяет создавать, читать, обновлять и удалять подписки, а также получать суммы по фильтрам.

## Технологии
- Go + Gin
- PostgreSQL
- Swagger (OpenAPI 3.1)

## Эндпоинты 

| Метод | Путь | Описание |
|-------|------|----------|
| POST  | /subscriptions | Создать подписку |
| GET   | /subscriptions | Получить все подписки |
| GET   | /subscriptions/{id} | Получить подписку по ID |
| PUT   | /subscriptions/{id} | Обновить подписку |
| DELETE| /subscriptions/{id} | Удалить подписку |
| GET   | /subscriptions/user/{user_id} | Все подписки пользователя |
| GET   | /subscriptions/sum | Сумма подписок с фильтром |

## Пример запроса
```json
{
  "service_name": "kino",
  "price": 330,
  "user_id": "60601fee-2bf1-4721-ae6f-763679a0c5ba",
  "start_date": "08-2025"
}


## swagger (OpenAPI)
api/openapi.yaml вставить в свагер едитор
