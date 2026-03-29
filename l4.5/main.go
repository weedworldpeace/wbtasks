package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"sync"
)

func main() {
	http.HandleFunc("/sum", badSumHandler)
	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

var p = sync.Pool{
	New: func() interface{} {
		return make([]int, 0, 100000)
	},
}

func badSumHandler(w http.ResponseWriter, r *http.Request) {
	sum := 0
	values := p.Get().([]int)
	defer p.Put(values[:0])
	for i := 0; i < 100000; i++ {
		sum += i
		values = append(values, i)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(fmt.Sprintf("sum=%d\n", sum)))
}
