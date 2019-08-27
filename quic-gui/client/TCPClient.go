package main

import (
	// "bufio"
	//"code.google.com/p/mahonia"
	"strconv"
	"log"
	"io"
	"net"
	"os"
	"github.com/andlabs/ui"
	_"github.com/andlabs/ui/winmanifest"
)


// TCPClient 开启客户端
func TCPClient(address ,file string, test bool)  {
	conn, err := net.Dial("tcp", address)
	defer conn.Close()
	if err != nil {
		log.Fatalf("connect server fail！")
	}
	if test {
		fi, err := os.Open(file)
		if err != nil {
			panic(err)
		}
		fiinfo, err := fi.Stat()
		sendByte:= strconv.FormatInt(fiinfo.Size(),10)
		result =" file send "+sendByte+" bytes\n"
		rs.SetText(result)
		//发送文件名
		_, err = conn.Write([] byte(file))
		if err != nil {
			log.Fatalf("name send error")
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
				log.Fatalf(err.Error())
			}
		}
		fi.Close()
	} else if test == false {
		//发送文件名
		_, err = conn.Write([] byte(file))
		if err != nil {
			log.Fatalf("name send error")
		}
		result="the name you has download is :"+file+"\n"
		rs.Append(result)
        //接受文件
		fi, err := os.Create(file)
		if err != nil {
		    log.Fatalf("file create error")
		}
		for {
			data := make([]byte, 1024)
			wd, err := conn.Read(data)
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
		fi.Close() 
	}
	ui.MsgBox(mainwin,
		"congratulation ,operation completed",
		"please first close this window")
}