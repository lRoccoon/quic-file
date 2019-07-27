package client

import (
	"bufio"
	//  "code.google.com/p/mahonia"
	"fmt"
	"io"
	"net"
	"os"
)

// TCPClient 开启客户端
func TCPClient(address string, test bool) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println("connect server fail！")
		return
	}
	defer conn.Close()
	if test {
		fmt.Println("send file to the destination,please input  filename:")
		reader := bufio.NewReader(os.Stdin)
		input, _, _ := reader.ReadLine()
		fi, err := os.Open(string(input))
		if err != nil {
			panic(err)
		}
		defer fi.Close()
		fiinfo, err := fi.Stat()
		fmt.Println("the size of file is ", fiinfo.Size(), "bytes") //fiinfo.Size() return int64 type
		//send filename
		_, err = conn.Write(input)
		if err != nil {
			fmt.Println("name send error")
		}
		buff := make([]byte, 1024)
		for {
			n, err := fi.Read(buff)
			if err != nil && err != io.EOF {
				panic(err)
			}
			if n == 0 {
				conn.Write([]byte("filerecvend"))
				break
			}
			_, err = conn.Write(buff)
			if err != nil {
				fmt.Println(err.Error())
			}
		}
	} else if test == false {
		data := make([]byte, 1024)
		wc, err := conn.Read(data)
		fi, err := os.Create("E:\\quic-file\\client\\" + string(data[0:wc]))
		if err != nil {
			fmt.Println("file create error")
		}
		fmt.Println("the name of file is " + string(data[0:wc]))
		for {
			wd, err := conn.Read(data)
			if err != nil {
				fmt.Println("connection read error")
			}
			if string(data[0:wd]) == "filerecvend" {
				fmt.Println("file write complete")
				break
			}
			_, err = fi.Write(data[0:wd])
			if err != nil {
				fmt.Println("file write error")
				break
			}
		}
	}
	conn.Close()
	fmt.Println("client close!")
}
