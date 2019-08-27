package main

import (
//	"bufio"
	// "code.google.com/p/mahonia"
	"fmt"
	"io"
	"net"
	"os"
	"log"
	"strconv"
	"path/filepath"
)


// TCPServer 开启服务器
func TCPServer(address string, port string, test bool)string {
	ip := net.ParseIP(address)
	por, err := strconv.Atoi(port)
	addr := net.TCPAddr{ip, por, ""}

	listener, err := net.ListenTCP("tcp", &addr) //TCPListener listen
	if err != nil {
		log.Fatalf("Initialize error", err.Error())
	} else {
		result="the server is listening\n"+"the address is "+address+" "+port+"\n"
		rs.SetText(result)
	}
	for {
		tcpcon, err := listener.AcceptTCP() //TCPConn client
		defer tcpcon.Close()
		if err != nil {
			log.Fatalf(err.Error())
		}
		if test {
			//接受文件名
			data := make([]byte, 1024)
			wc, err := tcpcon.Read(data)
			file :=string(data[0:wc])
			filename :=filepath.Base(file)
			//接受文件
			fi, err := os.Create(filename)
			if err != nil {
				log.Fatalf("file create error")
			}
			for {
				data := make([]byte, 1024)
				wd, err := tcpcon.Read(data)
				if err != nil {
					log.Fatalf("connection read error")
				}
				if string(data[0:wd]) == "filerecvend" {
					break
				}
				_, err = fi.Write(data[0:wd])
				if err != nil {
					log.Fatalf("file write error")
					break
				}
			}	
			result="the name of file you upload is " +filename+" \n"
			rs.Append(result)		
		} else {
			//接受文件名
		    data := make([]byte, 1024)
			wc, err := tcpcon.Read(data)
			result=string(data[0:wc])
			fi, err := os.Open(result)
			if err != nil {
				log.Fatalf("file open error")
			}
			result=("the name of file you download is " + string(data[0:wc]))
			rs.Append(result)
			//发送文件
			buff := make([]byte, 1024)
			for {
				n, err := fi.Read(buff)
				if err != nil && err != io.EOF {
					panic(err)
				}
				if n == 0 {
					tcpcon.Write([]byte("filerecvend"))
					break
				}
				_, err = tcpcon.Write(buff)
				if err != nil {
					fmt.Println("write error")
				}
			}
			fi.Close()
		}
	}
}
  
