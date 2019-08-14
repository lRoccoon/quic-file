package server

import (
	"bufio"
	// "code.google.com/p/mahonia"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

// TCPServer 开启服务器
func TCPServer(address string, port int, test bool) {
	ip := net.ParseIP(address)
	addr := net.TCPAddr{ip, port, ""}

	listener, err := net.ListenTCP("tcp", &addr) //TCPListener listen
	if err != nil {
	fmt.Println("Initialize error", err.Error())
		return
	} else {
		fmt.Println("the server is listening")
	}
	for {
		tcpcon, err := listener.AcceptTCP() //TCPConn client
		defer tcpcon.Close()
		if err != nil {
			fmt.Println(err.Error())
		}
		if test {
			//接受文件名
			data := make([]byte, 1024)
			wc, err := tcpcon.Read(data)
			fmt.Println("the name of file you upload is :" + string(data[0:wc]))
			//接受文件
			fi, err := os.Create(".\\server\\"+string(data[0:wc]))
			if err != nil {
				fmt.Println("file create error")
			}
			for {
				data := make([]byte, 1024)
				wd, err := tcpcon.Read(data)
				if err != nil {
					fmt.Println("connection read error")
				}
				if string(data[0:wd]) == "filerecvend" {
					break
				}
				_, err = fi.Write(data[0:wd])
				if err != nil {
					fmt.Println("file write error")
					break
				}
			}			
		} else {
			//接受文件名
		    data := make([]byte, 1024)
			wc, err := tcpcon.Read(data)
			fi, err := os.Open(".\\server\\"+string(data[0:wc]))
			if err != nil {
				fmt.Println("file open error")
			}
			fmt.Println("the name of file you download is " + string(data[0:wc]))
			//发送文件
			fiinfo, err := fi.Stat()
			fmt.Println("the size of file is ", fiinfo.Size(), "bytes")
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
		fmt.Println("file transmission end,and if you want to transmission again,please input ok,else input exit")
		inputReader := bufio.NewReader(os.Stdin)
		input, err := inputReader.ReadString('\n')
		var ex string
		ex ="exit\r\n"
		ts :=strings.Compare(input,ex)
		if ts==0{
			return 
		} else{
			fmt.Println("the server is listening")
		}
	}
}
  
