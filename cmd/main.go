package main

import (
	"log/slog"
	"time"
)

func main() {
	slog.Info("check 1")

	num := 1
	for {
		time.Sleep(10 * time.Second)
		slog.Info("message", "number", num)
		num++
	}
}
