package main

import (
	"fmt"
	"flag"
	"log"
	"strings"
	
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
	var file1 , file2 string
	if *isDownload{
        fmt.Println("please input the name you want to download and the filename is in the server directory")
		fmt.Scanf(file1)
		fmt.Println("please input the name after in the download and the filename is in the current directory,default is same as filename download")
		fmt.Scanf(file2)
		i :=strings.Compare(file2,"")
		if i==0{
			file2=file1
		}
	} else if *isUpload {
		fmt.Println("please input the name you want to upload and the filename is in the current directory")
		fmt.Scanf(file1)
		fmt.Println("please input the name after in the upload and the filename is in the server directory,default is same as filename upload")
		fmt.Scanf(file2)
		i :=strings.Compare(file2,"")
		if i==0{
			file2=file1
		}
	}
	if *isQUIC {
		if *isUpload {
			if *isServer {
				server.Server("[::]:8000",file2,true)
			}
			if *isClient {
				client.Client("127.0.0.1:8000",file1,true)
			}
		}
		if *isDownload {
			if *isServer {
				server.Server("[::]:8000",file1,false)
			}
			if *isClient {
				client.Client("127.0.0.1:8000",file2,false)
			}
		}
	}
	if *isTCP {
		if *isUpload {
			if *isServer {
				server.TCPServer("[::]", file2,9090 ,true)
			}
			if *isClient {
				client.TCPClient("127.0.0.1:9090",file1,true)
			}
		}
		if *isDownload {
			if *isServer {
				server.TCPServer("[::]",file1, 9090, false)
			}
			if *isClient {
				client.TCPClient("127.0.0.1:9090",file2,false)
			}
		}
	}
}
