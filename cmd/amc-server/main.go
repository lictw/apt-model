package main

import (
	"apt-model/server"
	"bufio"
	"os"
	"strings"
	"time"
)

func main() {
	server.Start()
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		command := strings.ToLower(scanner.Text())
		if command == "start" {
			server.Start()
		}
		if command == "stop" {
			server.Stop()
		}
		if command == "exit" {
			os.Exit(0)
		}
	}
	time.Sleep(time.Hour)
}

