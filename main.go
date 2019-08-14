package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/873314461/quic-file/client"
	"github.com/873314461/quic-file/server"
)

func main() {

	isTCP := flag.Bool("t", false, "TCP connection")
	isQUIC := flag.Bool("q", false, "QUIC connection")

	isServer := flag.Bool("s", false, "server mode")
	isClient := flag.Bool("c", false, "client mode")

	isUpload := flag.Bool("u", false, "upload file")
	isDownload := flag.Bool("d", false, "download file")
	flag.Parse()

	if (*isTCP && *isQUIC) || (!*isTCP && !*isQUIC) {
		log.Fatalln("Tcp  or  Quic")
	}
	if (*isServer && *isClient) || (!*isServer && !*isClient) {
		log.Fatalln("server or client?")
	}
	if (*isUpload && *isDownload) || (!*isUpload && !*isDownload) {
		log.Fatalln("upload or download?")
	}
	var file string
	if *isQUIC {
		if *isUpload {
			if *isServer {
				server.Server("[::]:8000", true)
			}
			if *isClient {
				fmt.Println("please input the name you want to upload and the filename is in the client directory")
				fmt.Scanf("%s", &file)
				client.Client("127.0.0.1:8000", file, true)
			}
		}
		if *isDownload {
			if *isServer {
				server.Server("[::]:8000", false)
			}
			if *isClient {
				fmt.Println("please input the name you want to download and the filename is in the server directory")
				fmt.Scanf("%s", &file)
				client.Client("127.0.0.1:8000", file, false)
			}
		}
	}
	if *isTCP {
		if *isUpload {
			if *isServer {
				server.TCPServer("[::]", 9090, true)
			}
			if *isClient {
				fmt.Println("please input the name you want to upload and the filename is in the client directory")
				fmt.Scanf("%s", &file)
				client.TCPClient("127.0.0.1:9090", file, true)
			}
		}
		if *isDownload {
			if *isServer {
				server.TCPServer("[::]", 9090, false)
			}
			if *isClient {
				fmt.Println("please input the name you want to download and the filename is in the server directory")
				fmt.Scanf("%s", &file)
				client.TCPClient("127.0.0.1:9090", file, false)
			}
		}
	}
}
