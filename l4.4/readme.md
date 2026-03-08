# GC Analyzer
    Анализ сборщика мусора и памяти в Go. Предоставляет метрики в формате Prometheus.

### Запуск
    git clone github.com/weedworldpeace/wbtasks && cd l4.4 && go run main.go

### Эндпоинты 
    http://localhost:6060/metrics
    http://localhost:6060/alloc - аллокация памяти для изменения метрик