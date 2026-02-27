## WarehouseControl
    CRUD склада с историей и ролями

### 1. Запуск
    git clone github.com/weedworldpeace/wbtasks && cd l3.7 && docker compose up -d && make migrate-up

### 2. Описание
    Приложение состоит из основного сервиса, принимающего запросы, и бд.
    Вся история работы с данными записывается с помощью триггеров внутри бд.
    Существует разграничение прав по ролям(admin, manager, viewer(не может create, update, delete)).
    Компоненты системы:

#### App Service (Порт: 8080)
    Принимает HTTP-запросы от пользователей
    Создаёт записи в PostgreSQL
    Отдает фронтенд
#### PostgreSQL (Порт: 5432)
    Хранит информацию о продуктах
    Записывает историю взаимодействия с помощью триггеров 

### 3. API
    ФРОНТ  /home               - листинг продуктов, истории и возможность создания, изменения и удаления сущности
    ФРОНТ  /auth               - авторизация с выбором роли(фронтенд редиректит все запросы сюда если бэкенд дает Forbidden(нет jwt cookie))
    POST   /api/v1/items/      - создать продукт
    GET    /api/v1/items/      - листинг всех продуктов(?from<>&to=<>)
    GET    /api/v1/items/{id}  - информация о продукте
    PUT    /api/v1/items/{id}  - изменить продукт
    DELETE /api/v1/items/{id}  - удалить продукт
    GET    /api/v1/history     - история(?from<>&to=<>)
    GET    /api/v1/auth        - получение куки с токеном

### 4. Сущности
    type Item struct {
        ID          string    `json:"id" db:"id"`
        Name        string    `json:"name" db:"name"`
        Description string    `json:"description" db:"description"`
        Quantity    int       `json:"quantity" db:"quantity"`
        Price       float64   `json:"price" db:"price"`
        CreatedAt   time.Time `json:"created_at" db:"created_at"`
        UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
    }

    type ItemHistory struct {
        ID        string    `json:"id" db:"id"`
        ItemID    string    `json:"item_id" db:"item_id"`
        Action    string    `json:"action" db:"action"`
        UserID    string    `json:"user_id" db:"user_id"`
        UserRole  string    `json:"user_role" db:"user_role"`
        OldData   *Item     `json:"old_data,omitempty" db:"old_data"`
        NewData   *Item     `json:"new_data,omitempty" db:"new_data"`
        ChangedAt time.Time `json:"changed_at" db:"changed_at"`
    }
