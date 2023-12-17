package main

import (
	"github.com/c12s/kuiper/startup"
)

func main() {
	server := startup.NewServer()
	server.Start()

}
