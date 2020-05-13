package main

import (
	"flag"
	"log"

	pingpong "github.com/jemgunay/echo-pingpong"
)

func main() {
	// parse flags
	port := flag.Uint64("port", 3000, "the port for the HTTP server to listen on")
	skillID := flag.String("skill_id", "", "the port for the HTTP server to listen on")

	flag.Parse()

	if *skillID == "" {
		log.Fatal("skill id flag not set")
	}

	log.Printf("starting HTTP server on port %d", *port)
	pingpong.Start(int(*port), *skillID)
}
