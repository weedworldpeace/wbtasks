## Микросервис “Календарь событий”
    CRUD событий с фоновым архиватором и фоновым нотифаером

### 1. Запуск
    git clone github.com/weedworldpeace/wbtasks && cd l4.3 && docker compose up -d && make migrate-up

### 2. Описание
    Приложение состоит из основного сервиса, принимающего запросы, и бд.
    Можно создавать события, обновлять, удалять, получать уведомления.
    Компоненты системы:

#### App Service (Порт: 8080)
    Принимает HTTP-запросы от пользователей
    Создаёт записи в PostgreSQL
    Фоново отправляет уведомления на почту
    Фоново архивирует старые события(в другую таблицу)
    Асинхронно логирует
#### PostgreSQL (Порт: 5432)
    Хранит информацию о событиях
#### MailHog (Порт: 1025 UI: 8025)

### 3. API
    POST   /create_event?email=<email> - создать событие
    POST   /update_event - обновить событие
    POST   /delete_event - удалить событие
    GET    /events_for_day?date=<timestamp>&user_id=<user_id>      - листинг всех событий на день
    GET    /events_for_week?date=<timestamp>&user_id=<user_id>      - листинг всех событий на неделю
    GET    /events_for_month?date=<timestamp>&user_id=<user_id>      - листинг всех событий на месяц

### 4. Сущности
    type Event struct {
        EventId   string    `json:"event_id"`
        Message   string    `json:"message" binding:"required"`
        Date      time.Time `json:"date" binding:"required"`
        CreatedAt time.Time `json:"created_at"`
        UpdatedAt time.Time `json:"updated_at"`
    }

    type UserEvent struct {
        UserId string `json:"user_id"`
        Event
    }
### 5. Тестирование
    curl -X POST http://127.0.0.1:8080/create_event?email=example@example.com -d '{"user_id": "445393bc-0863-4514-8272-167e54096140","message": "cre",      "date": "2026-03-04T15:00:00Z"}'  -  создание события

    curl -X POST http://127.0.0.1:8080/update_event?email=example@example.com -d '{"user_id": "445393bc-0863-4514-8272-167e54096140","event_id": "07fa23a8-ede2-419b-a0c7-9bf8d3813477","message": "upt","date": "2026-03-02T14:00:00Z"}'  -  обновление события

    curl -X POST http://127.0.0.1:8080/delete_event?email=example@example.com -d '{"user_id": "445393bc-0863-4514-8272-167e54096140","event_id": "07fa23a8-ede2-419b-a0c7-9bf8d3813477","message": "upt","date": "2026-03-02T14:00:00Z"}'  -  удаление события

    curl http://127.0.0.1:8080/events_for_week?user_id=445393bc-0863-4514-8272-167e54096140\&date=2026-03-01T08:18:03Z  -  получение событий