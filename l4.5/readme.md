## Оптимизация простого API-сервиса с профилировкой
    Простой сервис, который по эндпоинту заполняет слайс 10^5 чисел и считает их сумму
    API: localhost:8080/sum
### v1 слайс без начальной капасити: аллокация на каждый request, переаллокация при росте
    wrk -t 10 -c 100 -d 10s http://localhost:8080/sum

    Running 10s test @ http://localhost:8080/sum
    10 threads and 100 connections
    Thread Stats   Avg      Stdev     Max   +/- Stdev
        Latency    28.66ms   20.60ms 153.74ms   62.34%
        Req/Sec   370.63     45.19   505.00     69.20%
    36973 requests in 10.02s, 4.34MB read
    Requests/sec:   3688.14
    Transfer/sec:    443.01KB

    NumGC = 2811
    HeapAlloc = 1046936

### v2 слайс с начальной капасити: аллокация на каждый request, нет переаллокаций
    wrk -t 10 -c 100 -d 10s http://localhost:8080/sum

    Running 10s test @ http://localhost:8080/sum
    10 threads and 100 connections
    Thread Stats   Avg      Stdev     Max   +/- Stdev
        Latency     6.49ms    5.80ms  65.01ms   72.97%
        Req/Sec     1.77k   263.70     2.29k    66.30%
    176515 requests in 10.00s, 20.71MB read
    Requests/sec:  17644.91
    Transfer/sec:      2.07MB

    NumGC = 12704 (меньше latency больше alloc больше numgc, хоть сервис и стал быстрее)
    HeapAlloc = 4394744
### v3 слайс с начальной капасити из sync.Pool: аллокация лишь при нехватке в Pool
    wrk -t 10 -c 100 -d 10s http://localhost:8080/sum

    Running 10s test @ http://localhost:8080/sum
    10 threads and 100 connections
    Thread Stats   Avg      Stdev     Max   +/- Stdev
        Latency     1.16ms    1.43ms  26.21ms   87.83%
        Req/Sec    12.53k     2.63k   18.00k    71.30%
    1249494 requests in 10.03s, 146.57MB read
    Requests/sec: 124620.94
    Transfer/sec:     14.62MB

    NumGC = 205
    HeapAlloc = 37039992

### В утилите trace видно что в 1 версии GC вызывается постоянно, в 3 версии только ~50ms, из-за этого также почти не используются syscall. Latency упало в 28 раз, а RPS увеличилось в 33 раза. 