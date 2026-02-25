## SalesTracker
    CRUD с аналитикой и агрегированием данных

### 1. Запуск
    git clone github.com/weedworldpeace/wbtasks && cd l3.6 && docker compose up -d && make migrate-up

### 2. Описание
    Приложение состоит из основного сервиса, принимающего запросы и бд.
    Вся аналитика агрегируется внутри sql запросов
    Компоненты системы:

#### App Service (Порт: 8080)
    Принимает HTTP-запросы от пользователей
    Создаёт записи в PostgreSQL
    Отдает фронтенд
    Отдает аналитику
#### PostgreSQL (Порт: 5432)
    Хранит информацию о транзакциях 

### 3. API
    ФРОНТ  /                          - листинг транзакций и аналитика, возможность фильтрации по дате, создания, изменения и удаления транзакции
    POST   /api/v1/transactions/      - создать транзакцию
    GET    /api/v1/transactions/      - листинг всех транзакций(?from<>&to=<>)
    GET    /api/v1/transactions/{id}  - информация о транзакции
    PUT    /api/v1/transactions/{id}  - изменить некоторые поля транзакции
    DELETE /api/v1/transactions/{id}  - удалить транзакцию
    GET    /api/v1/analytics          - аналитика(?from<>&to=<>)

### 4. Сущности
    type Transaction struct {
	ID          string          `json:"id" db:"id" `
	UserID      string          `json:"user_id" db:"user_id"`
	Amount      float64         `json:"amount" db:"amount"`
	Type        TransactionType `json:"type" db:"type"`
	Category    string          `json:"category" db:"category"`
	Description string          `json:"description" db:"description"`
	CreatedAt   time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at" db:"updated_at"`
    }
