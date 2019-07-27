package server

import (
	"bufio"
	// "code.google.com/p/mahonia"
	"fmt"
	"io"
	"net"
	"os"
)

// TCPServer 开启服务器
func TCPServer(address string, port int, test bool) {
	ip := net.ParseIP(address)
	addr := net.TCPAddr{ip, port, ""}
	for {
		listener, err := net.ListenTCP("tcp", &addr) //TCPListener listen
		if err != nil {
			fmt.Println("Initialize error", err.Error())
			return
		}
		tcpcon, err := listener.AcceptTCP() //TCPConn client
		if err != nil {
			fmt.Println(err.Error())
			//continue
		}
		if test {
			data := make([]byte, 1024)
			wc, err := tcpcon.Read(data)
			fo, err := os.Create("E:\\quic-file\\server\\" + string(data[0:wc]))
			if err != nil {
				fmt.Println("file create error")
			}
			fmt.Println("the name of file is :" + string(data[0:wc]))
			for {
				c, err := tcpcon.Read(data)
				if err != nil {
					fmt.Println("read error")
				}
				if string(data[0:c]) == "filerecvend" {
					fmt.Println("write complete ")
					break
				}
				_, err = fo.Write(data[0:c])
				if err != nil {
					fmt.Println("file write error")
				}
			}
		} else {
			fmt.Println("please input the name you want download")
			reader := bufio.NewReader(os.Stdin)
			input, _, _ := reader.ReadLine()
			fi, _ := os.Open(string(input))
			if err != nil {
				panic(err)
			}
			defer fi.Close()
			fiinfo, err := fi.Stat()
			fmt.Println("the size of file is ", fiinfo.Size(), "bytes")
			//send filename
			_, err = tcpcon.Write(input)
			if err != nil {
				fmt.Println("file name send error")
			}
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
		}
		return
	}
}
