package main

import (
	"flag"
	"log"

	"github.com/873314461/quic-file/client"
	"github.com/873314461/quic-file/server"
)

func main() {
	isServer := flag.Bool("s", false, "server mode")
	isClient := flag.Bool("c", false, "client mode")
	flag.Parse()

	if (*isServer && *isClient) || (!*isServer && !*isClient) {
		log.Fatalln("server or client?")
	}
	if *isServer {
		server.Server("[::]:8000")
	}
	if *isClient {
		client.Client("127.0.0.1:8000", "send.bin")
	}

}
