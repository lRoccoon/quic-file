package client

import (
	// "bufio"
	//  "code.google.com/p/mahonia"
	"fmt"
	"io"
	"net"
	"os"
)

// TCPClient 开启客户端
func TCPClient(address ,file string, test bool) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println("connect server fail！")
		return
	}
	if test {
		fi, err := os.Open(".\\client\\"+file)
		if err != nil {
			panic(err)
		}
		fiinfo, err := fi.Stat()
		fmt.Println("the size of file is ", fiinfo.Size(), "bytes") //fiinfo.Size() return int64 type
		//发送文件名
		_, err = conn.Write([] byte(file))
		if err != nil {
			fmt.Println("name send error")
		}
		//发送文件
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
		fi.Close()
	} else if test == false {
		//发送文件名
		_, err = conn.Write([] byte(file))
		if err != nil {
			fmt.Println("name send error")
		}
		fmt.Println("the name you has download is :"+file)
        //接受文件
		fi, err := os.Create(".\\client\\"+file)
		if err != nil {
		    fmt.Println("file create error")
		}
		for {
			data := make([]byte, 1024)
			wd, err := conn.Read(data)
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
		fi.Close()
	}
	conn.Close()
	fmt.Println("client close!")
}