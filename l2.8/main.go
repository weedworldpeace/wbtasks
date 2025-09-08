package main

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/beevik/ntp"
)

func main() {
	lg := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	response, err := ntp.Query("0.beevik-ntp.pool.ntp.org")
	if err != nil {
		lg.Error(fmt.Sprintf("ntp query error: %v", err))
		os.Exit(1)
	}
	time := time.Now().Add(response.ClockOffset)
	fmt.Println(time)
}
