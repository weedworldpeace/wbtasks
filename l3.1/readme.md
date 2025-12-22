## DelayedNotifier 
Сервис для отложенной отправки уведомлений через RabbitMQ.

## Основные возможности

- **Отложенная отправка** уведомлений в указанное время
- **Поддержка RabbitMQ** для надежной доставки сообщений
- **Автоматические повторные попытки** с экспоненциальной задержкой
- **REST API** для управления уведомлениями


### 1. Запуск

git clone github.com/weedworldpeace/wbtasks && cd l3.1 
docker compose up

### 2. Использование

localhost:8080 - адрес сервера по умолчанию
localhost:8025 - адрес ui mailhog по умолчанию

Создание уведомления

    POST /notify
    Content-Type: application/json

    {
    "email": "user@example.com",
    "data": "Напоминание о встрече",
    "sending_date": "2024-01-15T14:30:00+03:00"
    }

    {
    "id": "550e8400-e29b-41d4-a716-446655440000"
    }

Получение статуса уведомления

    GET /notify/{id}

    {
    "email": "user@example.com",
    "id": "231b8979-a0f6-43e2-96d5-3df99991a0a1",
    "creation_date": "2025-12-22T09:46:52.016859642Z",
    "sending_date": "2026-01-30T01:49:30+03:00",
    "data": ""
    }

Отмена уведомления

    DELETE /notify/{id}

    {
    "message": "notification deleted"
    }
### 3. Тестирование

curl -X POST http://localhost:8080/notify   -H "Content-Type: application/json"   -d '{
    "email": "test@example.com",
    "data": "Тестовое уведомление",
    "sending_date": "'$(date -u +"%Y-%m-%dT%H:%M:%SZ")'" 
  }'

curl http://localhost:8080/notify/{id}

curl -X DELETE http://localhost:8080/notify/{id}

### 4. Примечание

Очередь отложенных уведомлений реализована с помощью плагина github.com/rabbitmq/rabbitmq-delayed-message-exchange и соответственно несет поставленные им лимиты.