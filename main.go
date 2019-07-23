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

	isUpload := flag.Bool("u", false, "upload file")
	isDownload := flag.Bool("d", false, "download file")
	flag.Parse()

	if (*isServer && *isClient) || (!*isServer && !*isClient) {
		log.Fatalln("server or client?")
	}
	if (*isUpload && *isDownload) || (!*isUpload && !*isDownload) {
		log.Fatalln("upload or download?")
	}
	if *isUpload {
		if *isServer {
			server.Server("[::]:8000", "uploads.txt", true)
		}
		if *isClient {
			client.Client("127.0.0.1:8000", "upload.txt", true)
		}
	}
	if *isDownload {
		if *isServer {
			server.Server("[::]:8000", "download.txt", false)
		}
		if *isClient {
			client.Client("127.0.0.1:8000", "downloads.txt", false)
		}
	}
}
