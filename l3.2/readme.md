## Shortener
Сервис сокращения URL с аналитикой

### 1. Запуск

git clone github.com/weedworldpeace/wbtasks && cd l3.2 
docker compose up

### 2. Использование

localhost:8080 - адрес сервера по умолчанию
localhost:5433 - адрес postgres по умолчанию

Создание короткой ссылки

    POST /shorten
    Content-Type: application/json

    {
    "url": "https://example.com/very/long/url"
    }

    {
    "original_url": "https://example.com/very/long/url",
    "short_code": "jxm456"
    }

Редирект

    GET /s/{short_code}

Аналитика

    GET /analytics/{short_code}

    {
    "short_code": "abc123",
    "total_clicks": 150,
    "clicks_by_day": [
        {"date": "2024-01-15", "count": 25},
        {"date": "2024-01-16", "count": 35}
    ],
    "user_agents": [
        "Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/120.0.0.0",
        "Mozilla/5.0 (iPhone) Safari/605.1.15"
    ]
    }
### 3. Тестирование

curl -X POST http://localhost:8080/shorten   -H "Content-Type: application/json"   -d '{
    "url": "https://example.com/very/long/url"
    }'

curl http://localhost:8080/s/{short_code}

curl http://localhost:8080/analytics/{short_code}