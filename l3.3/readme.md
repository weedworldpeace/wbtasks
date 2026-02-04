## CommentTree API

Сервис древовидных комментариев с поиском и неограниченной вложенностью. Поддерживает создание, удаление, поиск и навигацию по дереву комментариев.

### 1. Запуск

git clone github.com/weedworldpeace/wbtasks && cd l3.3 
docker compose up
make migrate-up

### 2. Использование

localhost:8080 - адрес сервера по умолчанию
localhost:5433 - адрес postgres по умолчанию

localhost:8080/ - фронтенд
localhost:8080/comments - бэкенд

Создание комментария

    curl -X POST http://localhost:8080/comments \
    -H "Content-Type: application/json" \
    -d '{
        "content": "Привет, это первый комментарий!"
    }'

    {
    "id": "123e4567-e89b-12d3-a456-426614174000",
    "parent_id": null,
    "content": "Привет, это первый комментарий!",
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z",
    "depth": 0
    }

Получение комментариев

    curl "http://localhost:8080/comments?page=1&limit=20&sort_by=created_at&order=desc" 

    curl "http://localhost:8080/api/comments?parent=123e4567-e89b-12d3-a456-426614174000&page=1&limit=10" (с родителем)

    curl "http://localhost:8080/api/comments?query=важный&page=1&limit=20" (поиск)


    {
    "comments": [
        {
        "id": "123e4567-e89b-12d3-a456-426614174000",
        "parent_id": null,
        "content": "Привет, это первый комментарий!",
        "created_at": "2024-01-15T10:30:00Z",
        "updated_at": "2024-01-15T10:30:00Z",
        "depth": 0
        }
    ],
    "page": 1,
    "limit": 20,
    "total": 1,
    "total_pages": 1
    }

Удаление

    curl -X DELETE http://localhost:8080/comments/123e4567-e89b-12d3-a456-426614174000

    204 No Content
